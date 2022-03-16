package maptemplate

import (
	"bufio"
	"fmt"
	"os"
	"text/template"

	_ "embed"

	"github.com/jamesjarvis/mappyboi/pkg/models"
)

//go:embed base_page.html
var basePage string

const maptemplate = `
<script>    
    
    var mappyboi = L.map(
        "mappyboi",
        {
            center: [48.0673, 12.8633],
            crs: L.CRS.EPSG3857,
            zoom: 6,
            zoomControl: true,
            preferCanvas: false,
        }
    );

    var tile_layer = L.tileLayer(
        "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
        {"attribution": "Data by \u0026copy; \u003ca href=\"http://openstreetmap.org\"\u003eOpenStreetMap\u003c/a\u003e, under \u003ca href=\"http://www.openstreetmap.org/copyright\"\u003eODbL\u003c/a\u003e.", "detectRetina": false, "maxNativeZoom": 18, "maxZoom": 18, "minZoom": 0, "noWrap": false, "opacity": 1, "subdomains": "abc", "tms": false}
    ).addTo(mappyboi);


    var heat_map = L.heatLayer(
        [ {{ range .GoLocations }}[{{ .Latitude }},{{ .Longitude }},1],{{ end }} ],
        {"blur": 4, "max": 11560, "maxZoom": 4, "minOpacity": 0.2, "radius": 7}
    ).addTo(mappyboi);
        
</script>
`

func GenerateHTML(filepath string, data *models.Data) error {
	tmpl, err := template.New("test").Parse(maptemplate)
	if err != nil {
		return fmt.Errorf("failed to parse html template: %w", err)
	}

	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %w", filepath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	_, err = w.WriteString(basePage)
	if err != nil {
		return fmt.Errorf("failed to write leaflet.js data to output file: %w", err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return fmt.Errorf("failed to execute html template: %w", err)
	}
	return nil
}
