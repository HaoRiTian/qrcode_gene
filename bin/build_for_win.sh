#!/bin/sh

Script_Dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
Work_Dir="$( cd ${Script_Dir}/.. && pwd )"
cd ${Work_Dir}

Version=$( git describe --abbrev=0 --tags )
if [ -z ${Version} ];then
  Version="Unknown"
fi

go generate
GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui -w -s" -o build/qrcode_gene-${Version}.exe