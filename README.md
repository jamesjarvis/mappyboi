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

```txt
USAGE:
   mappyboi v2 [global options] command [command options] [arguments...]

AUTHOR:
   James Jarvis <git@jamesjarvis.io>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --base_file FILE, --base FILE                        Base location history append only FILE, in .json or .json.gz (by suffixing with .gz mappyboi will compress the resulting file, using significantly less storage)
   --google_location_history FILE, --glh FILE           Google Takeout Location History FILE
   --gpx_directory DIRECTORY, --gpxd DIRECTORY          GPX DIRECTORY to load .gpx files from
   --fit_directory DIRECTORY, --fitd DIRECTORY          FIT DIRECTORY to load .fit files from
   --polarstep_directory DIRECTORY, --pstepd DIRECTORY  Polarstep DIRECTORY to load locations.json files from
   --output_type value, --ot value                      Output format, must be one of [ MAP ] (default: "MAP")
   --output_file FILE, --of FILE                        Output FILE to write to
   --output_reduce_points value, --rp value             If you struggle to open the file in a browser due to too many points, reduce the number of points by increasing this value. (default: 0)
   --output_randomise_points, --rand                    If you want to export the view of the points, but otherwise randomise the data to prevent perfect tracking, this will randomise the order. (default: false)
   --output_filter_start_date value, --from value       To filter the output to only include points on or after the provided timestamp
   --output_filter_end_date value, --to value           To filter the output to only include points on or before the provided timestamp
   --version, -v                                        print the version
   --help, -h                                           show help
```

```bash
mappyboi \
  --base_file all_locations.json \
  --google_location_history="/path/to/Location History.json" \
  --gpx_directory="/path/to/workout-routes" \
  --fit_directory="/path/to/workout-routes" \
  --polarstep_directory="/path/to/polarsteps" \
  --output_reduce_points 10 --output_randomise_points \
  --output_filter_start_date="2024-01-01T00:00:01Z"
  --output_filter_end_date="2024-03-01T00:00:01Z"
  --output_type MAP \
  --output_file heatmap.html
open heatmap.html
```
