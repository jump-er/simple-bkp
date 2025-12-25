package storage

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/studio-b12/gowebdav"
)

type Wd struct {
	url            string
	user           string
	password       string
	rootDir        string
	LocalFilePath  string
	WebdavFilePath string
	Client         *gowebdav.Client
}

func NewWd() *Wd {
	w := &Wd{
		url:      os.Getenv("WEB_DAV_URL"),
		user:     os.Getenv("WEB_DAV_USER"),
		password: os.Getenv("WEB_DAV_PASSWORD"),
		rootDir:  os.Getenv("WEB_DAV_ROOT_DIR"),
	}

	c := gowebdav.NewClient(
		w.url,
		w.user,
		w.password,
	)
	if err := c.Connect(); err != nil {
		logrus.Fatal(err)
	}
	w.Client = c

	c.SetTimeout(30 * time.Second)

	logrus.Infof("WebDav connect and auth to %s successfull", w.url)

	return w
}

func (w *Wd) CreateRoot() error {
	if _, err := w.Client.Stat(w.GetRootDir()); err != nil {
		err = w.Client.Mkdir(w.GetRootDir(), 0755)
		if err != nil {
			return err
		}
		logrus.Infof("WebDav root dir %s created successfull", w.GetRootDir())
	}
	logrus.Infof("WebDav root dir %s already exist", w.GetRootDir())

	return nil
}

func (w *Wd) Upload() error {
	localFilePath := w.LocalFilePath
	webdavFilePath := w.GetRootDir() + "/" + strings.Split(w.LocalFilePath, "/")[len(strings.Split(w.LocalFilePath, "/"))-1]

	file, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	md5s, sha256s, err := FileHashes(w.LocalFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	w.Client.SetHeader("Etag", md5s)
	w.Client.SetHeader("Sha256", sha256s)

	err = w.Client.WriteStreamWithLength(webdavFilePath, file, info.Size(), 0644)
	if err != nil {
		return err
	}
	logrus.Infof("WebDav file %s uploaded successfull", webdavFilePath)

	return nil
}

func (w *Wd) GetListFiles() ([]os.FileInfo, error) {
	l, err := w.Client.ReadDir(w.GetRootDir())
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
	return "/" + w.rootDir
}
