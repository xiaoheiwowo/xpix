package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/processor"
)

var (
	cropX      int
	cropY      int
	cropWidth  int
	cropHeight int
)

var cropCmd = &cobra.Command{
	Use:   "crop [image]",
	Short: "裁剪图像",
	Long:  `裁剪图像到指定的尺寸和位置`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		opts := processor.CropOptions{
			X:      cropX,
			Y:      cropY,
			Width:  cropWidth,
			Height: cropHeight,
		}

		if output == "" {
			output = addSuffix(inputPath, "_cropped")
		}

		return processor.Crop(inputPath, output, opts)
	},
}

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.Flags().IntVarP(&cropX, "x", "x", 0, "裁剪起始 X 坐标")
	cropCmd.Flags().IntVarP(&cropY, "y", "y", 0, "裁剪起始 Y 坐标")
	cropCmd.Flags().IntVarP(&cropWidth, "width", "w", 0, "裁剪宽度")
	cropCmd.Flags().IntVarP(&cropHeight, "height", "h", 0, "裁剪高度")
	cropCmd.Flags().StringVarP(&output, "output", "o", "", "输出文件路径")

	cropCmd.MarkFlagRequired("width")
	cropCmd.MarkFlagRequired("height")
}

