package entity

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	None Sorting = iota
	IdAsc
	IdDesc
	TitleAsc
	TitleDesc
	PriceAsc
	PriceDesc
	CreatedAtAsc
	CreatedAtDesc
	ModifiedAtAsc
	ModifiedAtDesc
)

var (
	// ErrValidation describes an error which occurred during validation
	ErrValidation = errors.New("validation failed")
	// ErrProductNotFound describes an error which occurred if a product cannot be found
	ErrProductNotFound = errors.New("product not found")
)

type (
	// Sorting describes in which order a list of products will be sorted
	Sorting int

	// Product is the main business object
	Product struct {
		// ID is a unique identifier for the product
		ID string
		// Title is a unique, non-empty title of the product
		Title string
		// Description is a non-empty description of the product
		Description string
		Price       uint64
		CreatedAt   time.Time
		ModifiedAt  time.Time
	}

	// FilterObject contains a filter with which products can be filtered
	FilterObject struct {
		Sort   Sorting
		Limit  uint
		Offset uint
		Free   bool
	}
)

// Validate checks if the product follows the business rules. A ErrValidation will be returned if not.
func (product Product) Validate() error {
	if product.Title == "" {
		return fmt.Errorf("%w: title must be set", ErrValidation)
	}

	if product.Description == "" {
		return fmt.Errorf("%w: description must be set", ErrValidation)
	}

	return nil
}

// NewId generates a new id for the given product
func (product *Product) NewId() {
	product.ID = uuid.NewString()
}

// Validate checks if the FilterObject follows the business rules. A ErrValidation will be returned if not.
func (filter FilterObject) Validate() error {
	if filter.Limit == 0 {
		return fmt.Errorf("%w: limit cannot be zero", ErrValidation)
	}

	return nil
}
