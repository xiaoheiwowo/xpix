package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/processor"
)

var infoCmd = &cobra.Command{
	Use:   "info [image]",
	Short: "显示图像的元数据信息",
	Long:  `显示图像的 EXIF 元数据信息，包括相机型号、拍摄参数、GPS 位置等`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		imagePath := args[0]
		return processor.ShowImageInfo(imagePath)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

