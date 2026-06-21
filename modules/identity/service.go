package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ecommerce/modules/auth"
	"ecommerce/modules/iam"
	"ecommerce/packages/events"
	"ecommerce/packages/types"
	"ecommerce/packages/upload"
)

const maxAvatarSize = 5 << 20 // 5 MB

type Service interface {
	CreateAdmin(ctx context.Context, req CreateAdminReq) (*User, error)
	ListAdmins(ctx context.Context) ([]User, error)
	RevokeAdmin(ctx context.Context, callerID, targetID string) error
	BanUser(ctx context.Context, callerID, targetID string) error
	UnbanUser(ctx context.Context, targetID string) error
	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileReq) (*User, error)
	UploadAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string) (*User, error)
}

type service struct {
	repo     *Repository
	authRepo *auth.Repository
	iamRepo  *iam.Repository
	bus      events.EventBus
	uploader upload.Uploader
}

func NewService(repo *Repository, authRepo *auth.Repository, iamRepo *iam.Repository, bus events.EventBus, uploader upload.Uploader) Service {
	return &service{repo: repo, authRepo: authRepo, iamRepo: iamRepo, bus: bus, uploader: uploader}
}

func (s *service) CreateAdmin(ctx context.Context, req CreateAdminReq) (*User, error) {
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, types.NewValidation("email, password, firstName and lastName are required")
	}
	if len(req.Password) < 6 {
		return nil, types.NewValidation("password must be at least 6 characters")
	}

	existing, _ := s.authRepo.GetProviderByEntityID(ctx, "emailpass", req.Email)
	if existing != nil {
		return nil, types.NewConflict("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	// identity.user.id == auth_identity.id by design
	authID := types.GenerateID("auth")

	identity := &auth.AuthIdentity{
		ID:          authID,
		AppMetadata: []byte(`{}`),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.authRepo.InsertAuthIdentity(ctx, identity); err != nil {
		return nil, fmt.Errorf("insert auth identity: %w", err)
	}

	providerMeta, _ := json.Marshal(map[string]string{"password": string(hashed)})
	provider := &auth.ProviderIdentity{
		ID:               types.GenerateID("prov"),
		AuthIdentityID:   authID,
		Provider:         "emailpass",
		EntityID:         req.Email,
		ProviderMetadata: providerMeta,
		UserMetadata:     []byte("{}"),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.authRepo.InsertProviderIdentity(ctx, provider); err != nil {
		return nil, fmt.Errorf("insert provider identity: %w", err)
	}

	user := &User{
		ID:        authID,
		Email:     req.Email,
		FirstName: &req.FirstName,
		LastName:  &req.LastName,
		Status:    StatusActive,
		Metadata:  []byte("{}"),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.InsertUser(ctx, user); err != nil {
		return nil, fmt.Errorf("insert identity user: %w", err)
	}

	if req.RoleID != nil {
		if _, err := s.iamRepo.GetRoleByID(ctx, *req.RoleID); err != nil {
			return nil, types.NewValidation("roleId not found")
		}
		if err := s.iamRepo.AssignRoleToPrincipal(ctx, authID, *req.RoleID); err != nil {
			return nil, fmt.Errorf("assign role: %w", err)
		}
	}

	return user, nil
}

func (s *service) ListAdmins(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *service) RevokeAdmin(ctx context.Context, callerID, targetID string) error {
	if callerID == targetID {
		return types.NewValidation("cannot revoke your own account")
	}
	target, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target.Status == StatusClosed {
		return types.NewValidation("account is already closed")
	}
	if err := s.repo.SetStatus(ctx, target.ID, StatusClosed); err != nil {
		return fmt.Errorf("set status closed: %w", err)
	}
	if err := s.authRepo.RevokeAuthTokensByAuthID(ctx, target.ID); err != nil {
		return fmt.Errorf("revoke tokens: %w", err)
	}
	return nil
}

func (s *service) BanUser(ctx context.Context, callerID, targetID string) error {
	if callerID == targetID {
		return types.NewValidation("cannot ban your own account")
	}
	target, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target.Status == StatusBan {
		return types.NewValidation("account is already banned")
	}
	if target.Status == StatusClosed {
		return types.NewValidation("cannot ban a closed account")
	}
	if err := s.repo.SetStatus(ctx, target.ID, StatusBan); err != nil {
		return fmt.Errorf("set status ban: %w", err)
	}
	if err := s.authRepo.RevokeAuthTokensByAuthID(ctx, target.ID); err != nil {
		return fmt.Errorf("revoke tokens: %w", err)
	}
	return nil
}

func (s *service) UnbanUser(ctx context.Context, targetID string) error {
	target, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return err
	}
	if target.Status != StatusBan {
		return types.NewValidation("account is not banned")
	}
	if err := s.repo.SetStatus(ctx, target.ID, StatusActive); err != nil {
		return fmt.Errorf("set status active: %w", err)
	}
	return nil
}

func (s *service) GetProfile(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *service) UpdateProfile(ctx context.Context, userID string, req *UpdateProfileReq) (*User, error) {
	if err := s.repo.UpdateProfile(ctx, userID, req); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return s.repo.GetByID(ctx, userID)
}

func (s *service) UploadAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string) (*User, error) {
	if err := upload.ValidateSize(size, maxAvatarSize); err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("identity/%s/avatar%s", userID, upload.ExtFromContentType(contentType))

	url, err := s.uploader.Upload(ctx, "avatars", objectName, file, size, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload avatar: %w", err)
	}

	url = fmt.Sprintf("%s?v=%d", url, time.Now().Unix())

	if err := s.repo.UpdateAvatarURL(ctx, userID, url); err != nil {
		return nil, fmt.Errorf("persist avatar url: %w", err)
	}

	return s.repo.GetByID(ctx, userID)
}
