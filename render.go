package rb

import (
	"net/http"

	"go.uber.org/zap"
)

// HeaderRender is an optional interface that can be implemented by a Render. If it wants to
// change the header code it should return a new code and not call WriteHeader itself.
type HeaderRender interface {
	RenderHeader(rc RenderCore, w http.ResponseWriter, r *http.Request, status int) (int, error)
}

// Render can be implemented to customize how data should be written to a response.
type Render interface {
	Render(rc RenderCore, w http.ResponseWriter, r *http.Request) error
}

// RenderOpts hold all options for a render
type RenderOpts struct {
	Code int
}

// RenderOption configures a render
type RenderOption func(*RenderOpts)

// Status configures any render such that the response header is written with the provided
// status code.
func Status(code int) RenderOption {
	return func(opts *RenderOpts) {
		opts.Code = code
	}
}

type renderCore struct {
	Templates
}

// NewRenderCore creates the render part of the Core
func NewRenderCore(view Templates) RenderCore {
	return &renderCore{view}
}

func (rc *renderCore) Render(w http.ResponseWriter, r *http.Request, rr Render, os ...RenderOption) error {
	var o RenderOpts
	for _, opt := range os {
		opt(&o)
	}

	L(r).Debug("start render", zap.Any("render", rr), zap.Any("options", o))

	if hr, ok := rr.(HeaderRender); ok {
		var err error
		L(r).Debug("render implemented header rendering", zap.Int("status_code", o.Code))
		o.Code, err = hr.RenderHeader(rc, w, r, o.Code)
		if err != nil {
			L(r).Error("error while rendering header", zap.Error(err))
			return err
		}
	}

	if o.Code < 1 {
		o.Code = http.StatusOK
		L(r).Debug("no explicit status provided, set default", zap.Int("status_code", o.Code))
	}

	w.WriteHeader(o.Code)
	L(r).Debug("wrote header", zap.Int("status_code", o.Code))

	if err := rr.Render(rc, w, r); err != nil {
		L(r).Error("error while rendering body", zap.Error(err))
		return err
	}

	L(r).Debug("render complete")
	return nil
}
