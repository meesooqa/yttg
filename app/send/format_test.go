package send

import (
	"strings"
	"testing"
	"time"

	"github.com/meesooqa/yttg/app/media"
)

func TestTelegramFormatter_Format_WithWebpageUrl(t *testing.T) {
	formatter := TelegramFormatter{}
	info := media.Info{
		Filename:   "test.mp3",
		Id:         "123",
		Title:      "Test Title",
		Duration:   120,
		Uploader:   "Uploader Name",
		Time:       media.TimeByTimestamp{Time: time.Unix(1622520000, 0)},
		WebpageUrl: "http://example.com",
		Thumbnail:  "http://example.com/thumb.jpg",
	}
	result := formatter.Format(info)

	// Проверяем наличие HTML-ссылки
	if !strings.Contains(result, `<a href="http://example.com">`) {
		t.Errorf("Ожидалась HTML-ссылка в результате, получено: %s", result)
	}
	// Проверяем, что мета-информация сформирована корректно (например, наличие имени автора без спецсимволов)
	if !strings.Contains(result, "#UploaderName") {
		t.Errorf("Ожидалась мета-информация с '#UploaderName', получено: %s", result)
	}
}

func TestTelegramFormatter_Format_WithoutWebpageUrl(t *testing.T) {
	formatter := TelegramFormatter{}
	info := media.Info{
		Filename:   "test.mp3",
		Id:         "123",
		Title:      "Test Title",
		Duration:   120,
		Uploader:   "Uploader Name",
		Time:       media.TimeByTimestamp{Time: time.Unix(1622520000, 0)},
		WebpageUrl: "",
		Thumbnail:  "http://example.com/thumb.jpg",
	}
	result := formatter.Format(info)

	// Если URL отсутствует, не должно быть ссылки (<a href=)
	if strings.Contains(result, "<a href=") {
		t.Errorf("Не ожидалась HTML-ссылка в результате, получено: %s", result)
	}
}
