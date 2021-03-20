package util

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Logger struct {
	filePath string
}

func NewLogger(fileName string) (*Logger, error) {
	logger := new(Logger)
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	rootDir := filepath.Join(u.HomeDir, ".aloneMP")
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		os.Mkdir(rootDir, 0755)
	}
	logger.filePath = filepath.Join(rootDir, fileName)
	return logger, nil
}

func (l *Logger) Write(content string) {
	f, err := os.OpenFile(l.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println(content)
}
