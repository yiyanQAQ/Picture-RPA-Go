package constants

const (
	COLOR_BG_CREAM     = "#F3F3F3"
	COLOR_PINK_ROSE    = "#E1E1E1"
	COLOR_SLATE_PURPLE = "#0078D4"
)

var CMD_TYPES = map[string]float64{
	"单击左键 (找图)": 1.0,
	"双击左键 (找图)": 2.0,
	"单击右键 (找图)": 3.0,
	"文本输入 (复制粘贴)": 4.0,
	"等待时长 (秒)": 5.0,
	"鼠标滚轮 (滑动)": 6.0,
	"执行快捷键 (Hotkey)": 7.0,
	"跳转:若找到图 (Goto)": 8.0,
	"跳转:若没找到图 (Goto)": 9.0,
	"结束流程 (Stop)": 10.0,
}

var REV_CMD_TYPES = map[float64]string{
	1.0:  "单击左键 (找图)",
	2.0:  "双击左键 (找图)",
	3.0:  "单击右键 (找图)",
	4.0:  "文本输入 (复制粘贴)",
	5.0:  "等待时长 (秒)",
	6.0:  "鼠标滚轮 (滑动)",
	7.0:  "执行快捷键 (Hotkey)",
	8.0:  "跳转:若找到图 (Goto)",
	9.0:  "跳转:若没找到图 (Goto)",
	10.0: "结束流程 (Stop)",
}

type Command struct {
	Type  float64 `json:"type"`
	Value string  `json:"value"`
	Retry int     `json:"retry"`
}
