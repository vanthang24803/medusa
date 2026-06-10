package payment

import (
	"context"
	"fmt"
)

// ProviderType defines the payment gateway provider type.
type ProviderType string

const (
	ProviderStripe ProviderType = "stripe"
	ProviderPayPal ProviderType = "paypal"
	ProviderMoMo   ProviderType = "momo"
	ProviderVNPay  ProviderType = "vnpay"
)

// Config holds the configuration parameters for the payment gateway.
type Config struct {
	Provider  ProviderType
	APIKey    string // Secret API key (for Stripe/PayPal)
	ClientID  string // Public client ID (for PayPal)
	Merchant  string // Merchant ID (for MoMo/VNPay)
	SecretKey string // Encryption/Signature key (for MoMo/VNPay webhook verification)
	Sandbox   bool   // Toggle sandbox/testing environment
}

// PaymentStatus represents the payment status.
type PaymentStatus string

const (
	StatusPending  PaymentStatus = "pending"
	StatusSuccess  PaymentStatus = "success"
	StatusFailed   PaymentStatus = "failed"
	StatusRefunded PaymentStatus = "refunded"
)

// CreatePaymentRequest contains the information needed to initiate a payment transaction.
type CreatePaymentRequest struct {
	OrderID     string
	Amount      int64  // Payment amount in the smallest currency unit (e.g. Cents or VND)
	Currency    string // E.g. "USD", "VND"
	Description string
	ReturnURL   string // Redirect URL after payment completion
	CancelURL   string // Redirect URL when user cancels the payment
}

// CreatePaymentResponse contains the output of the create transaction request (e.g., payment URL).
type CreatePaymentResponse struct {
	TransactionID string
	PaymentURL    string // URL to redirect the user to the payment gateway (MoMo, VNPay, Stripe)
}

// PaymentStatusResponse contains details about the transaction status.
type PaymentStatusResponse struct {
	TransactionID string
	OrderID       string
	Status        PaymentStatus
	Amount        int64
	RawResponse   string // Raw JSON/Response body from the payment gateway
}

// RefundRequest contains the refund request details.
type RefundRequest struct {
	TransactionID string
	Amount        int64
	Reason        string
}

// RefundResponse contains the refund request output.
type RefundResponse struct {
	RefundID string
	Status   string
}

// Payment defines the common interface for payment gateways (MoMo, VNPay, Stripe, PayPal, etc.).
type Payment interface {
	// CreatePayment initiates a payment transaction and returns the payment URL.
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error)

	// GetPaymentStatus retrieves the actual transaction status from the payment gateway.
	GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatusResponse, error)

	// Refund processes a refund for a previously successful transaction.
	Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error)

	// VerifyWebhook validates the signature and parses webhook notifications from the payment gateway.
	VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*PaymentStatusResponse, error)

	// Close cleans up any client connections/resources if applicable.
	Close() error
}

// NewPayment is a Factory function to create the corresponding Payment provider.
func NewPayment(cfg Config) (Payment, error) {
	switch cfg.Provider {
	case ProviderStripe:
		return newStripePayment(cfg)
	case ProviderPayPal:
		return newPayPalPayment(cfg)
	case ProviderMoMo:
		return newMoMoPayment(cfg)
	case ProviderVNPay:
		return newVNPayPayment(cfg)
	default:
		return nil, fmt.Errorf("payment: unsupported provider: %s", cfg.Provider)
	}
}

// ── Stripe Implementation (Stub) ─────────────────────────────────────────────

type stripePayment struct {
	cfg Config
}

func newStripePayment(cfg Config) (*stripePayment, error) {
	// TODO: Initialize Stripe API client
	return &stripePayment{cfg: cfg}, nil
}

func (p *stripePayment) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, nil
}

func (p *stripePayment) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *stripePayment) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, nil
}

func (p *stripePayment) VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *stripePayment) Close() error {
	return nil
}

// ── PayPal Implementation (Stub) ─────────────────────────────────────────────

type paypalPayment struct {
	cfg Config
}

func newPayPalPayment(cfg Config) (*paypalPayment, error) {
	// TODO: Initialize PayPal SDK client
	return &paypalPayment{cfg: cfg}, nil
}

func (p *paypalPayment) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, nil
}

func (p *paypalPayment) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *paypalPayment) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, nil
}

func (p *paypalPayment) VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *paypalPayment) Close() error {
	return nil
}

// ── MoMo Implementation (Stub) ───────────────────────────────────────────────

type momoPayment struct {
	cfg Config
}

func newMoMoPayment(cfg Config) (*momoPayment, error) {
	// TODO: Initialize MoMo API client configurations
	return &momoPayment{cfg: cfg}, nil
}

func (p *momoPayment) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, nil
}

func (p *momoPayment) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *momoPayment) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, nil
}

func (p *momoPayment) VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *momoPayment) Close() error {
	return nil
}

// ── VNPay Implementation (Stub) ──────────────────────────────────────────────

type vnpayPayment struct {
	cfg Config
}

func newVNPayPayment(cfg Config) (*vnpayPayment, error) {
	// TODO: Initialize VNPay API configurations
	return &vnpayPayment{cfg: cfg}, nil
}

func (p *vnpayPayment) CreatePayment(ctx context.Context, req CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, nil
}

func (p *vnpayPayment) GetPaymentStatus(ctx context.Context, transactionID string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *vnpayPayment) Refund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	return nil, nil
}

func (p *vnpayPayment) VerifyWebhook(ctx context.Context, payload []byte, headers map[string]string) (*PaymentStatusResponse, error) {
	return nil, nil
}

func (p *vnpayPayment) Close() error {
	return nil
}
