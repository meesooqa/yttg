package media

import (
	"bytes"
	"text/template"
)

const (
	tplInfo     = "yt-dlp -f bestaudio -x --embed-thumbnail --dump-json \"{{.Link}}\""
	tplDownload = "yt-dlp -f bestaudio -x --embed-thumbnail --audio-format={{.AudioFormat}} --no-progress -o {{.FileName}} -- \"{{.Link}}\""
)

// CmdInfo returns command to save media info from link
func CmdInfo(link string) string {
	tpl := tplInfo
	params := struct {
		Link string
	}{
		Link: link,
	}

	b := bytes.Buffer{}
	template.Must(template.New("ytdlp-info").Parse(tpl)).Execute(&b, params)

	return b.String()
}

// CmdDownload returns command to download media from link
func CmdDownload(link, fileName, audioFormat string) string {
	tpl := tplDownload
	params := struct {
		Link        string
		FileName    string
		AudioFormat string
	}{
		Link:        link,
		FileName:    fileName,
		AudioFormat: audioFormat,
	}

	b := bytes.Buffer{}
	template.Must(template.New("ytdlp-download").Parse(tpl)).Execute(&b, params)

	return b.String()
}
