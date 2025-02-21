package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetVoiceList 获取Edge TTS支持的语音列表
func GetVoiceList() ([]string, error) {
	// 执行edge-tts --list-voices命令
	cmd := exec.Command("edge-tts", "--list-voices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取语音列表失败: %v", err)
	}

	// 将输出转换为字符串并按行分割
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("无效的语音列表输出")
	}

	// 跳过标题行
	var formattedVoices []string
	for _, line := range lines[2:] { // 从第三行开始（跳过标题和分隔线）
		// 使用空格分割，并过滤掉空字符串
		fields := strings.Fields(line)
		if len(fields) < 2 { // 确保至少有Name和Gender字段
			continue
		}

		name := fields[0]
		gender := fields[1]

		// 从名称中提取语言代码和发音人名称
		parts := strings.Split(name, "-")
		if len(parts) < 3 {
			continue
		}
		locale := parts[0] + "-" + parts[1]
		
		// 解析语言代码
		language := getLanguageName(locale)
		// 转换性别显示
		genderDisplay := getGenderName(gender)
		// 提取发音人名称（去掉Neural后缀）
		voiceName := strings.TrimSuffix(parts[len(parts)-1], "Neural")

		// 格式化为"语言-性别-发音人|ShortName"形式
		formattedVoice := fmt.Sprintf("%s-%s-%s|%s", language, genderDisplay, voiceName, name)
		formattedVoices = append(formattedVoices, formattedVoice)
	}

	return formattedVoices, nil
}

// getLanguageName 将语言代码转换为语言名称
func getLanguageName(locale string) string {
	languageMap := map[string]string{
		"zh-CN": "中文(简体)",
		"zh-TW": "中文(繁体)",
		"zh-HK": "中文(香港)",
		"en-US": "英语(美国)",
		"en-GB": "英语(英国)",
		"en-AU": "英语(澳大利亚)",
		"en-CA": "英语(加拿大)",
		"en-IN": "英语(印度)",
		"en-IE": "英语(爱尔兰)",
		"en-KE": "英语(肯尼亚)",
		"en-NG": "英语(尼日利亚)",
		"en-NZ": "英语(新西兰)",
		"en-PH": "英语(菲律宾)",
		"en-SG": "英语(新加坡)",
		"en-TZ": "英语(坦桑尼亚)",
		"en-ZA": "英语(南非)",
		"ja-JP": "日语(日本)",
		"ko-KR": "韩语(韩国)",
		"fr-FR": "法语(法国)",
		"fr-CA": "法语(加拿大)",
		"fr-CH": "法语(瑞士)",
		"fr-BE": "法语(比利时)",
		"de-DE": "德语(德国)",
		"de-AT": "德语(奥地利)",
		"de-CH": "德语(瑞士)",
		"ru-RU": "俄语(俄罗斯)",
		"es-ES": "西班牙语(西班牙)",
		"es-MX": "西班牙语(墨西哥)",
		"es-AR": "西班牙语(阿根廷)",
		"es-CO": "西班牙语(哥伦比亚)",
		"es-PE": "西班牙语(秘鲁)",
		"es-VE": "西班牙语(委内瑞拉)",
		"es-CL": "西班牙语(智利)",
		"it-IT": "意大利语(意大利)",
		"pt-BR": "葡萄牙语(巴西)",
		"pt-PT": "葡萄牙语(葡萄牙)",
		"ar-EG": "阿拉伯语(埃及)",
		"ar-SA": "阿拉伯语(沙特)",
		"ar-AE": "阿拉伯语(阿联酋)",
		"ar-BH": "阿拉伯语(巴林)",
		"ar-KW": "阿拉伯语(科威特)",
		"ar-QA": "阿拉伯语(卡塔尔)",
		"hi-IN": "印地语(印度)",
		"th-TH": "泰语(泰国)",
		"vi-VN": "越南语(越南)",
		"id-ID": "印尼语(印尼)",
		"ms-MY": "马来语(马来西亚)",
		"tr-TR": "土耳其语(土耳其)",
		"pl-PL": "波兰语(波兰)",
		"nl-NL": "荷兰语(荷兰)",
		"nl-BE": "荷兰语(比利时)",
		"sv-SE": "瑞典语(瑞典)",
		"da-DK": "丹麦语(丹麦)",
		"fi-FI": "芬兰语(芬兰)",
		"el-GR": "希腊语(希腊)",
		"he-IL": "希伯来语(以色列)",
		"ro-RO": "罗马尼亚语(罗马尼亚)",
		"hu-HU": "匈牙利语(匈牙利)",
		"cs-CZ": "捷克语(捷克)",
		"sk-SK": "斯洛伐克语(斯洛伐克)",
		"uk-UA": "乌克兰语(乌克兰)",
		"bn-IN": "孟加拉语(印度)",
		"ta-IN": "泰米尔语(印度)",
		"te-IN": "泰卢固语(印度)",
		"ml-IN": "马拉雅拉姆语(印度)",
		"kn-IN": "卡纳达语(印度)",
		"mr-IN": "马拉地语(印度)",
		"gu-IN": "古吉拉特语(印度)",
		"pa-IN": "旁遮普语(印度)",
	}

	parts := strings.Split(locale, "-")
	if len(parts) >= 2 {
		if name, ok := languageMap[locale]; ok {
			return name
		}
	}
	return locale
}

// getGenderName 将性别代码转换为显示名称
func getGenderName(gender string) string {
	switch gender {
	case "Female":
		return "女"
	case "Male":
		return "男"
	default:
		return gender
	}
}

// getVoiceName 从ShortName中提取发音人名称
func getVoiceName(shortName string) string {
	// 移除Neural后缀
	name := strings.TrimSuffix(shortName, "Neural")
	// 提取最后一个部分作为名称
	parts := strings.Split(name, "-")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return shortName
}