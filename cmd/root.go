package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "simple-bkp",
	Short: "Утилита для простого резервного копирования файлов.",
}

var (
	srcDir                   = ""
	targetArchiveLocalPath   = ""
	targetArchiveNamePrefix  = ""
	remoteStorageType        = ""
	shouldRemoveLocalArchive = false
	archiveStorageDepth      = ""
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(makeBkpCmd)
	makeBkpCmd.Flags().StringVar(
		&srcDir, "src-dir", "", "Указывает путь к локальному каталогу, для которого будет осуществляться резервное копирование.")
	makeBkpCmd.Flags().StringVar(
		&targetArchiveLocalPath, "target-archive-local-path", "", "Указывает путь к локальному каталогу для архива.")
	makeBkpCmd.Flags().StringVar(
		&targetArchiveNamePrefix, "target-archive-name-prefix", "", "Указывает префикс имени архива.")
	makeBkpCmd.Flags().StringVar(
		&remoteStorageType, "remote-storage-type", "", "Указывает тип удаленного хранилища резервных копий (WebDav, s3, SMB...)")
	makeBkpCmd.Flags().BoolVar(
		&shouldRemoveLocalArchive, "remove-local-archive", false, "Указывает на необходимость удаления локальной копии архива.")
	makeBkpCmd.Flags().StringVar(
		&archiveStorageDepth, "archive-storage-depth", "3", "Глубина хранения резервных копий как локально, так и удаленно в днях.")

	rootCmd.AddCommand(getRemoteFilesCmd)
	getRemoteFilesCmd.Flags().StringVar(
		&remoteStorageType, "remote-storage-type", "", "Указывает тип удаленного хранилища резервных копий (WebDav, s3, SMB...)")

	rootCmd.AddCommand(versionCmd)
}
