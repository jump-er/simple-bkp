package cmd

import (
	"fmt"
	"os"
	"simple-bkp/storage"

	"github.com/sirupsen/logrus"
)

type RemoteStorage interface {
	CreateRoot() error
	Upload() error
	GetListFiles() ([]os.FileInfo, error)
}

var remoteStorage RemoteStorage

func selectRemoteStorage(src *src) (RemoteStorage, error) {
	switch remoteStorageType {
	case "webdav":
		rs := storage.NewWd()
		err := rs.ClientInit()
		if err != nil {
			logrus.Warnf("WebDav auth error %s", err)
		}
		rs.LocalFilePath = src.targetArchivePath + "/" + src.targetArchiveFullName
		rs.WebdavFilePath = rs.RootDir + "/" + src.targetArchiveFullName

		return rs, nil

	case "s3":
		fmt.Println("TODO s3")

	case "smb":
		fmt.Println("TODO smb")
	}

	return nil, nil
}

func managerRemoteStorage(src *src, cmd string) {
	var err error
	logrus.Infof("Remote storage type is '%s'", remoteStorageType)

	remoteStorage, err = selectRemoteStorage(src)
	if err != nil {
		logrus.Errorf("Select remote storage error: %s", err)
	}
	if remoteStorage == nil {
		logrus.Errorf("Unknown remote storage type '%s', check the remote-storage-type param", remoteStorageType)
	}

	switch cmd {
	case "makebkp":
		if err := runBackupProcessToRemoteStorage(remoteStorage); err != nil {
			logrus.Errorf("Run backup process to remote storage error: %s", err)
		}
	case "getRemoteFiles":
		if err = getFilesFromRemoteStorage(remoteStorage); err != nil {
			logrus.Errorf("Getting files from remote storage error: %s", err)
		}
	}
}
