# 📍 Mappyboi

![Heatmap example](https://user-images.githubusercontent.com/22618981/101420497-36556180-38ea-11eb-9417-d25dda5ae421.png)

Generates a heatmap of where you have been, using data from Google Takeout and Apple Health Export (or other assorted GPX files).

Example usage:

```bash
go run cmd/mappyboi/main.go --location_history="/path/to/Location History.json" --gpx_folder="/path/to/workout-routes"
open map.html
```
