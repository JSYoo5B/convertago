package googlechat

import (
	"encoding/json"
	"fmt"
)

// Icon displays a built-in icon, a custom HTTPS image, or a Material icon.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#Icon
type Icon struct {
	// KnownIcon names a built-in icon. Supply exactly one icon source.
	KnownIcon string `json:"knownIcon,omitempty" validate:"required_without_all=IconURL MaterialIcon,excluded_with=IconURL MaterialIcon"`
	// IconURL hosts a custom PNG or JPG icon over HTTPS.
	IconURL string `json:"iconUrl,omitempty" validate:"required_without_all=KnownIcon MaterialIcon,excluded_with=KnownIcon MaterialIcon,omitempty,http_url"`
	// MaterialIcon configures a Google Material icon.
	MaterialIcon *MaterialIcon `json:"materialIcon,omitempty" validate:"required_without_all=KnownIcon IconURL,excluded_with=KnownIcon IconURL,omitempty"`
	// AltText describes the icon or its action for accessibility.
	AltText string `json:"altText,omitempty"`
	// ImageType selects a square or circular crop.
	ImageType ImageType `json:"imageType,omitempty" validate:"omitempty,oneof=SQUARE CIRCLE"`
}

func (i Icon) MarshalJSON() ([]byte, error) {
	count := 0
	if i.KnownIcon != "" {
		count++
	}
	if i.IconURL != "" {
		count++
	}
	if i.MaterialIcon != nil {
		count++
	}
	if count != 1 {
		return nil, fmt.Errorf("googlechat: Icon requires exactly one source")
	}
	type Embed Icon
	return json.Marshal(Embed(i))
}

// MaterialIcon selects a Material icon and customizes its rendering.
//
// Reference: https://developers.google.com/workspace/chat/api/reference/rest/v1/cards#MaterialIcon
type MaterialIcon struct {
	// Name identifies the Material icon.
	Name string `json:"name" validate:"required"`
	// Fill switches from an outlined icon to a filled icon.
	Fill bool `json:"fill,omitempty"`
	// Weight selects 100, 200, 300, 400, 500, 600, or 700; omission uses 400.
	Weight int `json:"weight,omitempty" validate:"omitempty,oneof=100 200 300 400 500 600 700"`
	// Grade adjusts emphasis from -25 to 200; omission uses 0.
	Grade int `json:"grade,omitempty" validate:"oneof=-25 0 200"`
}
