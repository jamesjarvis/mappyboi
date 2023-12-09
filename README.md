# 📍 Mappyboi

![Heatmap example](https://user-images.githubusercontent.com/22618981/101420497-36556180-38ea-11eb-9417-d25dda5ae421.png)

Generates a heatmap of where you have been, using data from:
- Google Takeout (Location History)
- Apple Health Export
- Strava (or other assorted .gpx files)
- Strava (or other assorted .fit files)
- Polarsteps Data Export


If you have multiple gpx or fit directories, it is recommended to run multiple times, changing the directory path but keeping the same base file.
Only new points will be added.

## Installation

```bash
go install github.com/jamesjarvis/mappyboi/v2
```

## Usage

```bash
mappyboi --base_file all_locations.json --google_location_history="/path/to/Location History.json" --gpx_directory="/path/to/workout-routes" --fit_directory="/path/to/workout-routes" --polarstep_directory="/path/to/polarsteps" --output_reduce_points 10 --output_randomise_points --output_type MAP --output_file heatmap.html
open heatmap.html
```
