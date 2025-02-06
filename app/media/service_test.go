package media

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestMediaService_GetInfo_InvalidURL(t *testing.T) {
	svc := NewMediaService()
	// Передаём заведомо недействительный URL – ожидаем, что isValidURL вернёт false
	_, err := svc.GetInfo(context.Background(), "invalid")
	if err == nil {
		t.Error("Ожидался error для недействительного URL, получен nil")
	}
	// Ошибка формируется непосредственно в функции GetInfo
	if err.Error() != "URL is not valid" {
		t.Errorf("Ожидалась ошибка 'URL is not valid', получена: %v", err)
	}
}

func TestMediaService_Download_InvalidLink(t *testing.T) {
	// Этот тест является интеграционным и требует наличия yt-dlp в PATH.
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp не найден в PATH, пропускаем тест")
	}

	svc := NewMediaService()
	filename, err := svc.Download(context.Background(), "invalid", "testid")
	if err == nil {
		t.Error("Ожидалась ошибка для недействительной ссылки, получен nil")
	}
	errMsg := err.Error()
	if errMsg != "file is not downloaded" && !strings.Contains(errMsg, "exit status") {
		t.Errorf("Ожидалась ошибка 'file is not downloaded' или содержащая 'exit status', получена: %v", err)
	}
	if filename != "" {
		os.Remove(filename)
	}
}

func TestMediaService_GetInfo_Timeout(t *testing.T) {
	// Этот тест иллюстрирует использование контекста с таймаутом.
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp не найден в PATH, пропускаем тест")
	}
	svc := NewMediaService()
	// Устанавливаем очень короткий таймаут, чтобы команда не успевала выполниться.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	_, err := svc.GetInfo(ctx, "http://example.com")
	if err == nil {
		t.Error("Ожидался error из-за таймаута, получен nil")
	}
}
