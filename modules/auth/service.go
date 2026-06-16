package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"ecommerce/modules/customer"
	"ecommerce/packages/events"
	"ecommerce/packages/types"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Ping(ctx context.Context) error
	Register(ctx context.Context, req RegisterReq) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginReq) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshReq) (*RefreshResponse, error)
	Logout(ctx context.Context, authIdentityID string) error
	ValidateToken(ctx context.Context, tokenStr string) (string, string, error)
}

type service struct {
	repo      *Repository
	custRepo  *customer.Repository
	bus       events.EventBus
	jwtSecret []byte
	jwtIssuer string
}

func NewService(repo *Repository, custRepo *customer.Repository, bus events.EventBus, jwtSecret string) Service {
	return &service{repo: repo, custRepo: custRepo, bus: bus, jwtSecret: []byte(jwtSecret), jwtIssuer: "medusa"}
}

type customClaims struct {
	CustomerID string `json:"customerId"`
	jwt.RegisteredClaims
}

const (
	accessTokenTTL  = 30 * time.Minute
	refreshTokenTTL = 30 * 24 * time.Hour
)

func (s *service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *service) Register(ctx context.Context, req RegisterReq) (*RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, types.NewValidation("email and password are required")
	}
	if len(req.Password) < 6 {
		return nil, types.NewValidation("password must be at least 6 characters")
	}

	existing, _ := s.repo.GetProviderByEntityID(ctx, "emailpass", req.Email)
	if existing != nil {
		return nil, types.NewConflict("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	authID := types.GenerateID("auth")
	customerID := types.GenerateID("cus")

	appMeta, _ := json.Marshal(map[string]string{"customerId": customerID})

	identity := &AuthIdentity{
		ID:          authID,
		AppMetadata: appMeta,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.InsertAuthIdentity(ctx, identity); err != nil {
		return nil, fmt.Errorf("insert auth identity: %w", err)
	}

	providerMeta, _ := json.Marshal(map[string]string{"password": string(hashed)})
	provider := &ProviderIdentity{
		ID:               types.GenerateID("prov"),
		AuthIdentityID:   authID,
		Provider:         "emailpass",
		EntityID:         req.Email,
		ProviderMetadata: providerMeta,
		UserMetadata:     []byte("{}"),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.repo.InsertProviderIdentity(ctx, provider); err != nil {
		return nil, fmt.Errorf("insert provider identity: %w", err)
	}

	c := &customer.Customer{
		ID:         customerID,
		Email:      req.Email,
		FirstName:  &req.FirstName,
		LastName:   &req.LastName,
		HasAccount: true,
		Metadata:   []byte("{}"),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.custRepo.Insert(ctx, c); err != nil {
		return nil, fmt.Errorf("insert customer: %w", err)
	}

	return &RegisterResponse{
		Message: "registered successfully",
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginReq) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, types.NewValidation("email and password are required")
	}

	provider, err := s.repo.GetProviderByEntityID(ctx, "emailpass", req.Email)
	if err != nil {
		return nil, types.NewValidation("invalid email or password")
	}

	var meta struct {
		Password string `json:"password"`
	}
	_ = json.Unmarshal(provider.ProviderMetadata, &meta)

	if err := bcrypt.CompareHashAndPassword([]byte(meta.Password), []byte(req.Password)); err != nil {
		return nil, types.NewValidation("invalid email or password")
	}

	identity, err := s.repo.GetAuthIdentityByID(ctx, provider.AuthIdentityID)
	if err != nil {
		return nil, err
	}

	var appMeta struct {
		CustomerID string `json:"customerId"`
	}
	_ = json.Unmarshal(identity.AppMetadata, &appMeta)

	accessToken, err := s.generateAccessToken(identity.ID, appMeta.CustomerID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, identity.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, req RefreshReq) (*RefreshResponse, error) {
	if req.RefreshToken == "" {
		return nil, types.NewValidation("refresh token is required")
	}

	hash := sha256Hex(req.RefreshToken)
	stored, err := s.repo.GetAuthTokenByHash(ctx, hash)
	if err != nil {
		return nil, types.NewValidation("invalid refresh token")
	}

	if stored.RevokedAt != nil {
		return nil, types.NewValidation("token has been revoked")
	}

	if time.Now().UTC().After(stored.ExpiresAt) {
		return nil, types.NewValidation("token has expired")
	}

	_ = s.repo.RevokeAuthToken(ctx, stored.ID)

	identity, err := s.repo.GetAuthIdentityByID(ctx, stored.AuthIdentityID)
	if err != nil {
		return nil, err
	}

	var appMeta struct {
		CustomerID string `json:"customerId"`
	}
	_ = json.Unmarshal(identity.AppMetadata, &appMeta)

	accessToken, err := s.generateAccessToken(identity.ID, appMeta.CustomerID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.generateRefreshToken(ctx, identity.ID)
	if err != nil {
		return nil, err
	}

	return &RefreshResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *service) Logout(ctx context.Context, authIdentityID string) error {
	return s.repo.RevokeAuthTokensByAuthID(ctx, authIdentityID)
}

func (s *service) ValidateToken(ctx context.Context, tokenStr string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return "", "", types.NewValidation("invalid or expired token")
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok || !token.Valid {
		return "", "", types.NewValidation("invalid token claims")
	}

	return claims.Subject, claims.CustomerID, nil
}

func (s *service) generateAccessToken(authID, customerID string) (string, error) {
	now := time.Now().UTC()
	claims := customClaims{
		CustomerID: customerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   authID,
			Issuer:    s.jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *service) generateRefreshToken(ctx context.Context, authID string) (string, error) {
	token := types.GenerateID("rtok")
	hash := sha256Hex(token)
	now := time.Now().UTC()

	t := &AuthToken{
		ID:             types.GenerateID("atok"),
		AuthIdentityID: authID,
		TokenHash:      hash,
		Type:           "refresh",
		ExpiresAt:      now.Add(refreshTokenTTL),
		CreatedAt:      now,
	}
	if err := s.repo.InsertAuthToken(ctx, t); err != nil {
		return "", err
	}
	return token, nil
}

func sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}


