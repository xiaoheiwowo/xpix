package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/processor"
)

var (
	brightness float64
	contrast   float64
	saturation float64
	output     string
)

var adjustCmd = &cobra.Command{
	Use:   "adjust [image]",
	Short: "调整图像的亮度、对比度、饱和度等",
	Long: `对图像进行调色处理，支持：
  - 亮度调整 (--brightness)
  - 对比度调整 (--contrast)
  - 饱和度调整 (--saturation)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputPath := args[0]
		
		opts := processor.AdjustOptions{
			Brightness: brightness,
			Contrast:   contrast,
			Saturation: saturation,
		}

		if output == "" {
			output = addSuffix(inputPath, "_adjusted")
		}

		return processor.Adjust(inputPath, output, opts)
	},
}

func init() {
	rootCmd.AddCommand(adjustCmd)

	adjustCmd.Flags().Float64VarP(&brightness, "brightness", "b", 0, "亮度调整 (-100 到 100)")
	adjustCmd.Flags().Float64VarP(&contrast, "contrast", "c", 0, "对比度调整 (-100 到 100)")
	adjustCmd.Flags().Float64VarP(&saturation, "saturation", "s", 0, "饱和度调整 (-100 到 100)")
	adjustCmd.Flags().StringVarP(&output, "output", "o", "", "输出文件路径")
}

// addSuffix 为文件名添加后缀
func addSuffix(path, suffix string) string {
	ext := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext = path[i:]
			path = path[:i]
			break
		}
	}
	return fmt.Sprintf("%s%s%s", path, suffix, ext)
}

