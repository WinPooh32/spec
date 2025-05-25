package spec_test

import (
	"fmt"
	"testing"

	"github.com/WinPooh32/spec"
	"github.com/stretchr/testify/assert"
)

type Man struct {
	Age    int
	Mortal bool
	Diet   []string
}

type ImmortalSpec struct{}

func (ImmortalSpec) SatisfiedBy(man Man) error {
	if !man.Mortal {
		return nil
	}

	return fmt.Errorf("expected to be immortal")
}

type AncientSpec struct{}

func (AncientSpec) SatisfiedBy(man Man) error {
	const n = 500

	if man.Age >= n {
		return nil
	}

	return fmt.Errorf("expected to be older than %d years", n)
}

type DietSpec struct {
	Meals map[string]struct{}
}

func (diet DietSpec) SatisfiedBy(man Man) error {
	for _, meal := range man.Diet {
		if _, ok := diet.Meals[meal]; !ok {
			return fmt.Errorf("expected diet: %v", diet.Meals)
		}
	}

	return nil
}

func ExampleSpecification() {
	// Define composite specification.
	vampireSpec := spec.And(
		ImmortalSpec{},
		AncientSpec{},
		DietSpec{map[string]struct{}{"blood": {}}},
	)

	// Define man instance.
	man := Man{
		Age:    200,
		Mortal: true,
		Diet:   []string{"meat", "apple"},
	}

	// Test the man is a vampire.
	fmt.Println(vampireSpec.SatisfiedBy(man))

	// Output:
	// expected to be immortal
	// expected to be older than 500 years
	// expected diet: map[blood:{}]
}

type testMan struct {
	BirthYear int
	Name      string
	Age       int
	Mortal    bool
}

type NewtonSpec struct{}

func (NewtonSpec) SatisfiedBy(man testMan) error {
	if man.Name == "Newton" {
		return nil
	}

	return fmt.Errorf("expected to be named as %q", "Newton")
}

type OldSpec struct{}

func (OldSpec) SatisfiedBy(man testMan) error {
	const n = 20

	if man.Age > n {
		return nil
	}

	return fmt.Errorf("expected to be older than %d", n)
}

type MortalSpec struct{}

func (MortalSpec) SatisfiedBy(man testMan) error {
	if man.Mortal {
		return nil
	}

	return fmt.Errorf("expected to be mortal")
}

type BirthYearSpec struct {
	Year int
}

func (birth BirthYearSpec) SatisfiedBy(man testMan) error {
	if man.BirthYear == birth.Year {
		return nil
	}

	return fmt.Errorf("expected to born at %d", man.BirthYear)
}

func TestAnd_xy(t *testing.T) {
	oldManSpec := spec.And(
		OldSpec{},
		MortalSpec{},
	)

	man := testMan{
		Age:    50,
		Mortal: true,
	}

	err := oldManSpec.SatisfiedBy(man)
	assert.NoError(t, err)
}

func TestAnd_xy_not_satisfied(t *testing.T) {
	oldManSpec := spec.And(
		OldSpec{},
		MortalSpec{},
	)

	man := testMan{
		Age:    50,
		Mortal: false,
	}

	err := oldManSpec.SatisfiedBy(man)
	assert.Error(t, err)
}

func TestAnd_xyzz(t *testing.T) {
	IsaacNewtonSpec := spec.And(
		OldSpec{},
		NewtonSpec{},
		MortalSpec{},
		BirthYearSpec{Year: 1643},
	)

	man := testMan{
		Name:      "Newton",
		Age:       25,
		Mortal:    true,
		BirthYear: 1643,
	}

	err := IsaacNewtonSpec.SatisfiedBy(man)
	assert.NoError(t, err)
}

func TestAnd_xyzz_z_not_satisfied(t *testing.T) {
	modernNewtonSpec := spec.And(
		OldSpec{},
		NewtonSpec{},
		MortalSpec{},
		BirthYearSpec{Year: 2000},
	)

	man := testMan{
		Name:      "Newton",
		Age:       25,
		Mortal:    true,
		BirthYear: 1643,
	}

	err := modernNewtonSpec.SatisfiedBy(man)
	assert.Error(t, err)
}

func TestAnd_xyzz_y_not_satisfied(t *testing.T) {
	modernNewtonSpec := spec.And(
		OldSpec{},
		NewtonSpec{},
		MortalSpec{},
		BirthYearSpec{Year: 2000},
	)

	man := testMan{
		Name:      "Newton",
		Age:       10,
		Mortal:    true,
		BirthYear: 2000,
	}

	err := modernNewtonSpec.SatisfiedBy(man)
	assert.Error(t, err)
}
