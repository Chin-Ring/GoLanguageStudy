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
	"strings"
)

// 定义常量
const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

// 定义变量
var rootPath string
var fileSize uint64
var sizeUnit string
var matchFiles []FileInfo

// 结构体,用于存储文件路径和文件大小
type FileInfo struct {
	FilePath string
	FileSize uint64
}

// 初始化函数
func init() {

	flag.StringVar(&rootPath, "dir", "./", "文件查找的根路径")
	flag.Uint64Var(&fileSize, "size", 0, "查找文件大小大于多少的文件")
	flag.StringVar(&sizeUnit, "unit", "KB", "查找文件大小的单位(不区分大小写): [KB,MB,GB]")

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
		return fmt.Sprintf("%f GB", float64(size)/float64(GB))

	case size >= MB:
		return fmt.Sprintf("%f MB", float64(size)/float64(MB))

	case size >= KB:
		return fmt.Sprintf("%f KB", float64(size)/float64(KB))

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
		} else if result {
			fileBytesSize, _ := os.Stat(absolutePath)
			matchFiles = append(matchFiles, FileInfo{FilePath: absolutePath, FileSize: uint64(fileBytesSize.Size())})
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
		fmt.Println(StructData.FilePath, formatSize(StructData.FileSize))
	}

}
