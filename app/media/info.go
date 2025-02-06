package media

import (
	"encoding/json"
	"time"
)

// TimeByTimestamp is a custom type to handle date format "2006-01-02T15:04:05"
type TimeByTimestamp struct {
	time.Time
}

// UnmarshalJSON is a custom unmarshaler for TimeByTimestamp
func (t *TimeByTimestamp) UnmarshalJSON(b []byte) (err error) {
	var timestamp int64
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}

	t.Time = time.Unix(timestamp, 0)

	return nil
}

// Info represents `yt-dlp --dump-json` output
type Info struct {
	Filename   string
	Id         string          `json:"id"`
	Title      string          `json:"title"`
	Duration   int             `json:"duration"`
	Uploader   string          `json:"uploader"`
	Time       TimeByTimestamp `json:"timestamp"`
	WebpageUrl string          `json:"webpage_url"`
	Thumbnail  string          `json:"thumbnail"`
}
