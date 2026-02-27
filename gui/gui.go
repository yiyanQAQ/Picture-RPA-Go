package gui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"picture-rpa/automation"
	"picture-rpa/constants"

	"github.com/go-vgo/robotgo"
)

type PictureRPAGUI struct {
	Window    fyne.Window
	Commands  []constants.Command
	Table     *widget.Table
	Logs      *widget.List
	LogData   []string
	IsRunning bool
}

func NewPictureRPAGUI() *PictureRPAGUI {
	a := app.New()
	w := a.NewWindow("Picture RPA v3.0 - Logic & Snip")
	w.Resize(fyne.NewSize(800, 700))

	g := &PictureRPAGUI{
		Window:   w,
		Commands: []constants.Command{},
		LogData:  []string{},
	}

	g.buildUI()
	return g
}

func (g *PictureRPAGUI) buildUI() {
	// Top Frame (Form)
	cmdTypes := []string{
		"单击左键 (找图)",
		"双击左键 (找图)",
		"单击右键 (找图)",
		"文本输入 (复制粘贴)",
		"等待时长 (秒)",
		"鼠标滚轮 (滑动)",
		"执行快捷键 (Hotkey)",
		"跳转:若找到图 (Goto)",
		"跳转:若没找到图 (Goto)",
		"结束流程 (Stop)",
	}
	typeSelect := widget.NewSelect(cmdTypes, func(s string) {})
	typeSelect.SetSelectedIndex(0)

	valEntry := widget.NewEntry()
	valEntry.SetPlaceHolder("内容/路径:")

	retryEntry := widget.NewEntry()
	retryEntry.SetText("1")

	snipBtn := widget.NewButtonWithIcon("截图保存", theme.ContentCopyIcon(), func() {
		g.Window.Hide()
		time.Sleep(500 * time.Millisecond)

		if _, err := os.Stat("snips"); os.IsNotExist(err) {
			os.Mkdir("snips", 0755)
		}

		savePath := fmt.Sprintf("snips/snip_%d.png", time.Now().Unix())
		robotgo.SaveCapture(savePath)
		g.Window.Show()

		g.Log("截图已保存: "+savePath, "info")
		valEntry.SetText(savePath)
	})

	browseBtn := widget.NewButtonWithIcon("浏览图片", theme.FileIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				valEntry.SetText(reader.URI().Path())
			}
		}, g.Window)
		fd.Show()
	})

	addBtn := widget.NewButtonWithIcon("添加到流程", theme.ContentAddIcon(), func() {
		t := typeSelect.Selected
		v := valEntry.Text
		r, _ := strconv.Atoi(retryEntry.Text)
		g.Commands = append(g.Commands, constants.Command{
			Type:  constants.CMD_TYPES[t],
			Value: v,
			Retry: r,
		})
		g.Table.Refresh()
		g.Log("添加到流程: "+t, "info")
	})

	form := container.New(layout.NewVBoxLayout(),
		container.NewGridWithColumns(2,
			widget.NewLabel("操作类型:"), typeSelect,
			widget.NewLabel("内容/路径:"), valEntry,
			widget.NewLabel("重试次数:"), retryEntry,
		),
		container.NewHBox(snipBtn, browseBtn, layout.NewSpacer(), addBtn),
	)

	g.Table = widget.NewTable(
		func() (int, int) {
			return len(g.Commands), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("CellContent")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			if i.Row >= len(g.Commands) {
				return
			}
			cmd := g.Commands[i.Row]
			switch i.Col {
			case 0:
				l.SetText(strconv.Itoa(i.Row + 1))
			case 1:
				l.SetText(constants.REV_CMD_TYPES[cmd.Type])
			case 2:
				l.SetText(cmd.Value)
			case 3:
				l.SetText(strconv.Itoa(cmd.Retry))
			}
		},
	)
	g.Table.SetColumnWidth(0, 50)
	g.Table.SetColumnWidth(1, 150)
	g.Table.SetColumnWidth(2, 400)
	g.Table.SetColumnWidth(3, 100)

	g.Logs = widget.NewList(
		func() int {
			return len(g.LogData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(g.LogData[i])
		},
	)

	runOnceBtn := widget.NewButtonWithIcon("Once执行", theme.MediaPlayIcon(), func() {
		g.startRun(false)
	})
	runLoopBtn := widget.NewButtonWithIcon("Infinity循环", theme.ViewRefreshIcon(), func() {
		g.startRun(true)
	})
	stopBtn := widget.NewButtonWithIcon("Stop停止", theme.MediaStopIcon(), func() {
		g.IsRunning = false
		g.Log("ESC触发，已停止程序", "err")
	})

	saveBtn := widget.NewButton("Save保存", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				data, _ := json.MarshalIndent(g.Commands, "", "  ")
				writer.Write(data)
				writer.Close()
			}
		}, g.Window)
		fd.Show()
	})

	loadBtn := widget.NewButton("Load加载", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				data, _ := ioutil.ReadAll(reader)
				json.Unmarshal(data, &g.Commands)
				g.Table.Refresh()
				reader.Close()
			}
		}, g.Window)
		fd.Show()
	})

	controls := container.NewHBox(saveBtn, loadBtn, layout.NewSpacer(), runOnceBtn, runLoopBtn, stopBtn)

	content := container.NewBorder(form, controls, nil, nil, container.NewVSplit(g.Table, g.Logs))
	g.Window.SetContent(content)
}

func (g *PictureRPAGUI) Log(msg, tag string) {
	t := time.Now().Format("15:04:05")
	g.LogData = append(g.LogData, fmt.Sprintf("[%s] %s", t, msg))
	if g.Logs != nil {
		g.Logs.Refresh()
		g.Logs.ScrollToBottom()
	}
}

func (g *PictureRPAGUI) startRun(loop bool) {
	if len(g.Commands) == 0 {
		return
	}
	g.IsRunning = true
	ctx := &automation.ExecutionContext{
		IsRunning: &g.IsRunning,
		LogFunc:   g.Log,
		Commands:  g.Commands,
	}

	go func() {
		defer func() {
			g.IsRunning = false
			g.Log("RPA 等待中", "info")
		}()

		if loop {
			for g.IsRunning {
				ctx.ExecuteWorkflow()
			}
		} else {
			ctx.ExecuteWorkflow()
		}
	}()
}
