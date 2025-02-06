package media

import (
	"strings"
	"testing"
)

func TestCmdInfo(t *testing.T) {
	link := "http://example.com/video"
	cmd := CmdInfo(link)
	expected := "yt-dlp -f bestaudio -x --embed-thumbnail --dump-json \"" + link + "\""
	if cmd != expected {
		t.Errorf("Ожидалась команда:\n%q\nПолучена:\n%q", expected, cmd)
	}
}

func TestCmdDownload(t *testing.T) {
	link := "http://example.com/video"
	fileName := "output"
	audioFormat := "mp3"
	cmd := CmdDownload(link, fileName, audioFormat)
	expected := "yt-dlp -f bestaudio -x --embed-thumbnail --audio-format=" + audioFormat +
		" --no-progress -o " + fileName + " -- \"" + link + "\""
	// Убираем возможные лишние пробелы
	if strings.TrimSpace(cmd) != expected {
		t.Errorf("Ожидалась команда:\n%q\nПолучена:\n%q", expected, cmd)
	}
}
