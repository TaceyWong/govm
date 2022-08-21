#!/bin/sh

GOVM_BIN_DIR=$HOME/.govm/bin
mkdir -p $GOVM_BIN_DIR

GOVM_ARCH_BIN=''

THISOS=$(uname -s)
ARCH=$(uname -m)

case $THISOS in
   Linux*)
      case $ARCH in
        arm64)
          GOVM_ARCH_BIN="govm-linux-arm64"
          ;;
        aarch64)
          GOVM_ARCH_BIN="govm-linux-arm64"
          ;;
        *)
          GOVM_ARCH_BIN="govm-linux-amd64"
          ;;
      esac
      ;;
   Darwin*)
      case $ARCH in
        arm64)
          GOVM_ARCH_BIN="govm-darwin-arm64"
          ;;
        *)
          GOVM_ARCH_BIN="govm-darwin-amd64"
          ;;
      esac
      ;;
   Windows*)
      GOVM_ARCH_BIN="govm-windows-amd64.exe"
      ;;
esac

if [ -z "$GOVM_VERSION" ]
then
      GOVM_VERSION=master
      echo "使用最新版本的Govm\n"
else
      echo "使用$GOVM_VERSION版本的GoVM\n"
fi

curl -kLs https://github.com/TaceyWong/govm/releases/latest/download/$GOVM_ARCH_BIN -o $GOVM_BIN_DIR/govm

chmod +x $GOVM_BIN_DIR/govm

echo "成功安装到: $GOVM_BIN_DIR/govm"

echo "============================"
$GOVM_BIN_DIR/govm help
echo "============================"

echo
echo "***请将以下内容手动添加到你的环境变量(~/.bashrc或~/.zshrc之类)***"
echo
echo 'export PATH="$HOME/.govm/current/bin:$HOME/.govm/bin:$PATH"'
echo