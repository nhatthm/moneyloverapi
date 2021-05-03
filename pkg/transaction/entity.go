package transaction

import "time"

// BankTransaction represents a bank transaction (not moneylover transaction).
// This struct holds needed data to create and import transaction to moneylover.
type BankTransaction struct {
	ID            string    `json:"id"`
	AccountName   string    `json:"accountName"`
	AccountID     string    `json:"accountId"`
	AccountBank   string    `json:"accountBank"`
	Amount        float64   `json:"amount"`
	Category      string    `json:"category"`
	ReferenceText string    `json:"referenceText"`
	DisplayDate   time.Time `json:"displayDate"`
	PartnerName   string    `json:"partnerName"`
	PartnerID     string    `json:"partnerId"`
	PartnerBank   string    `json:"partnerBank"`
}
