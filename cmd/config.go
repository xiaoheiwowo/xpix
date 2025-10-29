package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置文件管理",
	Long:  `管理 xpix 配置文件`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "创建默认配置文件",
	Long:  `在默认位置创建配置文件模板`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := config.GetDefaultConfigPath()
		return config.CreateDefaultConfig(path)
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "显示配置文件路径",
	Long:  `显示当前使用的配置文件路径`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.ConfigPath != "" {
			fmt.Printf("当前配置文件: %s\n", config.ConfigPath)
		} else {
			fmt.Printf("默认配置路径: %s\n", config.GetDefaultConfigPath())
		}
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Long:  `显示当前生效的配置内容`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		fmt.Println("当前配置:")
		fmt.Println()
		fmt.Println("[watermark]")
		fmt.Printf("  font_size = %.3f\n", cfg.Watermark.FontSize)
		fmt.Printf("  position = \"%s\"\n", cfg.Watermark.Position)
		fmt.Printf("  opacity = %.2f\n", cfg.Watermark.Opacity)
		fmt.Printf("  color = \"%s\"\n", cfg.Watermark.Color)
		fmt.Printf("  margin = %d\n", cfg.Watermark.Margin)
		fmt.Println()
		fmt.Println("[output]")
		fmt.Printf("  quality = %d\n", cfg.Output.Quality)
		fmt.Printf("  format = \"%s\"\n", cfg.Output.Format)
		fmt.Println()
		fmt.Println("注意: 字体文件固定为 ~/Library/Fonts/CaskaydiaMonoNerdFont-Regular.ttf")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configPathCmd)
	configCmd.AddCommand(configShowCmd)
}

