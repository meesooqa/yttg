package send

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"task-queue-001/app/media"
)

type TelegramFormatter struct{}

// Format generates HTML message from provided media.Info
func (f *TelegramFormatter) Format(item media.Info) string {
	var titleLink, metaInfo string

	re := regexp.MustCompile(`[^\p{L}\d_]+`)
	author := re.ReplaceAllString(item.Uploader, "")
	metaInfo = fmt.Sprintf("#%s <code>%s</code>", author, item.Time.Format(time.DateTime))

	title := strings.TrimSpace(item.Title)
	if title != "" && item.WebpageUrl == "" {
		titleLink = fmt.Sprintf("%s", title)
	} else if title != "" && item.WebpageUrl != "" {
		titleLink = fmt.Sprintf("<a href=%q>%s</a>", item.WebpageUrl, title)
	}

	return metaInfo + "\n" + titleLink
}
