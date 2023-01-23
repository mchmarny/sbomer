package report

import (
	"github.com/pkg/errors"
)

type Request struct {
	Path   string
	Target string
}

func (r *Request) Validate() error {
	if r.Path == "" {
		return errors.New("no input provided")
	}
	if r.Target == "" {
		return errors.New("no target provided")
	}
	return nil
}
