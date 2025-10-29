# xpix

xpix 是一个强大的图像处理命令行工具，使用 Go 语言开发。

## 特性

- 🎨 **调色**：亮度、对比度、饱和度调整
- ✂️ **裁剪**：精确裁剪图像区域
- 🖼️ **水印**：支持文字和图片水印
- 📐 **尺寸调整**：灵活的图像缩放

## 安装

```bash
# 克隆项目
git clone https://github.com/xiaoheiwowo/xpix.git
cd xpix

# 安装依赖
go mod download

# 编译
go build -o xpix

# 或直接运行
go run main.go
```

## 配置

xpix 支持通过配置文件管理默认参数。

### 创建配置文件

```bash
# 创建默认配置文件到 ~/.config/xpix/config.toml
xpix config init

# 查看配置文件路径
xpix config path

# 查看当前配置
xpix config show
```

### 使用自定义配置

```bash
# 使用指定的配置文件
xpix -c ./my-config.toml watermark photo.jpg --text "© 2025"

# 使用默认配置（~/.config/xpix/config.toml）
xpix watermark photo.jpg --text "© 2025"
```

### 配置文件示例

配置文件使用 TOML 格式，参考 `config.example.toml`：

```toml
[watermark]
# 注意: 字体文件固定为 ~/Library/Fonts/CaskaydiaMonoNerdFont-Regular.ttf
# 请将字体文件放置到该位置

font_size = 0.033  # 字体大小（相对图像宽度）
position = "bottom-center"  # 水印位置（默认底部居中）
opacity = 0.7
color = "#FFFFFF"  # 白色
margin = 200  # 边距（像素）

[output]
quality = 95  # JPEG 质量
format = "auto"
```

**字体要求：**

工具使用固定字体路径：`~/Library/Fonts/CaskaydiaMonoNerdFont-Regular.ttf`

请确保该字体文件存在，否则水印功能无法使用。

## 使用示例

### 调整亮度、对比度

```bash
# 增加亮度
xpix adjust photo.jpg --brightness 20

# 增加对比度和饱和度
xpix adjust photo.jpg --contrast 15 --saturation 10 -o output.jpg
```

### 调整图像尺寸

```bash
# 调整宽度，保持宽高比
xpix resize photo.jpg --width 800

# 调整到指定尺寸（不保持宽高比）
xpix resize photo.jpg --width 800 --height 600 --keep-ratio=false
```

### 裁剪图像

```bash
# 从 (100, 100) 开始裁剪 500x500 的区域
xpix crop photo.jpg --x 100 --y 100 --width 500 --height 500
```

### 添加水印

```bash
# 添加文字水印（使用配置文件设置）
xpix watermark photo.jpg --text "© 2025 MyName"

# 添加中文水印
xpix watermark photo.jpg --text "版权所有 © 2025"

# 添加 Emoji 水印 🎨
xpix watermark photo.jpg --text "📸 My Photo 2025 ✨"

# 添加文字水印（自定义位置和透明度）
xpix watermark photo.jpg --text "© 2025" --position bottom-center --opacity 0.7

# 添加图片水印
xpix watermark photo.jpg --image logo.png --position top-left
```

### 查看图像信息

```bash
# 显示图像的 EXIF 元数据
xpix info photo.jpg
```

## 命令参考

### 全局参数

| 参数 | 简写 | 说明 |
|------|------|------|
| `--config` | `-c` | 配置文件路径（默认: `~/.config/xpix/config.toml`） |
| `--help` | `-h` | 显示帮助信息 |
| `--version` | `-v` | 显示版本信息 |

### `xpix config`

配置文件管理。

**子命令：**
- `xpix config init` - 创建默认配置文件
- `xpix config path` - 显示配置文件路径
- `xpix config show` - 显示当前配置

### `xpix info`

显示图像元数据信息。

| 参数 | 说明 |
|------|------|
| `[image]` | 图像文件路径 |

显示内容包括：
- 📁 文件信息（文件名、大小、修改时间）
- 📷 EXIF 元数据（相机型号、拍摄参数、镜头信息）
- 📍 GPS 位置信息（如果有）

### `xpix adjust`

调整图像的亮度、对比度、饱和度。

| 参数 | 简写 | 说明 | 范围 |
|------|------|------|------|
| `--brightness` | `-b` | 亮度调整 | -100 到 100 |
| `--contrast` | `-c` | 对比度调整 | -100 到 100 |
| `--saturation` | `-s` | 饱和度调整 | -100 到 100 |
| `--output` | `-o` | 输出文件路径 | - |

### `xpix resize`

调整图像尺寸。

| 参数 | 简写 | 说明 |
|------|------|------|
| `--width` | `-w` | 目标宽度 |
| `--height` | `-h` | 目标高度 |
| `--keep-ratio` | `-k` | 保持宽高比（默认: true） |
| `--output` | `-o` | 输出文件路径 |

### `xpix crop`

裁剪图像。

| 参数 | 简写 | 说明 |
|------|------|------|
| `--x` | `-x` | 裁剪起始 X 坐标 |
| `--y` | `-y` | 裁剪起始 Y 坐标 |
| `--width` | `-w` | 裁剪宽度（必需） |
| `--height` | `-h` | 裁剪高度（必需） |
| `--output` | `-o` | 输出文件路径 |

### `xpix watermark`

添加文字或图片水印。

| 参数 | 简写 | 说明 |
|------|------|------|
| `--text` | `-t` | 文字水印内容 |
| `--image` | - | 图片水印路径 |
| `--position` | `-p` | 水印位置（默认: bottom-right） |
| `--opacity` | - | 水印透明度 0-1（默认: 0.5） |
| `--output` | `-o` | 输出文件路径 |

支持的位置：`top-left`, `top-center`, `top-right`, `bottom-left`, `bottom-center`, `bottom-right`, `center`

## 项目结构

```
xpix/
├── main.go                 # 程序入口
├── go.mod                  # 依赖管理
├── cmd/                    # 命令行命令
│   ├── root.go            # 根命令
│   ├── adjust.go          # 调色命令
│   ├── resize.go          # 尺寸调整命令
│   ├── crop.go            # 裁剪命令
│   └── watermark.go       # 水印命令
└── internal/              # 内部包
    └── processor/         # 图像处理逻辑
        ├── adjust.go      # 调色处理
        ├── resize.go      # 尺寸调整处理
        ├── crop.go        # 裁剪处理
        └── watermark.go   # 水印处理
```

## 扩展指南

### 添加新命令

1. 在 `cmd/` 目录下创建新的命令文件（如 `filter.go`）
2. 在 `internal/processor/` 下实现处理逻辑
3. 在命令的 `init()` 函数中注册到 `rootCmd`

示例：

```go
// cmd/filter.go
package cmd

import (
    "github.com/spf13/cobra"
    "github.com/xiaoheiwowo/xpix/internal/processor"
)

var filterCmd = &cobra.Command{
    Use:   "filter [image]",
    Short: "应用滤镜",
    RunE: func(cmd *cobra.Command, args []string) error {
        // 实现逻辑
        return processor.ApplyFilter(args[0], output, opts)
    },
}

func init() {
    rootCmd.AddCommand(filterCmd)
}
```

## 特性说明

### Unicode 和 Emoji 支持

xpix 完全支持 Unicode 字符和 Emoji 表情水印：

```bash
# 中文水印
xpix watermark photo.jpg --text "我的照片 © 2025"

# 日文水印
xpix watermark photo.jpg --text "私の写真 © 2025"

# Emoji 水印
xpix watermark photo.jpg --text "📸 Travel 2025 ✈️🌍"

# 混合使用
xpix watermark photo.jpg --text "🎨 作品集 2025 ©️ Artist"
```

**字体要求：**

请将字体文件放置到固定位置：

```bash
~/Library/Fonts/CaskaydiaMonoNerdFont-Regular.ttf
```

该字体支持完整的 Unicode 字符和 Emoji 表情。

## 依赖

- [cobra](https://github.com/spf13/cobra) - CLI 框架
- [imaging](https://github.com/disintegration/imaging) - 图像处理库
- [gg](https://github.com/fogleman/gg) - 2D 图形绘制（用于文字水印）
- [toml](https://github.com/BurntSushi/toml) - TOML 配置文件解析
- [goexif](https://github.com/rwcarlsen/goexif) - EXIF 元数据读取

## License

MIT

