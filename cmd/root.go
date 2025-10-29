package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xiaoheiwowo/xpix/internal/config"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "xpix",
	Short: "xpix - 一个强大的图像处理命令行工具",
	Long: `xpix 是一个功能丰富的图像处理 CLI 工具。
支持调色、构图、水印、滤镜等多种图像处理功能。`,
	Version: "0.1.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 加载配置文件
		_, err := config.Load(cfgFile)
		return err
	},
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 全局配置标志
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", 
		fmt.Sprintf("配置文件路径 (默认: %s)", config.GetDefaultConfigPath()))
}

