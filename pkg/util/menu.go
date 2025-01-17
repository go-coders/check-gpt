package util

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Menu-specific emojis
const (
	EmojiDefault = "⭐"
	EmojiCustom  = "✏️"
	EmojiOpenAI  = "🤖"
	EmojiExit    = "👋"
	EmojiTool    = "🛠️"
	EmojiUpdate  = "🚀"
	EmojiBack    = "🔙"
)

// MenuItem represents a menu item
type MenuItem struct {
	ID       int
	Label    string
	Emoji    string
	Selected bool
}

// Menu represents a menu with title and items
type Menu struct {
	Title       string
	TitleEmoji  string
	Description string
	Items       []MenuItem
	Prompt      string
	ValidChoice func(string) bool
}

// Menu constants
var (
	MainMenu = Menu{
		Title:      "GPT 检测工具",
		TitleEmoji: EmojiTool,
	}

	MenuKey = Menu{
		Title:      "请选择测试模型",
		TitleEmoji: EmojiGear,
		Items: []MenuItem{
			{ID: 1, Label: "使用默认模型", Emoji: EmojiDefault},
			{ID: 2, Label: "自定义测试模型", Emoji: EmojiCustom},
		},
		Prompt: "请输入选择或输入自定义模型（多个模型用空格分隔）: ",
		ValidChoice: func(choice string) bool {
			return choice == "1" || choice == "2" || strings.Contains(choice, " ")
		},
	}

	MenuMain = Menu{
		Title:      "GPT 检测工具",
		TitleEmoji: EmojiTool,
		Description: fmt.Sprintf("%s项目地址: %shttps://github.com/go-coders/check-gpt%s",
			ColorGray, ColorBlue, ColorReset),
		Items: []MenuItem{
			{ID: 1, Label: "API Key 可用性测试", Emoji: EmojiKey},
			{ID: 2, Label: "API 中转链路检测", Emoji: EmojiLink},
			{ID: 3, Label: "检查更新", Emoji: EmojiGear},
			{ID: 4, Label: "退出", Emoji: EmojiExit},
		},
		Prompt: "请选择功能 (1-4): ",
		ValidChoice: func(choice string) bool {
			return choice >= "1" && choice <= "4"
		},
	}

	MenuUpdate = Menu{
		Title:      "检查更新",
		TitleEmoji: EmojiGear,
		Items: []MenuItem{
			{ID: 1, Label: "立即更新", Emoji: EmojiUpdate},
			{ID: 2, Label: "返回主菜单", Emoji: EmojiBack},
		},
		Prompt: "请选择操作: ",
		ValidChoice: func(choice string) bool {
			return choice == "1" || choice == "2"
		},
	}
)

// ShowMenu displays a menu and returns user's choice
func ShowMenu(menu Menu, input io.Reader, output io.Writer) (MenuItem, error) {
	printer := NewPrinter(output)

	// Clear screen and show title
	ClearConsole()
	printer.PrintTitle(menu.Title, menu.TitleEmoji)

	// Show description if present
	if menu.Description != "" {
		printer.Printf("%s", menu.Description)
		printer.PrintSeparator()
	}

	// Show menu items
	for _, item := range menu.Items {
		var format string
		if item.Selected {
			format = fmt.Sprintf("%d. %s  %s %s[已选择]%s\n",
				item.ID, item.Label, item.Emoji,
				ColorGray, ColorReset)
		} else {
			format = fmt.Sprintf("%d. %s  %s\n",
				item.ID, item.Label, item.Emoji)
		}
		printer.Print(format)
	}
	printer.Printf("\n%s", menu.Prompt)

	// Read user input
	reader := bufio.NewReader(input)
	for {
		choice, err := reader.ReadString('\n')
		if err != nil {
			return MenuItem{}, fmt.Errorf("读取选择失败: %v", err)
		}

		choice = strings.TrimSpace(choice)
		if menu.ValidChoice == nil || menu.ValidChoice(choice) {
			// get the item
			// convert choice to int
			choiceInt, err := strconv.Atoi(choice)
			if err != nil {
				return MenuItem{}, fmt.Errorf("选择无效: %v", err)
			}
			for _, item := range menu.Items {
				if item.ID == choiceInt {
					return item, nil
				}
			}
		}

		printer.Printf("%s", menu.Prompt)
	}
}

// ShowMainMenu displays the main menu and returns the user's choice
func ShowMainMenu(in io.Reader, out io.Writer) (MenuItem, error) {
	return ShowMenu(MenuMain, in, out)
}
