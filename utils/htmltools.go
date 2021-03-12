package utils

import (
	"html"
	"strings"

	"github.com/anacrolix/torrent"
)

// MakeNiceList - Generate html for torrent metainfo
func MakeNiceList(T *torrent.Torrent) string {
	files := T.Files()
	toret := "<!DOCTYPE html>"
	toret += `<html>
	<head>
		<title>` + html.EscapeString(T.Name()) + `</title>
	</head>
	<body>
`
	for _, file := range files {
		torrentpath := strings.Replace(file.Path(), T.Name(), "", -1)
		toret += "<a href=\"" + html.EscapeString(torrentpath) + "\">" + html.EscapeString(torrentpath) + "</a><br />\n"
	}
	toret += "</body>\n</html>"
	return toret
}
