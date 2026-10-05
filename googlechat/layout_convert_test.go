package googlechat

import (
	"github.com/JSYoo5B/convertago/internal/conversion"
	"testing"
)

func TestColumnAndCarouselWidgetRestrictions(t *testing.T) {
	source := struct {
		Columns struct {
			Column struct {
				Divider struct{} `googlechat:"divider"`
			} `googlechat:"column"`
		} `googlechat:"columns"`
	}{}
	if _, err := ToMessage(source, conversion.WithWarningAsError()); err == nil {
		t.Fatal("column accepted divider")
	}
	source2 := struct {
		Carousel struct {
			Card struct {
				Text string `googlechat:"decoratedText"`
			} `googlechat:"carouselCard"`
		} `googlechat:"carousel"`
	}{}
	if _, err := ToMessage(source2, conversion.WithWarningAsError()); err == nil {
		t.Fatal("carousel accepted decoratedText")
	}
	source3 := struct {
		Columns struct {
			Column struct {
				Text struct {
					Text      string `googlechat:"part"`
					Alignment string `googlechat:"part;slot=horizontalAlignment"`
				} `googlechat:"textParagraph"`
			} `googlechat:"column"`
		} `googlechat:"columns"`
	}{}
	source3.Columns.Column.Text.Text = "text"
	source3.Columns.Column.Text.Alignment = "CENTER"
	if _, err := ToMessage(source3, conversion.WithWarningAsError()); err == nil {
		t.Fatal("nested widget silently dropped alignment")
	}
}

func TestGoogleColumnCount(t *testing.T) {
	type column struct {
		Text string `googlechat:"textParagraph"`
	}
	source := struct {
		Columns struct {
			Columns []column `googlechat:"column"`
		} `googlechat:"columns"`
	}{}
	source.Columns.Columns = []column{{"one"}, {"two"}}
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	source.Columns.Columns = append(source.Columns.Columns, column{"three"})
	if _, err := ToMessage(source, conversion.WithWarningAsError()); err == nil {
		t.Fatal("accepted three columns")
	}
}

func TestGoogleCustomCropValidation(t *testing.T) {
	type crop struct {
		Type  string  `googlechat:"part;slot=type"`
		Ratio float64 `googlechat:"part;slot=aspectRatio"`
	}
	type image struct {
		URL  string `googlechat:"part;slot=imageUri"`
		Crop crop   `googlechat:"imageCropStyle"`
	}
	type item struct {
		Image image `googlechat:"imageComponent"`
	}
	source := struct {
		Grid struct {
			Item item `googlechat:"gridItem"`
		} `googlechat:"grid"`
	}{}
	source.Grid.Item.Image = image{"https://example.com/image.png", crop{"RECTANGLE_CUSTOM", 1.5}}
	if _, err := ToMessage(source); err != nil {
		t.Fatal(err)
	}
	source.Grid.Item.Image.Crop.Ratio = 0
	if _, err := ToMessage(source, conversion.WithWarningAsError()); err == nil {
		t.Fatal("nonpositive custom aspect ratio accepted")
	}
	source.Grid.Item.Image.Crop = crop{"CIRCLE", 1}
	if _, err := ToMessage(source, conversion.WithWarningAsError()); err == nil {
		t.Fatal("aspect ratio for circle accepted")
	}
}

func TestEveryGoogleRoleAvailable(t *testing.T) {
	profile, _ := conversion.Lookup("googlechat")
	for name, role := range profile.Roles {
		if role.Unavailable {
			t.Errorf("unimplemented role %s", name)
		}
	}
}

func TestGridBorderColorAndClick(t *testing.T) {
	type color struct {
		Red   float64 `googlechat:"part;slot=red"`
		Green float64 `googlechat:"part;slot=green"`
		Blue  float64 `googlechat:"part;slot=blue"`
	}
	type border struct {
		Type   string `googlechat:"part;slot=type"`
		Color  color  `googlechat:"color"`
		Radius int    `googlechat:"part;slot=cornerRadius"`
	}
	type grid struct {
		Items  []string `googlechat:"gridItem"`
		Border border   `googlechat:"borderStyle"`
		Click  string   `googlechat:"openLink"`
	}
	source := struct {
		Grid grid `googlechat:"grid"`
	}{grid{[]string{"one", "two"}, border{"STROKE", color{0.1, 0.2, 0.3}, 4}, "https://example.com"}}
	message, err := ToMessage(source)
	if err != nil {
		t.Fatal(err)
	}
	native := message.CardsV2[0].Card.Sections[0].Widgets[0].Content.(Grid)
	if native.BorderStyle.StrokeColor.Green != 0.2 || native.BorderStyle.CornerRadius != 4 || native.OnClick.OpenLink.URL != "https://example.com" {
		t.Fatalf("grid=%#v", native)
	}
	source.Grid.Border.Type = "NO_BORDER"
	if _, err := ToMessage(source, conversion.WithWarningAsError()); err == nil {
		t.Fatal("stroke color for no border accepted")
	}
}
