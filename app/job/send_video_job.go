package job

import (
	"context"
	"fmt"
	"strings"

	"github.com/meesooqa/yttg/app/media"
	"github.com/meesooqa/yttg/app/send"
)

// BrowserForCookies is type to represent browser code
type BrowserForCookies string

const (
	// Chrome --cookies-from-browser chrome
	Chrome BrowserForCookies = "chrome"
)

// ParseCookiesFromBrowser converts a string to CookiesFromBrowser enum value
// Returns error if the input string doesn't match any supported browser
func ParseCookiesFromBrowser(value string) (BrowserForCookies, error) {
	browser := BrowserForCookies(value)
	switch browser {
	case Chrome:
		return browser, nil
	default:
		return "", fmt.Errorf("unsupported browser: %s. Valid values: %s",
			value, strings.Join([]string{string(Chrome)}, ", "))
	}
}

// SendVideoJob download audio and send to Telegram
type SendVideoJob struct {
	BaseJob
	URL                string
	CookiesFromBrowser BrowserForCookies
	MediaService       media.Service
	TelegramFactory    send.ClientFactory
}

// Execute implements SendVideoJob
func (j SendVideoJob) Execute() error {
	fmt.Printf("Start processing URL: %s\n", j.URL)

	link := j.URL
	ctx := context.Background()

	ms := j.MediaService
	info, err := ms.GetInfo(ctx, link)
	if err != nil {
		return fmt.Errorf("failed to get info: %v", err)
	}

	filename, err := ms.Download(ctx, link, j.ID)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	info.Filename = filename

	tgClient, opts, err := j.TelegramFactory.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create Telegram client: %v", err)
	}

	if err := tgClient.Send(opts.Channel, *info); err != nil {
		return fmt.Errorf("failed to send to Telegram: %v", err)
	}

	return nil
}
