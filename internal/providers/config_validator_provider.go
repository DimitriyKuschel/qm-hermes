package providers

import (
	"github.com/gookit/validate"
	"queue-manager/internal/structures"
)

type CnfValidator struct {
	conf *structures.Config
}

func (c *CnfValidator) Validate() error {
	v := validate.Struct(c.conf)
	if !v.Validate() {
		return v.Errors.OneError()
	}
	return nil
}

func NewCnfValidator(c *structures.Config) *CnfValidator {
	return &CnfValidator{conf: c}
}
