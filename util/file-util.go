package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"
)

var KnowExtension = [4]string{".mp3", ".wav", ".flac", ".ogg"}

func GetKnowFiles(dir string) []string {
	var knowFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path != dir {
			if contains(KnowExtension, filepath.Ext(path)) {
				knowFiles = append(knowFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}

	return knowFiles
}

func GetKnowFilesInfo(dir string) []os.FileInfo {
	var knowFiles []os.FileInfo
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range files {
		if v.IsDir() {
			f := GetKnowFiles(filepath.Join(dir, v.Name()))
			if len(f) != 0 {
				knowFiles = append(knowFiles, v)
			}

		} else {
			if contains(KnowExtension, filepath.Ext(v.Name())) {
				knowFiles = append(knowFiles, v)
			}
		}
	}
	return knowFiles
}

func contains(s [4]string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func FormatProgDur(d time.Duration) string {
	// thanks to https://github.com/Depado
	h := math.Mod(d.Hours(), 24)
	m := math.Mod(d.Minutes(), 60)
	s := math.Mod(d.Seconds(), 60)
	tot := fmt.Sprintf("%02d:%02d:%02d", int(h), int(m), int(s))
	return tot
}
