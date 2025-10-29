package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/processor"
)

var (
	watermarkText     string
	watermarkImage    string
	watermarkPosition string
	watermarkOpacity  float64
)

var watermarkCmd = &cobra.Command{
	Use:   "watermark [image]",
	Short: "添加文字或图片水印",
	Long: `为图像添加水印，支持：
  - 文字水印 (--text)
  - 图片水印 (--image)
  - 位置控制 (--position: top-left, top-right, bottom-left, bottom-right, center)
  - 透明度控制 (--opacity)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]

		opts := processor.WatermarkOptions{
			Text:     watermarkText,
			Image:    watermarkImage,
			Position: watermarkPosition,
			Opacity:  watermarkOpacity,
		}

		if output == "" {
			output = addSuffix(inputPath, "_watermarked")
		}

		return processor.Watermark(inputPath, output, opts)
	},
}

func init() {
	rootCmd.AddCommand(watermarkCmd)

	watermarkCmd.Flags().StringVarP(&watermarkText, "text", "t", "", "文字水印内容")
	watermarkCmd.Flags().StringVar(&watermarkImage, "image", "", "图片水印路径")
	watermarkCmd.Flags().StringVarP(&watermarkPosition, "position", "p", "bottom-center", "水印位置")
	watermarkCmd.Flags().Float64Var(&watermarkOpacity, "opacity", 0.5, "水印透明度 (0-1)")
	watermarkCmd.Flags().StringVarP(&output, "output", "o", "", "输出文件路径")
}

