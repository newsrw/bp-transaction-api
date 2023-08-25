package domain

import "context"

// TxStatus transaction status
type TxStatus string

const (
	// TxStatusConfirmed transaction status "CONFIRMED" (Transaction has been processed and confirmed)
	TxStatusConfirmed TxStatus = "CONFIRMED"
	// TxStatusFailed transaction status "FAILED" (Transaction failed to process)
	TxStatusFailed TxStatus = "FAILED"
	// TxStatusPending transaction status "PENDING" (Transaction is awaiting processing)
	TxStatusPending TxStatus = "PENDING"
	// TxStatusDNE transaction status "DNE" (Transaction does not exist)
	TxStatusDNE TxStatus = "DNE"
)

// BroadcastTransactionRequest broadcast transaction request
type BroadcastTransactionRequest struct {
	Symbol    string `json:"symbol" validate:"required"`
	Price     uint64 `json:"price" validate:"required"`
	Timestamp uint64 `json:"timestamp" validate:"required"`
}

// BroadcastTransactionResponse broadcast transaction response
type BroadcastTransactionResponse struct {
	TxHash string `json:"tx_hash"`
}

// MonitoringTransactionStatusRequest monitoring transaction status request
type MonitoringTransactionStatusRequest struct {
	BroadcastTransactionResponse
}

// MonitoringTransactionStatusResponse monitoring transaction status response
type MonitoringTransactionStatusResponse struct {
	TxStatus TxStatus `json:"tx_status"`
}

// TransactionUsecase transaction's usecases interface
type TransactionUsecase interface {
	BroadcastTransaction(ctx context.Context, req *BroadcastTransactionRequest) (rsp BroadcastTransactionResponse, err error)
	MonitoringStatus(ctx context.Context, req *MonitoringTransactionStatusRequest) (rsp *MonitoringTransactionStatusResponse, err error)
}
