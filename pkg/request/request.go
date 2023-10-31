package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		context   echo.Context
		validator *validator.Validate
	}
)

func (c *contextWrapper) Bind(data any) error {
	if err := c.context.Bind(data); err != nil {
		log.Printf("Error: binding data: %v \n", err)
		return err
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("Error: validating data: %v \n", err)
		return err
	}

	return nil
}

func ContextWrapper(context echo.Context) contextWrapperService {
	return &contextWrapper{context, validator.New()}
}
