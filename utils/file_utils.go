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

// RenameFiles batch renames files
func RenameFiles(folderPath string, renameData []ExcelData) error {
	for _, data := range renameData {
		oldPath := filepath.Join(folderPath, data.OldName)
		newPath := filepath.Join(folderPath, data.NewName)

		// Check if source file exists
		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			return fmt.Errorf("文件不存在: %s", oldPath)
		}

		// Execute rename
		if err := os.Rename(oldPath, newPath); err != nil {
			return fmt.Errorf("重命名失败 %s -> %s: %v", oldPath, newPath, err)
		}
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