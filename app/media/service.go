package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Service interface {
	GetInfo(ctx context.Context, link string) (*Info, error)
	Download(ctx context.Context, link string, id string) (string, error)
}

type MediaService struct{}

func NewMediaService() *MediaService {
	return &MediaService{}
}

func (s *MediaService) GetInfo(ctx context.Context, link string) (*Info, error) {
	if !isValidURL(link) {
		return nil, fmt.Errorf("URL is not valid")
	}

	infoFile, err := os.CreateTemp(os.TempDir(), "yttg-info-*.json")
	if err != nil {
		return nil, err
	}
	defer os.Remove(infoFile.Name())

	cmdInfo := CmdInfo(link)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdInfo)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = infoFile
	log.Printf("[DEBUG] executing command: %s", cmdInfo)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var info Info
	infoData, err := os.ReadFile(infoFile.Name())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(infoData, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func (s *MediaService) Download(ctx context.Context, link string, id string) (string, error) {
	basename := "var/yttg/" + id
	audioFormat := "mp3"

	cmdDownload := CmdDownload(link, basename, audioFormat)
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdDownload)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("[DEBUG] executing command: %s", cmdDownload)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	filename := basename + "." + audioFormat
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return filename, errors.New("file is not downloaded")
	}
	return filename, nil
}
