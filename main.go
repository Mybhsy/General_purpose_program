package main

import (
	"errors"
	"fmt"
	"general_purpose_program/utils"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var window fyne.Window

// 创建主菜单
var mainMenu = fyne.NewMainMenu(
	fyne.NewMenu("帮助",
		fyne.NewMenuItem("关于", func() {
			showAboutDialog()
		}),
	),
)

// showAboutDialog 显示关于对话框
func showAboutDialog() {


	// 创建版本信息
	version := widget.NewLabel("版本: 1.0.0")

	// 创建作者信息
	author := widget.NewLabel("作者: 张三")

	// 创建邮箱信息
	email := widget.NewLabel("邮箱: zhangsan@example.com")

	// 创建版权信息
	copyright := widget.NewLabel("版权: © 2024 保留所有权利")

	// 创建描述信息
	description := widget.NewLabel("这是一个实用的工具软件，集成了以下功能：\n\n1. 文件批量重命名：\n   - 支持通过Excel表格配置新旧文件名，轻松完成大量文件的重命名操作\n   - 适用于批量整理照片、文档、音视频等各类文件\n   - 支持预览重命名结果，避免误操作\n   - 操作简单，只需准备Excel文件和选择目标文件夹即可完成\n\n2. 文字转语音：\n   - 基于微软Edge TTS引擎，提供专业级语音合成服务\n   - 支持中文、英文、日文等多种语言，可选择不同性别和风格的发音人\n   - 语速调节范围-50%至+50%，满足不同场景需求\n   - 音量可调节，确保输出音频清晰舒适\n   - 支持Excel表格批量导入文本，自动转换并保存为MP3格式\n   - 适用于配音、教学、有声书制作等多种应用场景")
	description.Wrapping = fyne.TextWrapWord

	// 创建可滚动的内容区域
	scrollContent := container.NewVBox(
		version,
		author,
		email,
		copyright,
		description,
	)

	// 创建滚动容器
	scroll := container.NewScroll(scrollContent)
	scroll.SetMinSize(fyne.NewSize(400, 200))

	// 创建主容器
	content := container.NewVBox(
		scroll,
	)

	// 显示对话框
	dialog := dialog.NewCustom("关于", "确认", content, window)
	dialog.Resize(fyne.NewSize(400, 300))
	dialog.Show()
}

func main() {
	// 创建应用
	myApp := app.New()
	
	// 创建主窗口
	window = myApp.NewWindow("通用工具项目")
	
	// 设置主菜单
	window.SetMainMenu(mainMenu)
	
	// 创建标签页容器
	tabs := container.NewAppTabs(
		container.NewTabItem("文件重命名", createRenameTab()),
		container.NewTabItem("文字转语音", createTTSTab()),
	)
	
	// 设置标签页位置
	tabs.SetTabLocation(container.TabLocationTop)
	
	// 设置窗口内容
	window.SetContent(tabs)
	
	// 设置窗口大小
	window.Resize(fyne.NewSize(800, 600))
	
	// 显示并运行
	window.ShowAndRun()
}

// createRenameTab 创建文件重命名标签页
func createRenameTab() fyne.CanvasObject {
	// 状态变量
	var excelPath, folderPath string
	statusLabel := widget.NewLabel("准备就绪")

	// Excel文件选择按钮
	selectExcelBtn := widget.NewButton("选择Excel文件", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			excelPath = reader.URI().Path()
			statusLabel.SetText("已选择Excel文件: " + excelPath)
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	})

	// 目标文件夹选择按钮
	selectFolderBtn := widget.NewButton("选择目标文件夹", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if uri == nil {
				return
			}
			folderPath = uri.Path()
			statusLabel.SetText("已选择目标文件夹: " + folderPath)
		}, window)
		fd.Show()
	})

	// 开始重命名按钮
	startRenameBtn := widget.NewButton("开始重命名", func() {
		if excelPath == "" || folderPath == "" {
			dialog.ShowError(errors.New("请先选择Excel文件和目标文件夹"), window)
			return
		}

		go func() {
			statusLabel.SetText("正在读取Excel文件...")
			data, err := utils.ReadExcelForRename(excelPath)
			if err != nil {
				dialog.ShowError(err, window)
				statusLabel.SetText("读取Excel文件失败")
				return
			}

			statusLabel.SetText("正在重命名文件...")
			if err := utils.RenameFiles(folderPath, data); err != nil {
				dialog.ShowError(err, window)
				statusLabel.SetText("重命名失败")
				return
			}

			statusLabel.SetText("重命名完成")
			dialog.ShowInformation("成功", "文件重命名完成", window)
		}()
	})

	// 布局
	return container.NewVBox(
		selectExcelBtn,
		selectFolderBtn,
		startRenameBtn,
		statusLabel,
	)
}

// createTTSTab 创建文字转语音标签页
func createTTSTab() fyne.CanvasObject {
	// 状态变量
	var excelPath, outputPath string
	var config utils.TTSConfig
	statusLabel := widget.NewLabel("准备就绪")

	// 初始化TTS配置
	config = utils.TTSConfig{
		Voice:  "zh-CN-XiaoxiaoNeural",
		Rate:   "+0%",
		Volume: "+0%",
	}

	// Excel文件选择按钮
	selectExcelBtn := widget.NewButton("选择Excel文件", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			excelPath = reader.URI().Path()
			statusLabel.SetText("已选择Excel文件: " + excelPath)
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	})

	// 输出文件夹选择按钮
	selectOutputFolderBtn := widget.NewButton("选择输出文件夹", func() {
		fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if uri == nil {
				return
			}
			outputPath = uri.Path()
			statusLabel.SetText("已选择输出文件夹: " + outputPath)
		}, window)
		fd.Show()
	})

	// 获取可用的语音列表
	voices, err := utils.GetVoiceList()
	if err != nil {
		dialog.ShowError(err, window)
		voices = []string{"zh-CN-XiaoxiaoNeural"} // 如果获取失败，使用默认值
	}

	// 语音参数设置
	voiceOptions := make([]string, len(voices))
	voiceMap := make(map[string]string)
	for i, voice := range voices {
		parts := strings.Split(voice, "|")
		if len(parts) == 2 {
			voiceMap[parts[1]] = parts[0]
			voiceOptions[i] = parts[0] // 使用显示名称作为选项
		} else {
			voiceMap[voice] = voice
			voiceOptions[i] = voice
		}
	}
	voiceSelect := widget.NewSelect(voiceOptions, func(value string) {
		// 从显示名称中提取实际的语音名称
		for shortName, displayName := range voiceMap {
			if displayName == value {
				config.Voice = shortName
				break
			}
		}
	})
	voiceSelect.SetSelected(voiceMap["zh-CN-XiaoxiaoNeural"]) // 设置默认选项为小娜的显示名称

	// 语速调节
	rateSlider := widget.NewSlider(-50, 50)
	rateSlider.OnChanged = func(value float64) {
		config.Rate = fmt.Sprintf("%+.0f%%", value)
	}
	rateSlider.SetValue(0)

	// 音量调节
	volumeSlider := widget.NewSlider(-50, 50)
	volumeSlider.OnChanged = func(value float64) {
		config.Volume = fmt.Sprintf("%+.0f%%", value)
	}
	volumeSlider.SetValue(0)

	// 开始转换按钮
	startConvertBtn := widget.NewButton("开始转换", func() {
		if excelPath == "" || outputPath == "" {
			dialog.ShowError(errors.New("请先选择Excel文件和输出文件夹"), window)
			return
		}

		go func() {
			statusLabel.SetText("正在读取Excel文件...")
			texts, err := utils.ReadExcelForTTS(excelPath)
			if err != nil {
				dialog.ShowError(err, window)
				statusLabel.SetText("读取Excel文件失败")
				return
			}

			statusLabel.SetText("正在转换语音...")
			if err := utils.BatchTextToSpeech(texts, outputPath, config); err != nil {
				dialog.ShowError(err, window)
				statusLabel.SetText(fmt.Sprintf("转换失败: %v", err))
				return
			}

			statusLabel.SetText("转换完成")
			dialog.ShowInformation("成功", "语音转换完成", window)
		}()
	})

	// 布局
	return container.NewVBox(
		selectExcelBtn,
		selectOutputFolderBtn,
		container.NewGridWithColumns(2,
			widget.NewLabel("语音选择"),
			voiceSelect,
			widget.NewLabel("语速调节"),
			rateSlider,
			widget.NewLabel("音量调节"),
			volumeSlider,
		),
		startConvertBtn,
		statusLabel,
	)
}