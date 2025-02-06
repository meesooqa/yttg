package job

import (
	"context"

	"task-queue-001/app/media"
	"task-queue-001/app/send"
)

// FakeMediaService реализует интерфейс media.Service.
type FakeMediaService struct {
	// Функции, которые можно настроить в тесте
	GetInfoFunc  func(ctx context.Context, link string) (*media.Info, error)
	DownloadFunc func(ctx context.Context, link, id string) (string, error)
}

func (ms *FakeMediaService) GetInfo(ctx context.Context, link string) (*media.Info, error) {
	return ms.GetInfoFunc(ctx, link)
}

func (ms *FakeMediaService) Download(ctx context.Context, link, id string) (string, error) {
	return ms.DownloadFunc(ctx, link, id)
}

// FakeTgClient — fake‑реализация send.Client.
type FakeTgClient struct {
	// Функция для настройки поведения метода Send.
	SendFunc func(channel string, info media.Info) error
}

func (c *FakeTgClient) Send(channel string, info media.Info) error {
	return c.SendFunc(channel, info)
}

// FakeTelegramFactory реализует send.ClientFactory.
type FakeTelegramFactory struct {
	// Фиктивный клиент, который будет возвращаться
	Client send.Client
	// Опции, которые будут возвращаться
	Opts *send.Options
	// Функция для имитации ошибки создания клиента (если необходимо)
	NewClientErr error
}

func (f *FakeTelegramFactory) NewClient() (send.Client, *send.Options, error) {
	if f.NewClientErr != nil {
		return nil, nil, f.NewClientErr
	}
	return f.Client, f.Opts, nil
}
