package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelData stores Excel file data
type ExcelData struct {
	OldName string
	NewName string
}

// ReadExcelForRename reads rename data from Excel file
func ReadExcelForRename(filePath string) ([]ExcelData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// Get first worksheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	var data []ExcelData
	// Skip header, start from second row
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) >= 2 {
			data = append(data, ExcelData{
				OldName: row[0],
				NewName: row[1],
			})
		}
	}

	return data, nil
}

// ProgressCallback 进度回调函数类型
type ProgressCallback func(current, total int, percentage float64)

// RenameFiles batch renames files and returns a list of failures
func RenameFiles(folderPath string, renameData []ExcelData, progressCallback ...ProgressCallback) error {
	var failures []string
	total := len(renameData)

	for i, data := range renameData {
		oldPath := filepath.Join(folderPath, data.OldName)
		newPath := filepath.Join(folderPath, data.NewName)

		// 计算进度百分比
        percentage := float64(i) / float64(total) * 100
		// 如果提供了回调函数，则调用它
        if len(progressCallback) > 0 && progressCallback[0] != nil {
            progressCallback[0](i, total, percentage)
        }

		// Check if source file exists
		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			failures = append(failures, fmt.Sprintf("文件不存在: %s", oldPath))
			continue
		}

		// Execute rename
		if err := os.Rename(oldPath, newPath); err != nil {
			failures = append(failures, fmt.Sprintf("重命名失败 %s -> %s: %v", oldPath, newPath, err))
			continue
		}
		
	}

	// 完成时调用回调函数，进度为100%
	if len(progressCallback) > 0 && progressCallback[0] != nil {
		progressCallback[0](total, total, 100.0)
	}

	// Return error with all failures if any
	if len(failures) > 0 {
		return fmt.Errorf("重命名过程中有 %d 个错误:\n%s", len(failures), strings.Join(failures, "\n"))
	} 
	return nil
}

// ReadExcelForTTS reads TTS text data from Excel file
func ReadExcelForTTS(filePath string) ([]string, error) {
	// 验证文件路径
	if filePath == "" {
		return nil, fmt.Errorf("Excel文件路径不能为空")
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Excel文件不存在: %s", filePath)
	}

	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("Excel文件中没有工作表")
	}

	// 读取工作表内容
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	// 检查是否有数据
	if len(rows) <= 1 { // 考虑到标题行
		return nil, fmt.Errorf("Excel文件中没有有效的文本数据")
	}

	var texts []string
	// 从第二行开始读取（跳过标题行）
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) >= 1 {
			// 验证文本内容
			text := strings.TrimSpace(row[0])
			if text != "" {
				texts = append(texts, text)
			}
		}
	}

	// 检查是否成功提取到文本
	if len(texts) == 0 {
		return nil, fmt.Errorf("未找到有效的文本内容")
	}

	return texts, nil
}