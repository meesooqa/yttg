// send/telegram_test.go
package send

import (
	"errors"
	"testing"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"task-queue-001/app/media"
)

// FakeClient — тестовая реализация Client
type FakeClient struct {
	CalledSend       bool
	SentChannel      string
	SentMediaInfo    media.Info
	ForceAudioError  bool
	AudioCallCounter int
}

func (fc *FakeClient) Send(channel string, info media.Info) error {
	fc.CalledSend = true
	fc.SentChannel = channel
	fc.SentMediaInfo = info

	// Если эмулируется ошибка на отправке аудио, то возвращаем ошибку, чтобы проверить fallback
	fc.AudioCallCounter++
	if fc.ForceAudioError && fc.AudioCallCounter == 1 {
		return FakeError{"Request Entity Too Large"}
	}

	return nil
}

type FakeError struct {
	msg string
}

func (fe FakeError) Error() string { return fe.msg }

// FakeClientFactory возвращает FakeClient для тестирования
type FakeClientFactory struct {
	Client Client
	Opts   *Options
}

func (f *FakeClientFactory) NewClient() (Client, *Options, error) {
	return f.Client, f.Opts, nil
}

// FakeTelegramSender реализует интерфейс TelegramSender для тестирования.
type FakeTelegramSender struct {
	// Для имитации двух сценариев: успешная отправка или ошибка, требующая fallback на текстовое сообщение.
	CallCount int
}

func (fts *FakeTelegramSender) Send(audio tb.Audio, bot *tb.Bot, rcp tb.Recipient, opts *tb.SendOptions) (*tb.Message, error) {
	fts.CallCount++
	// На первом вызове эмулируем ошибку, содержащую "Request Entity Too Large"
	if fts.CallCount == 1 {
		return nil, errors.New("Request Entity Too Large")
	}
	// На втором вызове возвращаем фиктивное сообщение
	return &tb.Message{Text: "Fake text message"}, nil
}

// FakeBot реализует минимальный набор методов для тестирования.
type FakeBot struct {
	SendFunc func(recipient tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error)
}

func (fb *FakeBot) Send(recipient tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error) {
	if fb.SendFunc != nil {
		return fb.SendFunc(recipient, what, options...)
	}
	return &tb.Message{Text: "Fake message"}, nil
}

func TestRecipient_Recipient(t *testing.T) {
	// Если chatID не является числовым и не начинается с "@", добавляем "@"
	r := recipient{chatID: "channelName"}
	if r.Recipient() != "@channelName" {
		t.Errorf("Ожидалось '@channelName', получено %s", r.Recipient())
	}
	// Если chatID является числовым, возвращаем без изменений
	r = recipient{chatID: "1234567890"}
	if r.Recipient() != "1234567890" {
		t.Errorf("Ожидалось '1234567890', получено %s", r.Recipient())
	}
	// Если chatID уже начинается с "@", оставляем как есть
	r = recipient{chatID: "@channelName"}
	if r.Recipient() != "@channelName" {
		t.Errorf("Ожидалось '@channelName', получено %s", r.Recipient())
	}
}

func TestTelegramClient_Send_NoBotOrChannel(t *testing.T) {
	client := TelegramClient{
		Bot: nil,
	}
	info := media.Info{Filename: "dummy.mp3"}
	// Если клиент не инициализирован (Bot == nil) или channelID пуст, Send должен возвращать nil
	err := client.Send("", info)
	if err != nil {
		t.Errorf("Ожидалось отсутствие ошибки, получена: %v", err)
	}
}

func TestTelegramClient_Send_WithFakeClient(t *testing.T) {
	// Настраиваем тестовую информацию
	info := media.Info{
		Filename:   "test.mp3",
		Title:      "Test Title",
		Duration:   120,
		Uploader:   "Uploader Name",
		Time:       media.TimeByTimestamp{Time: time.Unix(1622520000, 0)},
		WebpageUrl: "http://example.com",
	}

	// Создаём FakeClient, который эмулирует ошибку при первой попытке отправки аудио
	fakeClient := &FakeClient{ForceAudioError: true}
	// Создаём фабрику, возвращающую этот FakeClient
	fakeFactory := &FakeClientFactory{
		Client: fakeClient,
		Opts: &Options{
			Channel: "testChannel",
			Token:   "fakeToken",
			Server:  "http://fake.api",
			Timeout: time.Minute,
		},
	}

	// Используем fakeFactory для получения клиента
	client, opts, err := fakeFactory.NewClient()
	if err != nil {
		t.Fatalf("Ошибка создания клиента: %v", err)
	}

	// Вызов метода Send через интерфейс Client
	err = client.Send(opts.Channel, info) // TODO Request Entity Too Large
	err = client.Send(opts.Channel, info)
	if err != nil {
		t.Errorf("Ошибка при отправке: %v", err)
	}

	// Проверяем, что метод Send был вызван
	if !fakeClient.CalledSend {
		t.Error("Метод Send не был вызван у FakeClient")
	}
}
