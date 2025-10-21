package cmd

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type src struct {
	srcDir                   string
	targetArchivePath        string
	targetArchiveNamePrefix  string
	targetArchiveFullName    string
	shouldRemoveLocalArchive bool
}

func (s *src) zipFiles() error {
	arhiveFileName := s.targetArchiveFullName
	zipFile, err := os.Create(arhiveFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	filepath.Walk(s.srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(s.srcDir, path)
		if err != nil {
			return err
		}

		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		zipEntryWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntryWriter, fileToZip)
		return err
	})

	return nil
}

func (s *src) isPathExists() bool {
	if _, err := os.Stat(s.srcDir); os.IsNotExist(err) {
		return false
	}

	return true
}

func (s *src) nowTime() string {
	return time.Now().Format("2006-01-02--15-04")
}

func (s *src) setTargetArchiveFullName() {
	s.targetArchiveFullName = s.targetArchiveNamePrefix + "-" + s.nowTime() + ".zip"
}

func (s *src) handleLocalArchive() error {
	if shouldRemoveLocalArchive {
		err := os.RemoveAll(s.targetArchivePath + "/" + s.targetArchiveFullName)
		if err != nil {
			return err
		}
		logrus.Warnf("Archive %s is removed", s.targetArchivePath+"/"+s.targetArchiveFullName)
	}

	return nil
}

func NewSrc() *src {
	return &src{
		srcDir:                   srcDir,
		targetArchivePath:        targetArchiveLocalPath,
		targetArchiveNamePrefix:  targetArchiveNamePrefix,
		shouldRemoveLocalArchive: shouldRemoveLocalArchive,
	}
}
