package processor

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/xiaoheiwowo/xpix/internal/config"
)

// WatermarkOptions 水印选项
type WatermarkOptions struct {
	Text     string  // 文字水印内容
	Image    string  // 图片水印路径
	Position string  // 位置: top-left, top-right, bottom-left, bottom-right, center
	Opacity  float64 // 透明度 0-1
}

// Watermark 添加水印
func Watermark(inputPath, outputPath string, opts WatermarkOptions) error {
	// 打开图像
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开图像: %w", err)
	}

	var result image.Image

	// 根据类型处理水印
	if opts.Text != "" {
		result = addTextWatermark(img, opts)
	} else if opts.Image != "" {
		result, err = addImageWatermark(img, opts)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("请指定文字水印 (--text) 或图片水印 (--image)")
	}

	// 保存结果
	if err := imaging.Save(result, outputPath); err != nil {
		return fmt.Errorf("无法保存图像: %w", err)
	}

	fmt.Printf("✅ 图像已保存至: %s\n", outputPath)
	return nil
}

func addTextWatermark(img image.Image, opts WatermarkOptions) image.Image {
	bounds := img.Bounds()
	dc := gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawImage(img, 0, 0)

	// 获取配置
	cfg := config.Get()

	// 计算字体大小
	fontSize := float64(bounds.Dx()) * cfg.Watermark.FontSize

	// 加载字体（使用固定字体路径）
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("⚠️  无法获取用户目录: %v\n", err)
		return img
	}

	fontPath := home + "/Library/Fonts/CaskaydiaMonoNerdFont-Regular.ttf"

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		fmt.Printf("⚠️  无法加载字体 %s: %v\n", fontPath, err)
		fmt.Println("⚠️  请将字体文件放置到 ~/Library/Fonts/CaskaydiaMonoNerdFont-Regular.ttf")
		return img
	}

	// 使用配置的透明度（如果命令行未指定）
	opacity := opts.Opacity
	if opacity == 0 {
		opacity = cfg.Watermark.Opacity
	}

	// 解析颜色
	textColor := parseColor(cfg.Watermark.Color, opacity)
	dc.SetColor(textColor)

	// 使用配置的位置和边距
	position := opts.Position
	if position == "" {
		position = cfg.Watermark.Position
	}

	// 计算位置
	margin := float64(cfg.Watermark.Margin)
	x, y := calculateTextPosition(bounds.Dx(), bounds.Dy(), position, margin)

	// 绘制文字（支持 Unicode 和 Emoji）
	dc.DrawStringAnchored(opts.Text, x, y, 0.5, 0.5)

	return dc.Image()
}

// parseColor 解析颜色字符串（支持 #RRGGBB 格式）
func parseColor(colorStr string, opacity float64) color.Color {
	// 移除 # 前缀
	colorStr = strings.TrimPrefix(colorStr, "#")

	// 默认白色
	if colorStr == "" {
		return color.RGBA{R: 255, G: 255, B: 255, A: uint8(opacity * 255)}
	}

	// 解析十六进制颜色
	var r, g, b uint8 = 255, 255, 255

	if len(colorStr) == 6 {
		if val, err := strconv.ParseUint(colorStr[0:2], 16, 8); err == nil {
			r = uint8(val)
		}
		if val, err := strconv.ParseUint(colorStr[2:4], 16, 8); err == nil {
			g = uint8(val)
		}
		if val, err := strconv.ParseUint(colorStr[4:6], 16, 8); err == nil {
			b = uint8(val)
		}
	}

	return color.RGBA{R: r, G: g, B: b, A: uint8(opacity * 255)}
}

// calculateTextPosition 计算文字位置
func calculateTextPosition(width, height int, position string, margin float64) (float64, float64) {
	switch position {
	case "top-left":
		return margin, margin
	case "top-right":
		return float64(width) - margin, margin
	case "bottom-left":
		return margin, float64(height) - margin
	case "bottom-right":
		return float64(width) - margin, float64(height) - margin
	case "center":
		return float64(width) / 2, float64(height) / 2
	case "bottom-center":
		return float64(width) / 2, float64(height) - margin
	case "top-center":
		return float64(width) / 2, margin
	default:
		return float64(width) / 2, float64(height) - margin
	}
}

func addImageWatermark(img image.Image, opts WatermarkOptions) (image.Image, error) {
	// 打开水印图像
	watermark, err := imaging.Open(opts.Image)
	if err != nil {
		return nil, fmt.Errorf("无法打开水印图像: %w", err)
	}

	bounds := img.Bounds()

	// 调整水印大小（最大为原图的 1/5）
	maxSize := bounds.Dx() / 5
	if watermark.Bounds().Dx() > maxSize {
		watermark = imaging.Resize(watermark, maxSize, 0, imaging.Lanczos)
	}

	// 调整透明度
	watermark = adjustOpacity(watermark, opts.Opacity)

	// 计算位置
	x, y := calculateImagePosition(bounds.Dx(), bounds.Dy(), watermark.Bounds().Dx(), watermark.Bounds().Dy(), opts.Position)

	// 合成
	result := imaging.Overlay(img, watermark, image.Pt(x, y), 1.0)

	return result, nil
}

func calculateImagePosition(imgWidth, imgHeight, wmWidth, wmHeight int, position string) (int, int) {
	margin := 20
	switch position {
	case "top-left":
		return margin, margin
	case "top-right":
		return imgWidth - wmWidth - margin, margin
	case "bottom-left":
		return margin, imgHeight - wmHeight - margin
	case "bottom-right":
		return imgWidth - wmWidth - margin, imgHeight - wmHeight - margin
	case "center":
		return (imgWidth - wmWidth) / 2, (imgHeight - wmHeight) / 2
	case "bottom-center":
		return (imgWidth - wmWidth) / 2, imgHeight - wmHeight - margin
	case "top-center":
		return (imgWidth - wmWidth) / 2, margin
	default:
		return (imgWidth - wmWidth) / 2, imgHeight - wmHeight - margin
	}
}

// adjustOpacity 调整图像透明度
func adjustOpacity(img image.Image, opacity float64) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// 遍历每个像素，调整 alpha 通道
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			c.A = uint8(float64(c.A) * opacity)
			dst.Set(x, y, c)
		}
	}

	return dst
}
