package media

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeByTimestamp_UnmarshalJSON(t *testing.T) {
	// Пример timestamp в секундах (например, 1622520000 соответствует 2021-06-01T00:00:00 UTC)
	jsonInput := []byte("1622520000")
	var tbt TimeByTimestamp
	err := tbt.UnmarshalJSON(jsonInput)
	if err != nil {
		t.Fatalf("UnmarshalJSON вернул ошибку: %v", err)
	}
	expectedTime := time.Unix(1622520000, 0)
	if !tbt.Equal(expectedTime) {
		t.Errorf("Ожидалось время %v, получено %v", expectedTime, tbt.Time)
	}
}

func TestInfo_Unmarshal(t *testing.T) {
	jsonData := `{
		"id": "123",
		"title": "Test Video",
		"duration": 3600,
		"uploader": "TestUploader",
		"timestamp": 1622520000,
		"webpage_url": "http://example.com",
		"thumbnail": "http://example.com/thumb.jpg"
	}`
	var info Info
	err := json.Unmarshal([]byte(jsonData), &info)
	if err != nil {
		t.Fatalf("Unmarshal Info вернул ошибку: %v", err)
	}
	// Проверяем отдельные поля
	if info.Id != "123" {
		t.Errorf("Ожидалось id '123', получено %s", info.Id)
	}
	if info.Title != "Test Video" {
		t.Errorf("Ожидалось title 'Test Video', получено %s", info.Title)
	}
	// Можно добавить дополнительные проверки для остальных полей
}
