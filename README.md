# 📍 Mappyboi

![Heatmap example](https://user-images.githubusercontent.com/22618981/101420497-36556180-38ea-11eb-9417-d25dda5ae421.png)

Generates a heatmap of where you have been, using data from Google Takeout and Apple Health Export / Strava (or other assorted GPX files).
You can specify multiple folders of gpx files (for example if you have a strava folder and an apple health folder).

## Installation

```bash
go install github.com/jamesjarvis/mappyboi/v2
```

## Usage

```bash
mappyboi --base_file all_locations.json --google_location_history="/path/to/Location History.json" --gpx_directory="/path/to/workout-routes" --output_reduce_points 4.5 --output_type MAP --output_file heatmap.html
open heatmap.html
```
