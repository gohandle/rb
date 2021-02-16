package rbplayg

import (
	"context"
	"net/url"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	rb "github.com/gohandle/rb/rbv2"
)

type decoder struct{ *form.Decoder }

// AdaptDecoder adapts a form decoder to implement the rb values decoder
func AdaptDecoder(fdec *form.Decoder) rb.ValuesDecoder {
	return decoder{fdec}
}

func (d decoder) DecodeValues(v interface{}, values url.Values) error {
	return d.Decode(v, values)
}

type val struct{ *validator.Validate }

// AdaptValidator will adapt the playground validator to the rb validator struct
func AdaptValidator(v *validator.Validate) rb.StructValidator {
	return val{v}
}

func (vld val) ValidateStruct(ctx context.Context, v interface{}) error {
	return vld.Validate.StructCtx(ctx, v)
}
