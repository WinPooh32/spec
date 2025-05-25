# spec

Simple implementation of the [Composite Specification pattern](https://martinfowler.com/apsupp/spec.pdf) ([archive](https://web.archive.org/web/20250428234628/https://martinfowler.com/apsupp/spec.pdf)). Only conjunction operation (and) is implemented for simplicity and better performance.

Also it can provide textual description what specification rule is not satisfied by using an error interface instead of a bool predicate.

## Example

Define class:

```go
type Man struct {
    Age    int
    Mortal bool
    Diet   []string
}
```

Define specifications:

```go
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
```

Test specification is satisfied:

```go
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
```
