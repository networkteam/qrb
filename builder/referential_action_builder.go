package builder

// ReferentialActionBuilder provides typed methods for selecting a referential action (CASCADE, RESTRICT, etc.).
// It is parameterized on the parent builder type so that action methods return the correct builder for continued chaining.
type ReferentialActionBuilder[T any] struct {
	parent T
	setter func(T, string) T
}

// Cascade sets the referential action to CASCADE.
func (b ReferentialActionBuilder[T]) Cascade() T { return b.setter(b.parent, "CASCADE") }

// Restrict sets the referential action to RESTRICT.
func (b ReferentialActionBuilder[T]) Restrict() T { return b.setter(b.parent, "RESTRICT") }

// SetNull sets the referential action to SET NULL.
func (b ReferentialActionBuilder[T]) SetNull() T { return b.setter(b.parent, "SET NULL") }

// SetDefault sets the referential action to SET DEFAULT.
func (b ReferentialActionBuilder[T]) SetDefault() T { return b.setter(b.parent, "SET DEFAULT") }

// NoAction sets the referential action to NO ACTION.
func (b ReferentialActionBuilder[T]) NoAction() T { return b.setter(b.parent, "NO ACTION") }
