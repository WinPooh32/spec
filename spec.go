package spec

import (
	"errors"
)

type Specification[T any] interface {
	SatisfiedBy(value T) error
}

type Conjunction[T any] struct {
	x, y Specification[T]
	z    []Specification[T]
}

func (conj Conjunction[T]) SatisfiedBy(value T) error {
	errx := conj.x.SatisfiedBy(value)
	erry := conj.y.SatisfiedBy(value)

	p := (errx == nil) && (erry == nil)

	if len(conj.z) == 0 {
		if !p {
			return errors.Join(errx, erry)
		}

		return nil
	}

	var errs []error

	if !p {
		errs = append(errs, errx, erry)
	}

	for _, z := range conj.z {
		errz := z.SatisfiedBy(value)

		pz := (errz == nil)
		if !pz {
			errs = append(errs, errz)
		}

		p = p && pz
	}

	return errors.Join(errs...)
}

func And[T any](x, y Specification[T], z ...Specification[T]) Conjunction[T] {
	return Conjunction[T]{x, y, z}
}
