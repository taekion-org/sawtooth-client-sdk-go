package types

// BatchHeader represents a Sawtooth batch header.
type BatchHeader struct {
	SignerPublicKey	string			`json:"signer_public_key"`
	TransactionIds	[]string		`json:"transaction_ids"`
}

// Batch represents a Sawtooth batch.
type Batch struct {
	Header			BatchHeader 	`json:"header"`
	HeaderSignature	string     		`json:"header_signature"`
	Trace			bool         	`json:"trace"`
	Transactions	[]Transaction 	`json:"transactions"`
}

// BatchStatus represents a Sawtooth batch status.
type BatchStatus string

const (
	BATCH_STATUS_PENDING   BatchStatus = "PENDING"
	BATCH_STATUS_COMMITTED BatchStatus = "COMMITTED"
	BATCH_STATUS_INVALID   BatchStatus = "INVALID"
	BATCH_STATUS_UNKNOWN   BatchStatus = "UNKNOWN"
)
