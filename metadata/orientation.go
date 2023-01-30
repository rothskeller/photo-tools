package metadata

// Orientation is an enumerated type indicating the orientation of an image.
type Orientation uint

// Values for Orientation.  The names reflect what has to be done to the data
// to get to the state where (0,0) is the top left corner.  The values match
// those of the Orientation tag in EXIF.
const (
	Rotate0        Orientation = 1
	FlipX          Orientation = 2
	FlipXY         Orientation = 3
	FlipY          Orientation = 4
	Rotate90FlipX  Orientation = 5
	Rotate90       Orientation = 6
	Rotate270FlipX Orientation = 7
	Rotate270      Orientation = 8
)
