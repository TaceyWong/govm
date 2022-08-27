package govm

import (
	"testing"
)

func TestRegexGroup(t *testing.T) {
	l := []string{
		`<td class="link"><a href="go1.10.1.linux-386.tar.gz">go1.10.1.linux-386.tar.gz</a></td>`,
		`<td class="link"><a href="go1.10.1rc.windows-amd64.zip">go1.10.1rc.windows-amd64.zip</a></td>`,
		`<td class="link"><a href="go1.10.1rc1.freebsd-386.tar.gz">go1.10.1rc1.freebsd-386.tar.gz</a></td>`,
		`<td class="link"><a href="go1.8.linux-darwin.tar.gz">go1.8.linux-darwin.tar.gz</a></td>`,
		`<td class="link"><a href="go1.8rc.unkown-386.tar.gz">go1.8rc.unkown-386.tar.gz</a></td>`,
	}
	p := `<td class="link"><a href="(?P<link>\w.*?)">go(?P<version>.*)\.(?P<os>[a-z]+)-(?P<arch>[a-z]*[0-9]*)\.(?P<extension>[\w.]+)</a>`
	for _, s := range l {
		println("待解析匹配文本", s)
		println("------------------------------------------------------------")
		result := RegexGroup(p, s)
		println("URLBase:", result["link"])
		println("操作系统:", result["os"])
		println("CPU Arch:", result["arch"])
		println("压缩包后缀:", result["extension"])
		println("版本:", result["version"])
		println("------------------------------------------------------------")

	}

}
