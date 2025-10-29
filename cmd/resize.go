package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/processor"
)

var (
	resizeWidth  int
	resizeHeight int
	keepRatio    bool
)

var resizeCmd = &cobra.Command{
	Use:   "resize [image]",
	Short: "调整图像尺寸",
	Long:  `调整图像到指定的宽度和高度`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		opts := processor.ResizeOptions{
			Width:     resizeWidth,
			Height:    resizeHeight,
			KeepRatio: keepRatio,
		}

		if output == "" {
			output = addSuffix(inputPath, "_resized")
		}

		return processor.Resize(inputPath, output, opts)
	},
}

func init() {
	rootCmd.AddCommand(resizeCmd)

	resizeCmd.Flags().IntVarP(&resizeWidth, "width", "w", 0, "目标宽度")
	resizeCmd.Flags().IntVarP(&resizeHeight, "height", "h", 0, "目标高度")
	resizeCmd.Flags().BoolVarP(&keepRatio, "keep-ratio", "k", true, "保持宽高比")
	resizeCmd.Flags().StringVarP(&output, "output", "o", "", "输出文件路径")
}

