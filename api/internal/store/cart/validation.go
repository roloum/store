package cart

import (
	"errors"

	validator "github.com/go-playground/validator/v10"
)

const (
	//ErrCartIDIsEmpty Error describes when cartID is empty
	ErrCartIDIsEmpty = "CartIDIsEmpty"

	//ErrItemIDIsEmpty Error describes when itemID is empty
	ErrItemIDIsEmpty = "ItemIDIsEmpty"

	//ErrDescriptionIsEmpty Error describes when Item description is empty
	ErrDescriptionIsEmpty = "DescriptionIsEmpty"

	//ErrPriceIsEmpty Error describes when price is empty
	ErrPriceIsEmpty = "PriceIsEmpty"

	//ErrPriceIsInvalid Error describes when price is not valid number
	ErrPriceIsInvalid = "PriceIsInvalid"

	//ErrQuantityIsEmpty Error describes when Quantity is empty
	ErrQuantityIsEmpty = "QuantityIsEmpty"

	//ErrQuantityIsInvalid Error describes when quantity is not a valid number
	ErrQuantityIsInvalid = "QuantityIsInvalid"
)

var validate *validator.Validate

//init instantiates a validator
func init() {
	validate = validator.New()

	validate.RegisterValidation("validPrice", isValidPrice)
	validate.RegisterValidation("validQuantity", isValidQuantity)

}

//getValidationError Returns the first error reported by the validator
func getValidationError(verr error) error {

	//Retrieve first error
	err := verr.(validator.ValidationErrors)[0]

	switch err.Field() {
	case "CartID":
		return errors.New(ErrCartIDIsEmpty)
	case "ItemID":
		return errors.New(ErrItemIDIsEmpty)
	case "Description":
		return errors.New(ErrDescriptionIsEmpty)
	case "Price":
		switch err.Tag() {
		case "required":
			return errors.New(ErrPriceIsEmpty)
		case "validPrice":
			return errors.New(ErrPriceIsInvalid)
		}
	case "Quantity":
		switch err.Tag() {
		case "required":
			return errors.New(ErrQuantityIsEmpty)
		case "validQuantity":
			return errors.New(ErrQuantityIsInvalid)
		}
	}
	return nil
}

//isValidPrice Checks that the item's price is a valid number
func isValidPrice(fl validator.FieldLevel) bool {

	if price := fl.Field().Float(); price < 0 {
		return false
	}

	return true
}

//isValidQuantity Checks the item's quantity is a valid number
func isValidQuantity(fl validator.FieldLevel) bool {
	if quantity := fl.Field().Int(); quantity < 1 {
		return false
	}
	return true
}
