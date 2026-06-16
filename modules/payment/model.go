package payment

import "time"

type PaymentCollectionStatus string

const (
	PaymentCollectionStatusNotPaid           PaymentCollectionStatus = "not_paid"
	PaymentCollectionStatusAwaiting          PaymentCollectionStatus = "awaiting"
	PaymentCollectionStatusAuthorized        PaymentCollectionStatus = "authorized"
	PaymentCollectionStatusPartiallyCaptured PaymentCollectionStatus = "partially_captured"
	PaymentCollectionStatusCaptured          PaymentCollectionStatus = "captured"
	PaymentCollectionStatusCancelled         PaymentCollectionStatus = "cancelled"
	PaymentCollectionStatusRefunded          PaymentCollectionStatus = "refunded"
)

const (
	PaymentSessionStatusPending      PaymentSessionStatus = "pending"
	PaymentSessionStatusAuthorized   PaymentSessionStatus = "authorized"
	PaymentSessionStatusRequiresMore PaymentSessionStatus = "requires_more"
	PaymentSessionStatusError        PaymentSessionStatus = "error"
	PaymentSessionStatusCancelled    PaymentSessionStatus = "cancelled"
)

type PaymentSessionStatus string

// PaymentCollection — umbrella for all payments of 1 order. OrderID soft ref.
type PaymentCollection struct {
	ID               string                  `db:"id" json:"id"`
	OrderID          *string                 `db:"order_id" json:"orderId"` // soft ref
	CartID           *string                 `db:"cart_id" json:"cartId"`   // soft ref
	CurrencyCode     string                  `db:"currency_code" json:"currencyCode"`
	Amount           int64                   `db:"amount" json:"amount"`
	AuthorizedAmount *int64                  `db:"authorized_amount" json:"authorizedAmount"`
	CapturedAmount   *int64                  `db:"captured_amount" json:"capturedAmount"`
	RefundedAmount   *int64                  `db:"refunded_amount" json:"refundedAmount"`
	Status           PaymentCollectionStatus `db:"status" json:"status"`
	CreatedAt        time.Time               `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time               `db:"updated_at" json:"updatedAt"`
}

// PaymentSession — 1 attempt with 1 provider. Data: provider-specific.
type PaymentSession struct {
	ID                  string               `db:"id" json:"id"`
	PaymentCollectionID string               `db:"payment_collection_id" json:"paymentCollectionId"`
	ProviderID          string               `db:"provider_id" json:"providerId"`
	Status              PaymentSessionStatus `db:"status" json:"status"`
	Amount              int64                `db:"amount" json:"amount"`
	CurrencyCode        string               `db:"currency_code" json:"currencyCode"`
	Data                []byte               `db:"data" json:"-"`
	CreatedAt           time.Time            `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time            `db:"updated_at" json:"updatedAt"`
}

type Payment struct {
	ID                  string     `db:"id" json:"id"`
	PaymentCollectionID string     `db:"payment_collection_id" json:"paymentCollectionId"`
	ProviderID          string     `db:"provider_id" json:"providerId"`
	CurrencyCode        string     `db:"currency_code" json:"currencyCode"`
	Amount              int64      `db:"amount" json:"amount"`
	AuthorizedAt        *time.Time `db:"authorized_at" json:"authorizedAt"`
	CapturedAt          *time.Time `db:"captured_at" json:"capturedAt"`
	CancelledAt         *time.Time `db:"cancelled_at" json:"cancelledAt"`
	Data                []byte     `db:"data" json:"-"`
	CreatedAt           time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updatedAt"`
}

type Refund struct {
	ID        string    `db:"id" json:"id"`
	PaymentID string    `db:"payment_id" json:"paymentId"`
	Amount    int64     `db:"amount" json:"amount"`
	Note      *string   `db:"note" json:"note"`
	CreatedBy *string   `db:"created_by" json:"createdBy"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
