package asagi

import (
	"fmt"
	"strings"
)

//CameraSpecificProperties are the EXIF properties
//specific to the Camera.
//This list is most likely incomplete.
var cameraSpecificProperties = []string{
	"Camera Model",
	"Photographer",
	"Camera Software",
	"Equipment Make",
	"Lens Size",
	"Firmware Version",
	"Serial Number",
	"Lens Name",
	"Focal Length (35mm Equiv)",
	"Maximum Lens Aperture",
	"Sensing Method",
	"Color Filter Array Pattern",
	"Unique Image ID",
	"Image Width",
	"Image Height",
	"Compression Scheme",
	"Pixel Composition",
}

//ImageSpecificProperties are the EXIF
//properties specific to the Image.
//This list is most likely incomplete.
var imageSpecificProperties = []string{
	"Exposure Time",
	"Image Orientation",
	"Horizontal Resolution",
	"Vertical Resolution",
	"Image Created",
	"Image Data Arrangement",
	"F-Number",
	"Focal Length",
	"Color Space Information",
	"Flash",
	"Light Source",
	"ISO Speed Rating",
	"Metering Mode",
	"Exposure Bias",
	"Exposure Program",
	"Aperture Priority",
	"Lens Aperture",
	"Brightness",
	"Rendering",
	"Exposure Mode",
	"Focus Type",
	"White Balance",
	"Scene Capture Type",
	"Gain Control",
	"Contrast",
	"Shooting Mode",
	"Image Size",
	"Focus Mode",
	"Drive Mode",
	"Flash Mode",
	"Compression Setting",
	"Macro Mode",
	"Exposure Compensation",
	"Color Matrix",
	"Saturation",
	"Sharpness",
	"Subject Distance Range",
}

//hasExif checks if the Asagi exif field has any actual EXIF data in it.
func hasExif(exif map[string]interface{}) bool {
	for _, property := range cameraSpecificProperties {
		_, exists := exif[property]
		if exists {
			return true
		}
	}

	for _, property := range imageSpecificProperties {
		_, exists := exif[property]
		if exists {
			return true
		}
	}

	return false
}

//createExifTable generates the EXIF HTML table given the EXIF field
//from Asagi and the post number.
func createExifTable(exif map[string]interface{}, postNumber uint) string {
	var b strings.Builder

	fmt.Fprintf(&b, "<br><br><span class=\"abbr\">[EXIF data available. Click <a href=\"javascript:void(0)\" onclick=\"toggle('exif%d')\">here</a> to show/hide.]</span><br><table class=\"exif\" id=\"exif%d\"><tr><td colspan=\"2\"><b>Camera-Specific Properties:</b></td></tr>", postNumber, postNumber)

	for _, property := range cameraSpecificProperties {
		value, exists := exif[property]
		if exists {
			switch value.(type) {
			case string:
				{
					b.WriteString("<tr><td>")
					b.WriteString(property)
					b.WriteString("</td><td>")
					b.WriteString(value.(string))
					b.WriteString("</td></tr>")
				}
			}
		}
	}

	b.WriteString("<tr><td colspan=\"2\"><b>Image-Specific Properties:</b></td></tr>")

	for _, property := range imageSpecificProperties {
		value, exists := exif[property]
		if exists {
			switch value.(type) {
			case string:
				{
					b.WriteString("<tr><td>")
					b.WriteString(property)
					b.WriteString("</td><td>")
					b.WriteString(value.(string))
					b.WriteString("</td></tr>")
				}
			}
		}
	}

	b.WriteString("</table>")

	return b.String()
}
