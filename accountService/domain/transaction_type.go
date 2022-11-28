package domain

type TransactionType int64

const (
	internal_transaction TransactionType = iota
	external_transaction
)
