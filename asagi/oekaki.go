package asagi

import "strings"

//OekakiProperties are the possible EXIF field
//properties relating to Oekaki data.
var OekakiProperties = []string{
	"Time",
	"Painter",
	"Source",
}

//HasOekaki determines whether or not the post's
//exif field has any Oekaki data.
func HasOekaki(exif map[string]string) bool {
	for _, property := range OekakiProperties {
		_, hasProperty := exif[property]
		if hasProperty {
			return true
		}
	}

	return false
}

//CreateOekaki generates the HTML for the Oekaki data.
//This should then be appended to the end of the post's
//comment field.
func CreateOekaki(exif map[string]string) string {
	properties := make([]string, 0, 3)

	time, hasTime := exif["Time"]
	if hasTime {
		properties = append(properties, "Time: "+time)
	}

	painter, hasPainter := exif["Painter"]
	if hasPainter {
		properties = append(properties, "Painter: "+painter)
	}

	source, hasSource := exif["Source"]
	if hasSource {
		properties = append(properties, "Source: "+source)
	}

	return "<small><b>Oekaki Post</b>(" + strings.Join(properties, ", ") + ")</small>"
}
