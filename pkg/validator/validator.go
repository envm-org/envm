package validator

import (
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func New() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		log.Println("Validator initialized")
	})
	return validate
}
