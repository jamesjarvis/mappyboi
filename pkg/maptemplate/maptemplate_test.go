package maptemplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesjarvis/mappyboi/v2/pkg/types"
)

func TestGenerateHTMLEmbedsPointsAsJSONString(t *testing.T) {
	out := filepath.Join(t.TempDir(), "map.html")
	history := types.LocationHistory{Data: []types.Location{
		{Latitude: 51.4528557, Longitude: -0.2991395},
		{Latitude: -8.5, Longitude: 115.25},
	}}

	if err := GenerateHTML(out, history); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	html := string(raw)

	want := `L.heatLayer(
        JSON.parse("[[51.4528557,-0.2991395],[-8.5,115.25]]"),`
	if !strings.Contains(html, want) {
		t.Errorf("points should be embedded as a JSON string for JSON.parse, not a JS array literal; want %q in output", want)
	}
	if strings.Contains(html, "[51.4528557,-0.2991395,1]") {
		t.Error("output still contains the array literal form")
	}
}
