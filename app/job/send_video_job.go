package job

import (
	"context"
	"fmt"

	"github.com/meesooqa/yttg/app/media"
	"github.com/meesooqa/yttg/app/send"
)

// SendVideoJob download audio and send to Telegram
type SendVideoJob struct {
	BaseJob
	URL             string
	MediaService    media.Service
	TelegramFactory send.ClientFactory
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
