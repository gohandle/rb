package rbjet

import (
	"io"
	"reflect"

	"github.com/CloudyKit/jet/v6"
	rb "github.com/gohandle/rb/rbv2"
)

type adapted struct{ *jet.Set }

// Adapt adapts a jset tmplate set to implement the template directory
func Adapt(jset *jet.Set) rb.Templates {
	return adapted{jset}
}

func (a adapted) Lookup(n string) (rb.TemplateExecuter, error) {
	jset, err := a.Set.GetTemplate(n)
	if err != nil {
		return nil, err
	}

	return exec{jset}, nil
}

type exec struct{ *jet.Template }

func (e exec) Execute(w io.Writer, vars map[string]reflect.Value, data interface{}) (err error) {
	return e.Template.Execute(w, vars, data)
}
