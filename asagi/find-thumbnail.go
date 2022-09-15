package asagi

import "strings"

//Finds the thumbnail on the folder
func findThumbnail(imageFolder string, board string, previewOrig string) string {
	var b strings.Builder
	b.Grow(len(imageFolder) + 1 + len(board) + 1 + 7 + 4 + 1 + 2 + 1 + len(previewOrig))

	b.WriteString(imageFolder)
	b.WriteString("/")
	b.WriteString(board)
	b.WriteString("/thumb/")
	b.WriteString(previewOrig[:4])
	b.WriteString("/")
	b.WriteString(previewOrig[4:6])
	b.WriteString("/")
	b.WriteString(previewOrig)

	return b.String()
}
