package matchers

import (
	"fmt"
	"reflect"
)

// CopyInto is an interface that defines a method for copying values into another value.
type CopyInto interface {
	Matcher
	CopyInto(any) error
}

// CopyIntoFunc is a function type that defines how to copy values from one type to another.
type CopyIntoFunc[T any] func(dst T, src T) error

// copyIntoFunc is a matcher that copies the value from a given type using a custom function.
type copyIntoFunc[T any] struct {
	source   T
	copyFunc CopyIntoFunc[T]
}

// Match implements the Matcher interface for copyIntoFunc types.
func (c *copyIntoFunc[T]) Match(actual any) (bool, error) {
	// copy into always matches
	return true, nil
}

// CopyInto implements the CopyInto interface for copyIntoFunc types.
func (c *copyIntoFunc[T]) CopyInto(destination any) error {
	dst := destination.(T)
	return c.copyFunc(dst, c.source)
}

// NewCopyIntoFunc creates a CopyInto matcher that copies the value from a given type using a custom function.
func NewCopyIntoFunc[T any](toCopy T, copyFunc CopyIntoFunc[T]) CopyInto {
	return &copyIntoFunc[T]{
		source:   toCopy,
		copyFunc: copyFunc,
	}
}

// DeepCopyInto is an interface that defines a method for deep copying values into another value.
type DeepCopyInto[T any] interface {
	DeepCopyInto(T)
}

// deepCopyInto is a matcher that copies the value from a DeepCopyInto type.
type deepCopyInto[T DeepCopyInto[T]] struct {
	source T
}

// Match implements the Matcher interface for DeepCopyInto types.
func (d *deepCopyInto[T]) Match(any) (bool, error) {
	return true, nil
}

// CopyInto implements the CopyInto interface for DeepCopyInto types.
func (d *deepCopyInto[T]) CopyInto(destination any) (err error) {
	dst, ok := destination.(T)
	if !ok {
		return fmt.Errorf("destination is not a DeepCopyInto[%T]", reflect.TypeFor[T]())
	}
	d.source.DeepCopyInto(dst)

	return
}

// NewDeepCopyInto creates a CopyInto matcher that copies the value from a DeepCopyInto type.
func NewDeepCopyInto[T DeepCopyInto[T]](toCopy T) CopyInto {
	return &deepCopyInto[T]{
		source: toCopy,
	}
}

// NewCopyIntoPointer creates a CopyInto matcher that copies the value from a pointer.
func NewCopyIntoPointer[T any](toCopy *T) CopyInto {
	return NewCopyIntoFunc[*T](toCopy, func(dst *T, src *T) error {
		if src == nil {
			return nil // nothing to copy
		}
		if dst == nil {
			return fmt.Errorf("destination pointer is nil, cannot copy into it")
		}
		*dst = *src
		return nil
	})
}
