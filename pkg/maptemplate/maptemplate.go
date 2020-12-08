package maptemplate

import (
	"bufio"
	"os"
	"text/template"

	"github.com/jamesjarvis/mappyboi/pkg/models"
)

const maptemplate = `<!DOCTYPE html>
<head>    
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
    
    <script>
        L_NO_TOUCH = false;
        L_DISABLE_3D = false;
    </script>
    
    <script src="https://cdn.jsdelivr.net/npm/leaflet@1.5.1/dist/leaflet.js"></script>
    <script src="https://code.jquery.com/jquery-1.12.4.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.2.0/js/bootstrap.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Leaflet.awesome-markers/2.0.2/leaflet.awesome-markers.js"></script>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/leaflet@1.5.1/dist/leaflet.css"/>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap.min.css"/>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.2.0/css/bootstrap-theme.min.css"/>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.6.3/css/font-awesome.min.css"/>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/Leaflet.awesome-markers/2.0.2/leaflet.awesome-markers.css"/>
    <link rel="stylesheet" href="https://rawcdn.githack.com/python-visualization/folium/master/folium/templates/leaflet.awesome.rotate.css"/>
    
    <style>
        html, body {
            width: 100%;
            height: 100%;
            margin: 0;
            padding: 0;
        }

        #map {
            position:absolute;
            top:0;
            bottom:0;
            right:0;
            left:0;
        }

        #mappyboi {
            position: relative;
            width: 100.0%;
            height: 100.0%;
            left: 0.0%;
            top: 0.0%;
        }
    </style>
        
    <script src="https://leaflet.github.io/Leaflet.heat/dist/leaflet-heat.js"></script>
</head>
<body>    
    <div class="folium-map" id="mappyboi" ></div>
</body>
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
		return err
	}

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
