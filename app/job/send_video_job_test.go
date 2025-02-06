package job

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"task-queue-001/app/media"
	"task-queue-001/app/send"
)

func TestSendVideoJob_Execute_Success(t *testing.T) {
	expectedURL := "http://example.com/video"
	expectedFilename := "video_123.mp3"

	// Настраиваем fakeMediaService: GetInfo и Download возвращают корректные данные
	fakeMediaService := &FakeMediaService{
		GetInfoFunc: func(ctx context.Context, link string) (*media.Info, error) {
			if link != expectedURL {
				return nil, errors.New("wrong URL")
			}
			return &media.Info{
				Id:         "123",
				Title:      "Test Video",
				Duration:   120,
				Uploader:   "Uploader",
				Time:       media.TimeByTimestamp{Time: time.Now()},
				WebpageUrl: "http://example.com",
			}, nil
		},
		DownloadFunc: func(ctx context.Context, link, id string) (string, error) {
			if id != "job1" {
				return "", errors.New("wrong job id")
			}
			return expectedFilename, nil
		},
	}

	// Настраиваем fakeTelegramClient: Send возвращает nil (успешно)
	fakeTgClient := &FakeTgClient{
		SendFunc: func(channel string, info media.Info) error {
			if channel != "testChannel" {
				return errors.New("wrong channel")
			}
			// Можно добавить проверки содержимого info, если нужно
			return nil
		},
	}

	fakeTgFactory := &FakeTelegramFactory{
		Client: fakeTgClient,
		Opts: &send.Options{
			Channel: "testChannel",
			Token:   "dummyToken",
			Server:  "http://dummy.api",
			Timeout: time.Minute,
		},
	}

	job := SendVideoJob{
		BaseJob:         BaseJob{ID: "job1", Status: StatusQueued},
		URL:             expectedURL,
		MediaService:    fakeMediaService,
		TelegramFactory: fakeTgFactory,
	}

	err := job.Execute()
	if err != nil {
		t.Errorf("Ожидалось успешное выполнение, получена ошибка: %v", err)
	}
}

func TestSendVideoJob_Execute_GetInfoError(t *testing.T) {
	fakeMediaService := &FakeMediaService{
		GetInfoFunc: func(ctx context.Context, link string) (*media.Info, error) {
			return nil, errors.New("failed to fetch info")
		},
		DownloadFunc: func(ctx context.Context, link, id string) (string, error) {
			return "irrelevant", nil
		},
	}

	// Фабрика и клиент не будут вызваны, но зададим их для полноты картины
	fakeTgFactory := &FakeTelegramFactory{
		Client: &FakeTgClient{
			SendFunc: func(channel string, info media.Info) error { return nil },
		},
		Opts: &send.Options{Channel: "testChannel"},
	}

	job := SendVideoJob{
		BaseJob:         BaseJob{ID: "job2", Status: StatusQueued},
		URL:             "http://example.com/failure",
		MediaService:    fakeMediaService,
		TelegramFactory: fakeTgFactory,
	}

	err := job.Execute()
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "failed to get info") {
		t.Errorf("Ожидалась ошибка на этапе GetInfo, получена: %v", err)
	}
}
