package types

// BlockHeader represents a Sawtooth block header.
type BlockHeader struct {
	BatchIds			[]string	`json:"block_ids"`
	BlockNum			string		`json:"block_num"`
	Consensus			string		`json:"consensus"`
	PreviousBlockId		string		`json:"previous_block_id"`
	SignerPublicKey		string		`json:"signer_public_key"`
	StateRootHash		string		`json:"state_root_hash"`
}

// Block represents a Sawtooth block.
type Block struct {
	Batches         []Batch       `json:"batches"`
	Header          BlockHeader `json:"header"`
	HeaderSignature string      `json:"header_signature"`
}
