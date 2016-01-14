package errors

import (
	"errors"
	"strings"
)

type Errors struct {
	Errors []string
}

func (errs *Errors) Add(err string) {
	(*errs).Errors = append((*errs).Errors, err)
}

func (errs Errors) HaveError() bool {
	return len(errs.Errors) > 0
}

func (errs Errors) Errorify() error {
	return errors.New(strings.Join(errs.Errors, "\n"))
}
