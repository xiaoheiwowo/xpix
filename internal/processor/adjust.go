package processor

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
)

// AdjustOptions 调色选项
type AdjustOptions struct {
	Brightness  float64 // -100 到 100
	Contrast    float64 // -100 到 100
	Saturation  float64 // -100 到 100
	Exposure    float64 // -100 到 100
	Sharpen     float64 // 0 到 100
	Gamma       float64 // 0.1 到 3.0
	Temperature int     // 色温 K (2000-10000，6500 为标准日光)
	Dehaze      float64 // 0 到 100
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

	// 去雾处理（先处理）
	if opts.Dehaze > 0 {
		result = applyDehaze(result, opts.Dehaze)
	}

	// 色温调整
	if opts.Temperature != 0 && opts.Temperature != 6500 {
		result = applyTemperature(result, opts.Temperature)
	}

	// 曝光调整
	if opts.Exposure != 0 {
		result = imaging.AdjustBrightness(result, opts.Exposure)
	}

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

	// Gamma 调整
	if opts.Gamma != 0 && opts.Gamma != 1.0 {
		result = imaging.AdjustGamma(result, opts.Gamma)
	}

	// 锐化处理（最后处理）
	if opts.Sharpen > 0 {
		sharpness := opts.Sharpen / 10.0 // 转换为合适的范围
		result = imaging.Sharpen(result, sharpness)
	}

	return result
}

// applyTemperature 应用色温调整
func applyTemperature(img image.Image, kelvin int) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// 将色温转换为 RGB 调整因子
	// 参考标准：6500K 为标准日光白平衡
	rFactor, gFactor, bFactor := kelvinToRGB(kelvin)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			r := clamp(float64(oldColor.R) * rFactor)
			g := clamp(float64(oldColor.G) * gFactor)
			b := clamp(float64(oldColor.B) * bFactor)

			dst.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: oldColor.A,
			})
		}
	}

	return dst
}

// kelvinToRGB 将色温（开尔文）转换为 RGB 调整因子
// 基于 Tanner Helland 的算法
func kelvinToRGB(kelvin int) (float64, float64, float64) {
	// 限制色温范围
	if kelvin < 1000 {
		kelvin = 1000
	}
	if kelvin > 40000 {
		kelvin = 40000
	}

	temp := float64(kelvin) / 100.0
	var r, g, b float64

	// 计算红色
	if temp <= 66 {
		r = 255
	} else {
		r = temp - 60
		r = 329.698727446 * math.Pow(r, -0.1332047592)
		if r < 0 {
			r = 0
		}
		if r > 255 {
			r = 255
		}
	}

	// 计算绿色
	if temp <= 66 {
		g = temp
		g = 99.4708025861*math.Log(g) - 161.1195681661
		if g < 0 {
			g = 0
		}
		if g > 255 {
			g = 255
		}
	} else {
		g = temp - 60
		g = 288.1221695283 * math.Pow(g, -0.0755148492)
		if g < 0 {
			g = 0
		}
		if g > 255 {
			g = 255
		}
	}

	// 计算蓝色
	if temp >= 66 {
		b = 255
	} else if temp <= 19 {
		b = 0
	} else {
		b = temp - 10
		b = 138.5177312231*math.Log(b) - 305.0447927307
		if b < 0 {
			b = 0
		}
		if b > 255 {
			b = 255
		}
	}

	// 转换为调整因子（相对于 6500K 标准日光）
	return r / 255.0, g / 255.0, b / 255.0
}

// applyDehaze 应用去雾效果
func applyDehaze(img image.Image, strength float64) image.Image {
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)

	// 将强度转换为合适的因子 (0-100 -> 0-1)
	factor := strength / 100.0

	// 去雾：增加对比度和饱和度，降低灰度
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			r := float64(oldColor.R)
			g := float64(oldColor.G)
			b := float64(oldColor.B)

			// 计算灰度值
			gray := (r + g + b) / 3

			// 增加对比度（拉伸色阶）
			r = clamp(gray + (r-gray)*(1+factor*0.5))
			g = clamp(gray + (g-gray)*(1+factor*0.5))
			b = clamp(gray + (b-gray)*(1+factor*0.5))

			// 增加饱和度
			max := math.Max(math.Max(r, g), b)
			min := math.Min(math.Min(r, g), b)
			if max > min {
				saturation := 1 + factor*0.3
				r = clamp(max - (max-r)*saturation)
				g = clamp(max - (max-g)*saturation)
				b = clamp(max - (max-b)*saturation)
			}

			dst.Set(x, y, color.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: oldColor.A,
			})
		}
	}

	return dst
}

// clamp 将值限制在 0-255 范围内
func clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}
