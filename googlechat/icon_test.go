package googlechat_test

import (
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestIcon_Validate(t *testing.T) {
	v := validator.New()
	googlechat.RegisterValidation(v)
	for _, test := range []struct {
		name  string
		icon  googlechat.Icon
		valid bool
	}{
		{"known icon", googlechat.Icon{KnownIcon: "EMAIL"}, true},
		{"HTTPS icon", googlechat.Icon{IconURL: "https://example.com/icon.png"}, true},
		{"material icon", googlechat.Icon{MaterialIcon: &googlechat.MaterialIcon{Name: "mail", Weight: 700, Grade: -25}}, true},
		{"missing source", googlechat.Icon{}, false},
		{"two sources", googlechat.Icon{KnownIcon: "EMAIL", IconURL: "https://example.com/icon.png"}, false},
		{"HTTP icon", googlechat.Icon{IconURL: "http://example.com/icon.png"}, false},
		{"invalid weight", googlechat.Icon{MaterialIcon: &googlechat.MaterialIcon{Name: "mail", Weight: 150}}, false},
		{"invalid grade", googlechat.Icon{MaterialIcon: &googlechat.MaterialIcon{Name: "mail", Grade: 100}}, false},
		{"unknown image type", googlechat.Icon{KnownIcon: "EMAIL", ImageType: "TRIANGLE"}, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := v.Struct(test.icon)
			if test.valid {
				require.NoError(t, err)
			} else {
				require.ErrorAs(t, err, new(validator.ValidationErrors))
			}
		})
	}
}
