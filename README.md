# govm

GoVM: A  Go Version Manager. 

**快速配置**: 一条命令安装Go并管理版本.

**权限自由**: 无需root或sudo权限

**支持平台**: 支持Linux、FreeBSD、Mac.

**操作灵活**: 自由管理切换不同版本的Go


```shell
GoVM: Go版本管理器.[0.0.1.dev]

使用指南:
    govm use <版本>              安装并设置使用 <版本>
    govm ls                      list的别名
    govm ls-remote               远程版本列表 (包括 rc|beta 版本)
    govm install <version>       安装 <版本> (官方二进制或GOVM_REGISTRY环境变量)
    govm uninstall <version>     卸载<版本>
    govm list                    已经安装的版本(仅限GoVM管理的版本)
    govm self-update             GoVM自身升级
    govm addpath                 将GoVM相关程序加入环境变量
    govm env                     显示GoVM环境信息
    govm help                    显示此帮助信息
使用例子:
    govm use 1.16                使用1.16   版本的go
    govm use 1.16.1              使用1.16.1 版本的go
    govm use 1.16rc1             使用1.16rc1版本的go
    govm use 1.16@latest         使用1.16最新版本的go
    govm use 1.16@dev-latest     使用1.16最新版本的go, 包括rc和beta
    govm use latest              使用最新可用版本的go
    govm use dev-latest          使用最新可用版本的go,包括rc和beta
安装路径:(也可以简单使用 govm addpath 自动完成)
    将下面信息添加到你的~/.bashrc或~/.zshrc把GoVM加入环境变量
    export PATH="$HOME/.govm/current/bin:$HOME/.govm/bin:$PATH"
```

## 安装和更新

使用`curl`安装

```shell
curl -sLk https://git.io/gobrew | sh -
```
或使用`go`安装
```shell
go install github.com/kevincobain2000/gobrew/cmd/gobrew@latest
```
在你的shell配置文件中将`govm`添加到环境变量(.bashrc 或 .zshrc).
```shell
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"
```
重载配置，一切完成！