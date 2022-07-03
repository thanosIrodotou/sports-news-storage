package validator

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	*validator.Validate
	Translator ut.Translator
}

func New() (*Validator, error) {
	english := en.New()
	uni := ut.New(english, english)

	trans, ok := uni.GetTranslator("en")
	if !ok {
		return nil, errors.New("invalid translation")
	}

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	if err := entranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	return &Validator{
		Validate:   v,
		Translator: trans,
	}, nil
}
