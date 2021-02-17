package rb

// TranslateOpts configures options for translating a message
type TranslateOpts struct {
	PluralCount interface{}
}

// TranslateOption allows configuring the translate
type TranslateOption func(*TranslateOpts)

// PluralCount configures the plural version of the translated message
func PluralCount(c interface{}) TranslateOption {
	return func(o *TranslateOpts) {
		o.PluralCount = c
	}
}
