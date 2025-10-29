package processor

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

// AdjustOptions 调色选项
type AdjustOptions struct {
	Brightness float64 // -100 到 100
	Contrast   float64 // -100 到 100
	Saturation float64 // -100 到 100
}

// Adjust 调整图像的亮度、对比度、饱和度
func Adjust(inputPath, outputPath string, opts AdjustOptions) error {
	// 打开图像
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("无法打开图像: %w", err)
	}

	// 应用调整
	result := adjustImage(img, opts)

	// 保存结果
	if err := imaging.Save(result, outputPath); err != nil {
		return fmt.Errorf("无法保存图像: %w", err)
	}

	fmt.Printf("✅ 图像已保存至: %s\n", outputPath)
	return nil
}

func adjustImage(img image.Image, opts AdjustOptions) image.Image {
	result := img

	// 亮度调整
	if opts.Brightness != 0 {
		result = imaging.AdjustBrightness(result, opts.Brightness)
	}

	// 对比度调整
	if opts.Contrast != 0 {
		result = imaging.AdjustContrast(result, opts.Contrast)
	}

	// 饱和度调整
	if opts.Saturation != 0 {
		result = imaging.AdjustSaturation(result, opts.Saturation)
	}

	return result
}

