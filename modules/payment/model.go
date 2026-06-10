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

// PaymentCollection — umbrella for all payments of 1 order. OrderID soft ref.
type PaymentCollection struct {
	ID               string                  `db:"id" json:"id"`
	OrderID          *string                 `db:"order_id" json:"order_id"` // soft ref
	CartID           *string                 `db:"cart_id" json:"cart_id"`   // soft ref
	CurrencyCode     string                  `db:"currency_code" json:"currency_code"`
	Amount           int64                   `db:"amount" json:"amount"`
	AuthorizedAmount *int64                  `db:"authorized_amount" json:"authorized_amount"`
	CapturedAmount   *int64                  `db:"captured_amount" json:"captured_amount"`
	RefundedAmount   *int64                  `db:"refunded_amount" json:"refunded_amount"`
	Status           PaymentCollectionStatus `db:"status" json:"status"`
	CreatedAt        time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time               `db:"updated_at" json:"updated_at"`
}

type PaymentSessionStatus string

const (
	PaymentSessionStatusPending      PaymentSessionStatus = "pending"
	PaymentSessionStatusAuthorized   PaymentSessionStatus = "authorized"
	PaymentSessionStatusRequiresMore PaymentSessionStatus = "requires_more"
	PaymentSessionStatusError        PaymentSessionStatus = "error"
	PaymentSessionStatusCancelled    PaymentSessionStatus = "cancelled"
)

// PaymentSession — 1 attempt with 1 provider. Data: provider-specific.
type PaymentSession struct {
	ID                  string               `db:"id" json:"id"`
	PaymentCollectionID string               `db:"payment_collection_id" json:"payment_collection_id"`
	ProviderID          string               `db:"provider_id" json:"provider_id"`
	Status              PaymentSessionStatus `db:"status" json:"status"`
	Amount              int64                `db:"amount" json:"amount"`
	CurrencyCode        string               `db:"currency_code" json:"currency_code"`
	Data                []byte               `db:"data" json:"-"`
	CreatedAt           time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time            `db:"updated_at" json:"updated_at"`
}

type Payment struct {
	ID                  string     `db:"id" json:"id"`
	PaymentCollectionID string     `db:"payment_collection_id" json:"payment_collection_id"`
	ProviderID          string     `db:"provider_id" json:"provider_id"`
	CurrencyCode        string     `db:"currency_code" json:"currency_code"`
	Amount              int64      `db:"amount" json:"amount"`
	AuthorizedAt        *time.Time `db:"authorized_at" json:"authorized_at"`
	CapturedAt          *time.Time `db:"captured_at" json:"captured_at"`
	CancelledAt         *time.Time `db:"cancelled_at" json:"cancelled_at"`
	Data                []byte     `db:"data" json:"-"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
}

type Refund struct {
	ID        string    `db:"id" json:"id"`
	PaymentID string    `db:"payment_id" json:"payment_id"`
	Amount    int64     `db:"amount" json:"amount"`
	Note      *string   `db:"note" json:"note"`
	CreatedBy *string   `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
