package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

//Valid returns true if there are no errors, otherwise fales
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}


// Form create a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

//New initializes a form struck
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Has check if form is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}

	return true
}

//Required checks for required fileds
func(f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//MinLength checks for string minimum length
func(f *Form) MinLength(field string, length int, r *http.Request) bool {
	x:= r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

//IsEmail checks for valid address
func(f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}