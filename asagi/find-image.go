package asagi

import "strings"

//Finds the image in the folder
func findImage(imageFolder string, board string, mediaOrig string) string {
	var b strings.Builder
	b.Grow(len(imageFolder) + 1 + len(board) + 1 + 7 + 4 + 1 + 2 + 1 + len(mediaOrig))

	b.WriteString(imageFolder)
	b.WriteString("/")
	b.WriteString(board)
	b.WriteString("/image/")
	b.WriteString(mediaOrig[:4])
	b.WriteString("/")
	b.WriteString(mediaOrig[4:6])
	b.WriteString("/")
	b.WriteString(mediaOrig)

	return b.String()
}
