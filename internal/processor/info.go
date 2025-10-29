package processor

import (
	"fmt"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
)

// ShowImageInfo 显示图像的元数据信息
func ShowImageInfo(imagePath string) error {
	// 打开图像文件
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("无法打开图像文件: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("无法获取文件信息: %w", err)
	}

	// 解析 EXIF 数据
	exifData, err := exif.Decode(file)

	// 显示基本文件信息
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📁 文件信息\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("文件名:   %s\n", fileInfo.Name())
	fmt.Printf("文件大小: %s\n", formatFileSize(fileInfo.Size()))
	fmt.Printf("修改时间: %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 如果没有 EXIF 数据
	if err != nil {
		fmt.Println("⚠️  该图像没有 EXIF 元数据")
		return nil
	}

	// 显示 EXIF 信息
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📷 EXIF 元数据\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 相机信息
	if make, err := exifData.Get(exif.Make); err == nil {
		fmt.Printf("相机品牌: %s\n", formatExifValue(make))
	}
	if model, err := exifData.Get(exif.Model); err == nil {
		fmt.Printf("相机型号: %s\n", formatExifValue(model))
	}
	if lens, err := exifData.Get(exif.LensModel); err == nil {
		fmt.Printf("镜头型号: %s\n", formatExifValue(lens))
	}

	fmt.Println()

	// 拍摄参数
	if dateTime, err := exifData.DateTime(); err == nil {
		weekday := getChineseWeekday(dateTime.Weekday())
		fmt.Printf("拍摄时间: %s %s\n", dateTime.Format("2006-01-02 15:04:05"), weekday)
	}
	if iso, err := exifData.Get(exif.ISOSpeedRatings); err == nil {
		fmt.Printf("ISO:      %s\n", formatExifValue(iso))
	}
	if fNumber, err := exifData.Get(exif.FNumber); err == nil {
		fmt.Printf("光圈:     f/%s\n", formatExifValue(fNumber))
	}
	if expTime, err := exifData.Get(exif.ExposureTime); err == nil {
		fmt.Printf("快门速度: %s 秒\n", formatExifValue(expTime))
	}
	if focalLength, err := exifData.Get(exif.FocalLength); err == nil {
		fmt.Printf("焦距:     %s mm\n", formatExifValue(focalLength))
	}
	if expBias, err := exifData.Get(exif.ExposureBiasValue); err == nil {
		fmt.Printf("曝光补偿: %s EV\n", formatExifValue(expBias))
	}

	fmt.Println()

	// 图像信息
	if width, err := exifData.Get(exif.PixelXDimension); err == nil {
		if height, err := exifData.Get(exif.PixelYDimension); err == nil {
			fmt.Printf("图像尺寸: %s x %s 像素\n", formatExifValue(width), formatExifValue(height))
		}
	}
	if orientation, err := exifData.Get(exif.Orientation); err == nil {
		fmt.Printf("方向:     %s\n", formatExifValue(orientation))
	}

	// GPS 信息
	lat, lon, err := exifData.LatLong()
	if err == nil {
		fmt.Println()
		fmt.Printf("📍 GPS 位置: %.6f, %.6f\n", lat, lon)
	}

	// 软件信息
	if software, err := exifData.Get(exif.Software); err == nil {
		fmt.Println()
		fmt.Printf("编辑软件: %s\n", formatExifValue(software))
	}

	// 显示所有 EXIF 标签（可选，用于调试）
	// fmt.Println()
	// fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	// fmt.Printf("所有 EXIF 标签\n")
	// fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	// exifData.Walk(func(name exif.FieldName, tag *tiff.Tag) error {
	// 	fmt.Printf("%s: %s\n", name, tag.String())
	// 	return nil
	// })

	return nil
}

// formatExifValue 格式化 EXIF 值
func formatExifValue(tag *tiff.Tag) string {
	if tag == nil {
		return ""
	}

	// 移除引号
	str := tag.String()
	if len(str) > 2 && str[0] == '"' && str[len(str)-1] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// formatDateTime 格式化日期时间
func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// getChineseWeekday 获取中文星期几
func getChineseWeekday(weekday time.Weekday) string {
	weekdays := map[time.Weekday]string{
		time.Sunday:    "星期日",
		time.Monday:    "星期一",
		time.Tuesday:   "星期二",
		time.Wednesday: "星期三",
		time.Thursday:  "星期四",
		time.Friday:    "星期五",
		time.Saturday:  "星期六",
	}
	return weekdays[weekday]
}
