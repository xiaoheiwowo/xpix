package processor

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

// ResizeOptions 调整尺寸选项
type ResizeOptions struct {
	Width     int  // 目标宽度
	Height    int  // 目标高度
	KeepRatio bool // 保持宽高比
}

// Resize 调整图像尺寸
func Resize(inputPath, outputPath string, opts ResizeOptions) error {
	// 打开图像
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开图像: %w", err)
	}

	var result image.Image

	if opts.KeepRatio {
		// 保持宽高比
		if opts.Width > 0 && opts.Height == 0 {
			result = imaging.Resize(img, opts.Width, 0, imaging.Lanczos)
		} else if opts.Height > 0 && opts.Width == 0 {
			result = imaging.Resize(img, 0, opts.Height, imaging.Lanczos)
		} else {
			// 两者都指定时，缩略图模式
			result = imaging.Fit(img, opts.Width, opts.Height, imaging.Lanczos)
		}
	} else {
		// 强制调整到指定尺寸
		result = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	}

	// 保存结果
	if err := imaging.Save(result, outputPath); err != nil {
		return fmt.Errorf("无法保存图像: %w", err)
	}

	fmt.Printf("✅ 图像已保存至: %s\n", outputPath)
	return nil
}

