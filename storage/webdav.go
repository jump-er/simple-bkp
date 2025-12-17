package storage

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

type Wd struct {
	url            string
	user           string
	password       string
	RootDir        string
	LocalFilePath  string
	WebdavFilePath string
	Client         *gowebdav.Client
}

func (w *Wd) ClientInit() error {
	c := gowebdav.NewClient(
		w.url,
		w.user,
		w.password,
	)
	if err := c.Connect(); err != nil {
		return err
	}
	logrus.Info("WebDav auth successfull")

	w.Client = c

	return nil
}

func (w *Wd) CreateRoot() error {
	if _, err := w.Client.Stat(w.RootDir); err != nil {
		err = w.Client.Mkdir(w.RootDir, 0644)
		if err != nil {
			return err
		}
		logrus.Infof("WebDav root dir %s created successfull", w.RootDir)
	}
	logrus.Infof("WebDav root dir %s already exist", w.RootDir)

	return nil
}

func (w *Wd) Upload() error {
	b, err := os.ReadFile(w.LocalFilePath)
	if err != nil {
		return err
	}

	err = w.Client.Write(w.WebdavFilePath, b, 0644)
	if err != nil {
		return err
	}

	logrus.Info("WebDav file uploaded successfull")
	return nil
}

func (w *Wd) IsRootDirExists() error {
	w.Client.Stat(w.RootDir)
	return nil
}

func (w *Wd) GetListFiles() ([]os.FileInfo, error) {
	l, err := w.Client.ReadDir(w.RootDir)
	if err != nil {
		return []os.FileInfo{}, err
	}

	return l, nil
}

func (w *Wd) Remove(f string) error {
	if err := w.Client.Remove(f); err != nil {
		return nil
	}
	return nil
}

func (w *Wd) GetRootDir() string {
	return w.RootDir
}

func NewWd() *Wd {
	return &Wd{
		url:            os.Getenv("WEB_DAV_URL"),
		user:           os.Getenv("WEB_DAV_USER"),
		password:       os.Getenv("WEB_DAV_PASSWORD"),
		RootDir:        os.Getenv("WEB_DAV_ROOT_DIR"),
		LocalFilePath:  "",
		WebdavFilePath: "",
	}
}
