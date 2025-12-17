package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var makeBkpCmd = &cobra.Command{
	Use:   "makebkp",
	Short: "Создать бэкап.",
	Long:  `Команда создаёт архив из указанного каталога и отправляет его в хранилище.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("Let's start creating local backup...")
		src := managerLocal()

		logrus.Infof("Let's start creating remote backup...")
		managerRemoteStorage(src, cmd.Use)

		if err := src.handleLocalArchive(); err != nil {
			logrus.Errorf("Removing archive %s error: %s", src.targetArchiveFullName, err)
		}

		logrus.Infof("Cleanup local...")
		if err := src.cleanUpArchivesLocal(archiveStorageDepth); err != nil {
			logrus.Errorf("Cleanup archives local error: %s", err)
		}

		logrus.Info("Done")
	},
}

func managerLocal() *src {
	src := NewSrc()

	src.setTargetArchiveFullName()

	if !src.isPathExists() {
		logrus.Fatalf("%s dir not found, check the path", src.srcDir)
	}
	logrus.Infof("Local dir %s is found", src.srcDir)

	if err := src.zipFiles(); err != nil {
		panic("ZIP archiving error")
	}
	logrus.Infof("Backup archive %s created locally", src.targetArchivePath+"/"+src.targetArchiveFullName)

	return src
}

func runBackupProcessToRemoteStorage(rs RemoteStorage) error {
	if err := rs.CreateRoot(); err != nil {
		return err
	}

	if err := rs.Upload(); err != nil {
		return err
	}

	return nil
}
