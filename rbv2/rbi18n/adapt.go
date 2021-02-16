package rbi18n

import (
	"net/http"

	rb "github.com/gohandle/rb/rbv2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type bundle struct{ *i18n.Bundle }

// Adapt adapts an i18n.Bundle to allow translation
func Adapt(b *i18n.Bundle) rb.TranslateCore {
	return bundle{b}
}

func (b bundle) Translate(w http.ResponseWriter, r *http.Request, mid string, opts ...rb.TranslateOption) string {
	var o rb.TranslateOpts
	for _, opt := range opts {
		opt(&o)
	}

	loc := i18n.NewLocalizer(b.Bundle, r.Header.Get("Accept-Language"))

	s, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:   mid,
		PluralCount: o.PluralCount,
	})

	if err != nil {
		rb.L(r).Error("failed to translate",
			zap.Error(err), zap.String("message_id", mid), zap.Any("translate_opts", o))
		return mid
	}

	return s
}
