# Player animation pack

Every sheet in this folder uses the same layout:

- Frame size: 120 × 120 pixels
- Grid: 4 columns × 8 rows
- Columns: sequential animation frames
- Rows: south, southeast, east, northeast, north, northwest, west, southwest

animations.json is the source of truth for sheet paths, frame timing, looping,
and grid dimensions. Keep gameplay rules in Go, but update this manifest when
you replace art or change the frame layout.

Available clips: idle, walk, sprint, gun, sword, fist, hurt, heal, and death.
