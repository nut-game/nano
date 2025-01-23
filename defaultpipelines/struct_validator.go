package defaultpipelines

import (
	"context"
)

// StructValidator is the interface that must be implemented
// by a struct validator for the request arguments on nano.
//
// The default struct validator used by nano is https://github.com/go-playground/validator.
type StructValidator interface {
	Validate(context.Context, interface{}) (context.Context, interface{}, error)
}

// StructValidatorInstance holds the default validator
// on start but can be overridden if needed.
var StructValidatorInstance StructValidator = &DefaultValidator{}
