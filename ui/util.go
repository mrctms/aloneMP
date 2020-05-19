package ui

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func getKnowFiles(dir string) []os.FileInfo {
	var knowFiles []os.FileInfo
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range files {
		if v.IsDir() {
			f := getKnowFiles(filepath.Join(dir, v.Name()))
			if len(f) != 0 {
				knowFiles = append(knowFiles, v)
			}

		} else {
			if contains(knowExtension, filepath.Ext(v.Name())) {
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
