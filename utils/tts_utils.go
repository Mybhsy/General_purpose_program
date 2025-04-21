package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// TTSConfig 存储TTS配置信息
type TTSConfig struct {
	Voice string
	Rate  string
	Volume string
}

// TextToSpeech 将文本转换为语音
func TextToSpeech(text, outputPath string, config TTSConfig) error {
	// 验证输入参数
	if text = strings.TrimSpace(text); text == "" {
		return fmt.Errorf("文本内容不能为空")
	}

	if outputPath = strings.TrimSpace(outputPath); outputPath == "" {
		return fmt.Errorf("输出路径不能为空")
	}

	// 验证配置参数
	if config.Voice = strings.TrimSpace(config.Voice); config.Voice == "" {
		return fmt.Errorf("语音配置不能为空")
	}

	if config.Rate = strings.TrimSpace(config.Rate); config.Rate == "" {
		config.Rate = "+0%" // 使用默认语速
	} else if !strings.HasSuffix(config.Rate, "%") || !(strings.HasPrefix(config.Rate, "+") || strings.HasPrefix(config.Rate, "-")) {
		return fmt.Errorf("无效的语速格式，应为'+N%'或'-N%'的格式，如'+0%'")
	}

	if config.Volume = strings.TrimSpace(config.Volume); config.Volume == "" {
		config.Volume = "+0%" // 使用默认音量
	} else if !strings.HasSuffix(config.Volume, "%") || !(strings.HasPrefix(config.Volume, "+") || strings.HasPrefix(config.Volume, "-")) {
		return fmt.Errorf("无效的音量格式，应为'+N%'或'-N%'的格式，如'+0%'")
	}

	// 检查edge-tts命令是否可用
	_, err := exec.LookPath("edge-tts")
	if err != nil {
		return fmt.Errorf("未找到edge-tts命令，请确保已安装: %v", err)
	}

	// 确保输出目录存在
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 生成输出文件名 ‌在Windows系统中，文件名的最大长度为255个字符‌‌,使用128大多数情况下足够
	maxLen := 128
	if len(text) < maxLen {
		maxLen = len(text)
	}
	fileName := strings.ReplaceAll(text[:maxLen], " ", "_")
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName = strings.ReplaceAll(fileName, "\\", "_")
	fileName = strings.ReplaceAll(fileName, "?", "_")
	fileName = strings.ReplaceAll(fileName, "*", "_")
	fileName = strings.ReplaceAll(fileName, ":", "_")
	fileName = strings.ReplaceAll(fileName, "<", "_")
	fileName = strings.ReplaceAll(fileName, ">", "_")
	fileName = strings.ReplaceAll(fileName, "|", "_")
	fileName = strings.ReplaceAll(fileName, "\"", "_")
	outputFile := filepath.Join(outputPath, fileName+".mp3")

	// 构建edge-tts命令
	cmd := exec.Command("edge-tts",
		"--voice", config.Voice,
		"--rate", config.Rate,
		"--volume", config.Volume,
		"--text", text,
		"--write-media", outputFile,
	)

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// 检查命令执行结果
	if err != nil || 
	   strings.Contains(strings.ToLower(outputStr), "error") || 
	   strings.Contains(outputStr, "失败") || 
	   strings.Contains(strings.ToLower(outputStr), "invalid") || 
	   strings.Contains(strings.ToLower(outputStr), "failed") {
		// 清理可能生成的无效文件
		os.Remove(outputFile)
		return fmt.Errorf("TTS转换失败: %v\n命令输出: %s", err, outputStr)
	}

	// 验证输出文件
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("音频文件生成失败，未找到输出文件: %s", outputFile)
		}
		return fmt.Errorf("检查输出文件失败: %v", err)
	}

	// 验证文件大小
	if fileInfo.Size() < 100 {
		os.Remove(outputFile)
		return fmt.Errorf("生成的音频文件大小异常（%d字节），可能转换失败", fileInfo.Size())
	}

	return nil
}

// BatchTextToSpeech 批量将文本转换为语音
func BatchTextToSpeech(texts []string, outputPath string, config TTSConfig, progressCallback ...ProgressCallback) error {
    total := len(texts)
	retryCount := 10 // 重试次数

    for i, text := range texts {
        if text == "" {
            continue // 跳过空文本
        }
        
        // 计算进度百分比
        percentage := float64(i) / float64(total) * 100
        
        // 如果提供了回调函数，则调用它
        if len(progressCallback) > 0 && progressCallback[0] != nil {
            progressCallback[0](i, total, percentage)
        }
        
		for j := 0; j < retryCount; j++ {
            if err := TextToSpeech(text, outputPath, config); err != nil {
                if (j < retryCount - 1) {
					return fmt.Errorf("第 %d/%d 个文本转换失败（文本内容：%s）：%v", i+1, total, text, err)	
				}
            } else {
				break // 转换成功，跳出重试循环
			}
		}

    }
    
    // 完成时调用回调函数，进度为100%
    if len(progressCallback) > 0 && progressCallback[0] != nil {
        progressCallback[0](total, total, 100.0)
    }
    
    return nil
}