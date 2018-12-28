package main

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {

	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateBody will decode the body into payload and then validate against the payloads struct tags
func ValidateBody(payload interface{}, body io.Reader) error {

	err := json.NewDecoder(body).Decode(payload)
	// request.Body is empty then we'll get an EOF, we'll let the validator handle this case
	if err.Error() == "EOF" {
	} else if err != nil {
		return err
	}

	return validate.Struct(payload)
}
