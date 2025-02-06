package send

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/meesooqa/yttg/app/media"
)

type Options struct {
	Channel string
	Server  string
	Token   string
	Timeout time.Duration
}

type Client interface {
	Send(channel string, info media.Info) error
}

type ClientFactory interface {
	NewClient() (Client, *Options, error)
}

type EnvClientFactory struct{}

func (f *EnvClientFactory) NewClient() (Client, *Options, error) {
	return NewTelegramClientFromEnv()
}

// TelegramSender is the interface for sending messages to telegram
type TelegramSender interface {
	Send(tb.Audio, *tb.Bot, tb.Recipient, *tb.SendOptions) (*tb.Message, error)
}

type TelegramClient struct {
	Bot            *tb.Bot
	Timeout        time.Duration
	TelegramSender TelegramSender
	Formatter      TelegramFormatter
}

func optionsFromEnv() *Options {
	telegramTimeout, _ := strconv.Atoi(os.Getenv("TELEGRAM_TIMEOUT"))
	return &Options{
		Channel: os.Getenv("TELEGRAM_CHAN"),
		Server:  os.Getenv("TELEGRAM_SERVER"),
		Token:   os.Getenv("TELEGRAM_TOKEN"),
		Timeout: time.Duration(telegramTimeout) * time.Minute,
	}
}

// NewTelegramClientFromEnv init telegram client from ENV
func NewTelegramClientFromEnv() (client Client, opts *Options, err error) {
	opts = optionsFromEnv()
	client, err = newTelegramClient(
		opts.Token,
		opts.Server,
		opts.Timeout,
		&TelegramSenderImpl{},
		TelegramFormatter{},
	)
	return
}

// newTelegramClient init telegram client
func newTelegramClient(token, apiURL string, timeout time.Duration, tgs TelegramSender, tf TelegramFormatter) (Client, error) {
	log.Printf("[INFO] create telegram client for %s, timeout: %s", apiURL, timeout)
	if timeout == 0 {
		timeout = time.Second * 60
	}

	if token == "" {
		return TelegramClient{
			Bot:     nil,
			Timeout: timeout,
		}, nil
	}

	bot, err := tb.NewBot(tb.Settings{
		URL:    apiURL,
		Token:  token,
		Client: &http.Client{Timeout: timeout},
	})
	if err != nil {
		return TelegramClient{}, err
	}

	result := TelegramClient{
		Bot:            bot,
		Timeout:        timeout,
		TelegramSender: tgs,
		Formatter:      tf,
	}
	return result, err
}

func (client TelegramClient) Send(channelID string, item media.Info) (err error) {
	if client.Bot == nil || channelID == "" {
		return nil
	}

	message, err := client.sendAudio(channelID, item)
	if err != nil && strings.Contains(err.Error(), "Request Entity Too Large") {
		message, err = client.sendText(channelID, item)
	}

	if err != nil {
		return errors.Wrapf(err, "can't send to telegram for %+v", item.Filename)
	}

	log.Printf("[DEBUG] telegram message sent: \n%s", message.Text)
	//log.Printf("[DEBUG] telegram message sent: \n%s", message.Text, message.Caption)
	return nil
}

func (client TelegramClient) sendText(channelID string, item media.Info) (*tb.Message, error) {
	message, err := client.Bot.Send(
		recipient{chatID: channelID},
		client.Formatter.Format(item),
		tb.ModeHTML,
		tb.NoPreview,
	)

	return message, err
}

func (client TelegramClient) sendAudio(channelID string, item media.Info) (*tb.Message, error) {
	defer os.Remove(item.Filename)

	audio := tb.Audio{
		File:      tb.FromDisk(item.Filename),
		MIME:      "audio/mpeg",
		Caption:   client.getMessageHTML(item),
		Title:     item.Title,
		Performer: item.Uploader,
		Duration:  item.Duration,
	}

	return client.TelegramSender.Send(audio, client.Bot, recipient{chatID: channelID}, &tb.SendOptions{ParseMode: tb.ModeHTML})
}

// getMessageHTML generates HTML message from provided media.Info
func (client TelegramClient) getMessageHTML(item media.Info) string {
	return client.Formatter.Format(item)
}

type recipient struct {
	chatID string
}

func (r recipient) Recipient() string {
	if _, err := strconv.ParseInt(r.chatID, 10, 64); err != nil && !strings.HasPrefix(r.chatID, "@") {
		return "@" + r.chatID
	}

	return r.chatID
}

// TelegramSenderImpl is a TelegramSender implementation that sends messages to Telegram for real
type TelegramSenderImpl struct{}

// Send sends a message to Telegram
func (tg *TelegramSenderImpl) Send(audio tb.Audio, bot *tb.Bot, rcp tb.Recipient, opts *tb.SendOptions) (*tb.Message, error) {
	return audio.Send(bot, rcp, opts)
}
