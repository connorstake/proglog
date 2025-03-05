package server

// Record is a struct that implements the IRecord interface.
type Record struct {
	Value  []byte `json:"value"` // Changed to string to accept plain text
	Offset uint64 `json:"offset"`
}
