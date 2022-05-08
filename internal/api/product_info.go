// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ProductInfo product info
//
// swagger:model ProductInfo
type ProductInfo struct {

	// basket
	Basket *BasketWithoutRequired `json:"basket,omitempty"`

	// products i ds
	// Required: true
	ProductsIDs *int64 `json:"productsIDs"`

	// quantity
	Quantity int64 `json:"quantity,omitempty"`
}

// Validate validates this product info
func (m *ProductInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBasket(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProductsIDs(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProductInfo) validateBasket(formats strfmt.Registry) error {
	if swag.IsZero(m.Basket) { // not required
		return nil
	}

	if m.Basket != nil {
		if err := m.Basket.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("basket")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("basket")
			}
			return err
		}
	}

	return nil
}

func (m *ProductInfo) validateProductsIDs(formats strfmt.Registry) error {

	if err := validate.Required("productsIDs", "body", m.ProductsIDs); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this product info based on the context it is used
func (m *ProductInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBasket(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProductInfo) contextValidateBasket(ctx context.Context, formats strfmt.Registry) error {

	if m.Basket != nil {
		if err := m.Basket.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("basket")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("basket")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProductInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProductInfo) UnmarshalBinary(b []byte) error {
	var res ProductInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
