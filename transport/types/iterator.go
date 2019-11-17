package types

// CommonIterator is an interface that provides methods common to all iterators.
type CommonIterator interface {
	Next() bool
	Error() error
}

// BlockIterator is an interface that represents an iterator over blocks.
type BlockIterator interface {
	CommonIterator
	Current() (*Block, error)
}

// BatchIterator is an interface that represents an iterator over batches.
type BatchIterator interface {
	CommonIterator
	Current() (*Batch, error)
}

// TransactionIterator is an interface that represents an iterator over transactions.
type TransactionIterator interface {
	CommonIterator
	Current() (*Transaction, error)
}

// StateIterator is an interface that represents an iterator over state.
type StateIterator interface {
	CommonIterator
	Current() (*State, error)
}
