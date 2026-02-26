package automation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"picture-rpa/constants"
)

type ExecutionContext struct {
	IsRunning  *bool
	LogFunc    func(string, string)
	Conf       float64
	LoopCount  int
	Commands   []constants.Command
}

func toInterfaces(ss []string) []interface{} {
	is := make([]interface{}, len(ss))
	for i, s := range ss {
		is[i] = s
	}
	return is
}

func (ctx *ExecutionContext) ExecuteWorkflow() {
	i := 0
	for i < len(ctx.Commands) && *ctx.IsRunning {
		cmd := ctx.Commands[i]
		nextStep := i + 1

		ctx.LogFunc(fmt.Sprintf("Step %d: %s", i+1, constants.REV_CMD_TYPES[cmd.Type]), "info")

		switch cmd.Type {
		case 1.0, 2.0, 3.0: // Click, Double Click, Right Click (with image)
			clicks := 1
			button := "left"
			if cmd.Type == 2.0 {
				clicks = 2
			} else if cmd.Type == 3.0 {
				button = "right"
			}
			if ctx.mouseClick(clicks, button, cmd.Value, cmd.Retry) {
				ctx.LogFunc("Found and clicked", "info")
			} else {
				ctx.LogFunc("Not found: "+filepath.Base(cmd.Value), "warn")
			}
		case 4.0: // Text Input
			clipboard.WriteAll(cmd.Value)
			robotgo.KeyTap("v", "command")
			time.Sleep(500 * time.Millisecond)
		case 5.0: // Sleep
			val, _ := strconvToFloat(cmd.Value)
			duration := time.Duration(val * float64(time.Second))
			time.Sleep(duration)
		case 6.0: // Scroll
			var scrollAmount int
			fmt.Sscanf(cmd.Value, "%d", &scrollAmount)
			robotgo.Scroll(0, scrollAmount)
		case 7.0: // Hotkey
			keys := strings.Split(strings.ToLower(cmd.Value), "+")
			if len(keys) > 1 {
				tapKey := keys[len(keys)-1]
				modifiers := toInterfaces(keys[:len(keys)-1])
				robotgo.KeyTap(tapKey, modifiers...)
			} else {
				robotgo.KeyTap(keys[0])
			}
		case 8.0: // Jump if found
			if ctx.findImage(cmd.Value) {
				ctx.LogFunc(fmt.Sprintf("Found, jump to %d", cmd.Retry), "info")
				nextStep = cmd.Retry - 1
			}
		case 9.0: // Jump if not found
			if !ctx.findImage(cmd.Value) {
				ctx.LogFunc(fmt.Sprintf("Not found, jump to %d", cmd.Retry), "info")
				nextStep = cmd.Retry - 1
			}
		case 10.0: // Stop
			ctx.LogFunc("Workflow forced to stop", "err")
			*ctx.IsRunning = false
			return
		}

		i = nextStep
	}
}

func strconvToFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func (ctx *ExecutionContext) mouseClick(clickTimes int, button string, imgPath string, retry int) bool {
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		ctx.LogFunc("File not found: "+imgPath, "err")
		return false
	}

	maxAttempts := retry
	if retry <= 1 {
		maxAttempts = 5
	}
	if retry == -1 {
		maxAttempts = 999999
	}

	for i := 0; i < maxAttempts; i++ {
		if !*ctx.IsRunning {
			break
		}
		x, y := robotgo.FindPic(imgPath)
		if x != -1 && y != -1 {
			robotgo.Move(x, y)
			robotgo.Click(button, clickTimes == 2)
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

func (ctx *ExecutionContext) findImage(imgPath string) bool {
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		return false
	}
	x, y := robotgo.FindPic(imgPath)
	return x != -1 && y != -1
}
