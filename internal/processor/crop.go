package processor

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

// CropOptions 裁剪选项
type CropOptions struct {
	X      int // 起始 X 坐标
	Y      int // 起始 Y 坐标
	Width  int // 裁剪宽度
	Height int // 裁剪高度
}

// Crop 裁剪图像
func Crop(inputPath, outputPath string, opts CropOptions) error {
	// 打开图像
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开图像: %w", err)
	}

	// 裁剪
	rect := image.Rect(opts.X, opts.Y, opts.X+opts.Width, opts.Y+opts.Height)
	result := imaging.Crop(img, rect)

	// 保存结果
	if err := imaging.Save(result, outputPath); err != nil {
		return fmt.Errorf("无法保存图像: %w", err)
	}

	fmt.Printf("✅ 图像已保存至: %s\n", outputPath)
	return nil
}

