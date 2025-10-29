package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config 全局配置
type Config struct {
	Watermark WatermarkConfig `toml:"watermark"`
	Output    OutputConfig    `toml:"output"`
}

// WatermarkConfig 水印配置
type WatermarkConfig struct {
	FontSize float64 `toml:"font_size"` // 字体大小（相对于图像宽度的比例，默认 0.033）
	Position string  `toml:"position"`  // 位置: top-left, top-center, top-right, bottom-left, bottom-center, bottom-right, center
	Opacity  float64 `toml:"opacity"`   // 透明度 0-1
	Color    string  `toml:"color"`     // 颜色 (hex 格式，如 #FFFFFF)
	Margin   int     `toml:"margin"`    // 边距（像素）
}

// OutputConfig 输出配置
type OutputConfig struct {
	Quality int    `toml:"quality"` // JPEG 质量 (1-100)
	Format  string `toml:"format"`  // 输出格式: auto, jpeg, png
}

var (
	// GlobalConfig 全局配置实例
	GlobalConfig *Config
	// ConfigPath 配置文件路径
	ConfigPath string
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Watermark: WatermarkConfig{
			FontSize: 0.01,
			Position: "bottom-center",
			Opacity:  0.7,
			Color:    "#FFFFFF",
			Margin:   160,
		},
		Output: OutputConfig{
			Quality: 95,
			Format:  "auto",
		},
	}
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 如果没有指定配置文件，尝试从默认位置加载
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	cfg := DefaultConfig()

	// 如果配置文件存在，则加载
	if _, err := os.Stat(configPath); err == nil {
		if _, err := toml.DecodeFile(configPath, cfg); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
		fmt.Printf("✅ 已加载配置文件: %s\n", configPath)
	} else {
		// 配置文件不存在，使用默认配置
		if configPath != getDefaultConfigPath() {
			// 如果是用户指定的配置文件但不存在，返回错误
			return nil, fmt.Errorf("配置文件不存在: %s", configPath)
		}
		// 使用默认配置
		fmt.Println("⚠️  未找到配置文件，使用默认配置")
	}

	ConfigPath = configPath
	GlobalConfig = cfg
	return cfg, nil
}

// getDefaultConfigPath 获取默认配置文件路径
func getDefaultConfigPath() string {
	// 尝试使用 XDG_CONFIG_HOME
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "xpix", "config.toml")
	}

	// 使用 ~/.config/xpix/config.toml
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "xpix", "config.toml")
}

// GetDefaultConfigPath 获取默认配置路径（导出给外部使用）
func GetDefaultConfigPath() string {
	return getDefaultConfigPath()
}

// CreateDefaultConfig 创建默认配置文件
func CreateDefaultConfig(path string) error {
	if path == "" {
		path = getDefaultConfigPath()
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 创建配置文件
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	fmt.Printf("✅ 已创建默认配置文件: %s\n", path)
	return nil
}

// Get 获取全局配置
func Get() *Config {
	if GlobalConfig == nil {
		GlobalConfig = DefaultConfig()
	}
	return GlobalConfig
}

