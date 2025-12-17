package cmd

import (
	"fmt"
	"os"
	"simple-bkp/storage"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type RemoteStorage interface {
	CreateRoot() error
	Upload() error
	GetListFiles() ([]os.FileInfo, error)
	Remove(f string) error
	GetRootDir() string
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
		logrus.Warnf("Unknown remote storage type '%s', check the remote-storage-type param (or you don't need a backup on a remote storage)", remoteStorageType)
		return
	}

	switch cmd {
	case "makebkp":
		if err := runBackupProcessToRemoteStorage(remoteStorage); err != nil {
			logrus.Errorf("Run backup process to remote storage error: %s", err)
		}
		logrus.Info("Cleanup remote...")
		if err = cleanUpArchivesRemote(remoteStorage, archiveStorageDepth); err != nil {
			logrus.Error(err)
		}
	case "getRemoteFiles":
		if err = getFilesFromRemoteStorage(remoteStorage); err != nil {
			logrus.Errorf("Getting files from remote storage error: %s", err)
		}
	}
}

func cleanUpArchivesRemote(rs RemoteStorage, days string) error {
	d, err := strconv.Atoi(days)
	if err != nil {
		return err
	}

	remoteFiles, err := rs.GetListFiles()
	if err != nil {
		return err
	}

	for _, f := range remoteFiles {
		modTime := f.ModTime()
		if time.Since(modTime) > (time.Duration(d) * 24 * time.Hour) {
			if err := rs.Remove(rs.GetRootDir() + "/" + f.Name()); err != nil {
				return fmt.Errorf("remove old remote file error: %w", err)
			}
			logrus.Infof("Removed old remote file: %s (age %v days)", "/"+rs.GetRootDir()+"/"+f.Name(), int(time.Since(modTime).Hours()/24))
		} else {
			logrus.Infof("Remote file is recent: %s", "/"+rs.GetRootDir()+"/"+f.Name())
		}
	}

	if (len(remoteFiles) - 1) > d {
		logrus.Warn("More remote files than you need")
	}

	return nil
}
