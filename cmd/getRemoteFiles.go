package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getRemoteFilesCmd = &cobra.Command{
	Use:   "getRemoteFiles",
	Short: "Получить список файлов в удаленном хранилище.",
	Long:  `Команда подключается к удаленному хранилищу, и возвращает список файлов`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("Getting files from remote storage...")

		src := NewSrc()
		src.setTargetArchiveFullName()

		managerRemoteStorage(src, cmd.Use)

		logrus.Info("Done")
	},
}

func getFilesFromRemoteStorage(rs RemoteStorage) error {
	fl, err := rs.GetListFiles()
	if err != nil {
		return err
	}

	for _, f := range fl {
		fmt.Println("/"+rs.GetRootDir()+"/"+f.Name(), ":", f.Size())
	}

	return nil
}
