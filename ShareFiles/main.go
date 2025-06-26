/*
FileName: main.go
Create on: 2025-06-24
Author: ChinRing
Description: 实现查找文件大小大于输入值的文件
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 定义常量
const (
	KB = 1024      // KB 转换 Byte
	MB = KB * 1024 // MB 转换 Byte
	GB = MB * 1024 // GB 转换 Byte
)

// 定义变量
var rootPath string       // 目录路径
var fileSize uint64       // 文件大小
var sizeUnit string       // 大小单位
var shareKeyWord string   // 关键词
var fileExt string        // 文件拓展名
var matchFiles []FileInfo // 处理后符合条件的文件

// 结构体,用于存储文件路径和文件大小
type FileInfo struct {
	FilePath string
	FileSize uint64
}

// 关键字与拓展名函数结构一致,可以使用多态去实现
type FileFilter interface {
	Match(file string) bool
}

type KeyWordFilter struct {
	KeyWord string
}

type ExtFilter struct {
	Ext string
}

func (kw KeyWordFilter) Match(file string) bool {
	return strings.Contains(filepath.Base(file), kw.KeyWord)
}

func (e ExtFilter) Match(file string) bool {
	ext := e.Ext
	if ext != "" && ext[0] != '.' {
		ext = "." + ext
	}
	return filepath.Ext(file) == ext
}

func checkFilter(filter FileFilter, file string) bool {
	return filter.Match(file)
}

// 初始化函数
func init() {

	flag.StringVar(&rootPath, "dir", "./", "指定查找的文件根目录")
	flag.Uint64Var(&fileSize, "size", 0, "指定查找的文件大小")
	flag.StringVar(&sizeUnit, "unit", "KB", "指定查找文件大小的单位, 默认KB(不区分大小写): [KB,MB,GB]")
	flag.StringVar(&shareKeyWord, "keyword", "", "通过关键字查找")
	// flag.StringVar(&fileExt, "ext", "", "通过文件拓展名查找") // TODO

}

// 单位转换用于计算文件大小是否超过预定值
func convertToBytes(size uint64, unit string) (uint64, error) {

	switch strings.ToUpper(unit) {

	case "KB":
		return size * KB, nil

	case "MB":
		return size * MB, nil

	case "GB":
		return size * GB, nil

	default:
		return 0, errors.New("单位不为[KB,MB,GB]中的一种")

	}

}

// 单位转换输出
func formatSize(size uint64) string {

	switch {

	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))

	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))

	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))

	default:
		return fmt.Sprintf("%d B", size)

	}

}

// 判断文件大小
func checkFileSize(file string, size uint64, unit string) (bool, error) {

	var actualSize uint64

	fileInfo, _ := os.Stat(file)
	actualSize = uint64(fileInfo.Size())

	limitSize, err := convertToBytes(size, unit)
	if err != nil {
		return false, err
	}
	return actualSize > limitSize, nil

}

// 遍历目录
func ReadDir(path string) []FileInfo {

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		absolutePath := filepath.Join(path, file.Name())
		result, err := checkFileSize(absolutePath, fileSize, sizeUnit)

		if err != nil {
			log.Fatal(err)
		} else if checkFilter(KeyWordFilter{KeyWord: shareKeyWord}, absolutePath) && result {
			fileBytesSize, _ := os.Stat(absolutePath)
			matchFiles = append(
				matchFiles,
				FileInfo{FilePath: absolutePath, FileSize: uint64(fileBytesSize.Size())})
		}

	}

	// 递归遍历子目录
	for _, file := range files {

		if file.IsDir() {
			subPath := filepath.Join(path, file.Name())
			ReadDir(subPath)
		}

	}
	return matchFiles
}

// 程序主入口
func main() {

	file, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	mainFile := filepath.Base(file)

	// 修改 flag.Usage, 输出使用信息
	flag.Usage = func() {
		fmt.Printf("用法: %s [参数]", mainFile)
		fmt.Println("可用参数: ")
		flag.PrintDefaults()
	}

	// 解析参数
	flag.Parse()

	result := ReadDir(rootPath)
	for _, StructData := range result {
		pattern := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(shareKeyWord))
		text := StructData.FilePath
		highlighted := pattern.ReplaceAllString(text, "\033[31m$0\033[0m")
		fmt.Println(highlighted, formatSize(StructData.FileSize))
	}

}
