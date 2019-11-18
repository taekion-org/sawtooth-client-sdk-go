package types

// TransactionHeader represents a Sawtooth transaction header.
type TransactionHeader struct {
	BatcherPublicKey	string		`json:"batcher_public_key"`
	Dependencies		[]string	`json:"dependencies"`
	FamilyName			string		`json:"family_name"`
	FamilyVersion		string		`json:"family_version"`
	Inputs				[]string	`json:"inputs"`
	Nonce				string		`json:"nonce"`
	Outputs				[]string	`json:"outputs"`
	PayloadSHA256		string		`json:"payload_sha256"`
	SignerPublicKey		string		`json:"signer_public_key"`
}

// Transaction represents a Sawtooth transaction.
type Transaction struct {
	Header				TransactionHeader	`json:"header"`
	HeaderSignature		string				`json:"header_signature"`
	Payload				[]byte				`json:"payload"`
}
