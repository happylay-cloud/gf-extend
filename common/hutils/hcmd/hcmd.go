package hcmd

import (
	"os/exec"
	"runtime"

	"github.com/gogf/gf/frame/g"
)

// 不同平台启动指令
var commands = map[string]string{
	"windows": "start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

// OpenBrowser 打开浏览器
func OpenBrowser(url string) error {
	command, ok := commands[runtime.GOOS]
	if !ok {
		g.Log().Line(false).Error("在" + runtime.GOOS + "平台无法打开浏览器")
	}

	// 获取cmd命令
	cmd := exec.Command(command, url)

	// Run会阻塞等待命令完成
	return cmd.Run()
}
