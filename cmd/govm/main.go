package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	cs "github.com/mitchellh/colorstring"

	"github.com/schollz/progressbar/v3"
)

var args []string
var actionArg = ""
var versionArg = ""
var version = "0.0.1.dev"

var allowedArgs = []string{"h", "help", "ls", "list", "ls-remote", "install", "use", "uninstall", "self-update"}

func init() {
	log.SetFlags(0)

	if !isArgAllowed() {
		log.Println("[Info] Invalid usage")
		cs.Println(usage())
		return
	}

	flag.Parse()
	args = flag.Args()
	if len(args) == 0 {
		cs.Println(usage())
		return
	}

	actionArg = args[0]
	if len(args) == 2 {
		versionArg = args[1]
		versionArgSlice := strings.Split(versionArg, ".")
		if len(versionArgSlice) == 3 && versionArgSlice[2] == "0" {
			versionArg = versionArgSlice[0] + "." + versionArgSlice[1]
		}
	}
}

func main() {

	switch actionArg {
	case "h", "help":
		cs.Print(usage())
	case "ls", "list":
		tmp("list")
	case "ls-remote":
		tmp("ls-remote")
	case "install":
		tmp("install")
	case "use":
		tmp("use")
	case "uninstall":
		tmp("uninstall")
	case "self-update":
		tmp("self-update")
	}
}

func isArgAllowed() bool {
	ok := true
	if len(os.Args) > 1 {
		ok = IsIn(allowedArgs, os.Args[1])
		if !ok {
			return false
		}
	}

	if len(os.Args) > 2 {
		ok = IsIn(allowedArgs, os.Args[1])
		if !ok {
			return false
		}
	}

	return ok
}

func IsIn(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func usage() string {
	msg := `
[light_red][bold]GoVM[reset]: Go版本管理器.[[red]` + version + `[reset]]

[light_green][underline]使用指南[reset]:
    [magenta]govm[reset] [light_gray]use <版本>              [yellow]安装并设置使用 <版本>[reset]
    [magenta]govm[reset] [light_gray]ls                      [yellow]list的别名[reset]
    [magenta]govm[reset] [light_gray]ls-remote               [yellow]远程版本列表 (包括 rc|beta 版本)[reset]
    [magenta]govm[reset] [light_gray]install <version>       [yellow]安装 <版本> (官方二进制或GOVM_REGISTRY环境变量)[reset]
    [magenta]govm[reset] [light_gray]uninstall <version>     [yellow]卸载<版本>[reset]
    [magenta]govm[reset] [light_gray]list                    [yellow]已经安装的版本(仅限GoVM管理的版本)[reset]
    [magenta]govm[reset] [light_gray]self-update             [yellow]GoVM自身升级[reset]
    [magenta]govm[reset] [light_gray]addpath                 [yellow]将GoVM相关程序加入环境变量[reset]
    [magenta]govm[reset] [light_gray]env                     [yellow]显示GoVM环境信息[reset]
    [magenta]govm[reset] [light_gray]help                    [yellow]显示此帮助信息[reset]
[light_green][underline]使用例子[reset]:
    [magenta]govm[reset] [light_gray]use 1.16                [yellow]使用1.16   版本的go[reset]
    [magenta]govm[reset] [light_gray]use 1.16.1              [yellow]使用1.16.1 版本的go[reset]
    [magenta]govm[reset] [light_gray]use 1.16rc1             [yellow]使用1.16rc1版本的go[reset]
    [magenta]govm[reset] [light_gray]use 1.16@latest         [yellow]使用1.16最新版本的go[reset]
    [magenta]govm[reset] [light_gray]use 1.16@dev-latest     [yellow]使用1.16最新版本的go, 包括rc和beta[reset]
    [magenta]govm[reset] [light_gray]use latest              [yellow]使用最新可用版本的go[reset]
    [magenta]govm[reset] [light_gray]use dev-latest          [yellow]使用最新可用版本的go,包括rc和beta[reset]
[light_green][underline]安装路径[reset]:(也可以简单使用 govm addpath 自动完成)
    [light_gray]将下面信息添加到你的~/.bashrc或~/.zshrc把[light_red][bold][underline]GoVM[reset][light_gray]加入环境变量[reset]
    [yellow][underline]export PATH="$HOME/.govm/current/bin:$HOME/.govm/bin:$PATH"
`
	return msg
}

func tmp(text string) {
	bar := progressbar.NewOptions(1000,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		// progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[light_cyan][1/3][reset] "+text+"..."),
		// https://github.com/mitchellh/colorstring/blob/d06e56a500db/colorstring.go#L124
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[light_red]=[reset]",
			SaucerHead:    "[light_red]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	for i := 0; i < 1000; i++ {
		bar.Add(1)
		time.Sleep(5 * time.Millisecond)
	}
}
