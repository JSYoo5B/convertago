package googlechat_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/JSYoo5B/convertago/googlechat"
)

func TestNestedLayoutsRejectMissingContent(t *testing.T) {
	var paragraph *googlechat.TextParagraph
	for _, widget := range []any{
		googlechat.ColumnWidget{},
		googlechat.ColumnWidget{Content: paragraph},
		googlechat.NestedWidget{},
		googlechat.NestedWidget{Content: paragraph},
	} {
		if _, err := json.Marshal(widget); err == nil {
			t.Fatal("expected missing nested content to fail")
		}
	}
}

func TestNestedLayoutStrings(t *testing.T) {
	column := googlechat.Columns{ColumnItems: []googlechat.Column{{Widgets: []googlechat.ColumnWidget{
		{Content: googlechat.TextParagraph{Text: "Status"}},
		{Content: googlechat.Image{ImageURL: "https://example.com/chart.png", AltText: "Chart"}},
	}}}}
	carousel := googlechat.Carousel{CarouselCards: []googlechat.CarouselCard{{
		Widgets:       []googlechat.NestedWidget{{Content: googlechat.TextParagraph{Text: "Status"}}},
		FooterWidgets: []googlechat.NestedWidget{{Content: googlechat.ButtonList{Buttons: []googlechat.Button{{Text: "Report"}}}}},
	}}}
	if got := column.String(); got != "Status\nChart" {
		t.Fatalf("column text = %q", got)
	}
	if got := carousel.String(); got != "Status\nReport" {
		t.Fatalf("carousel text = %q", got)
	}
}

func TestColumnWidgetHasNoAlignmentWrapper(t *testing.T) {
	payload, err := json.Marshal(googlechat.ColumnWidget{Content: googlechat.TextParagraph{Text: "Status"}})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), "horizontalAlignment") || strings.Contains(string(payload), "Content") {
		t.Fatalf("unexpected implementation fields: %s", payload)
	}
}
