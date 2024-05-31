package main

import (
	"fmt"

	"github.com/milin2436/BaiduPCS-Go/sdk"
)

func init() {
	sdk.SdkInit()
}

func tdn() {
	defer sdk.SdkClose()
	op := map[string]string{}
	op["saveto"] = "/home/super/Documents"
	sdk.RunDownload([]string{"/update/view/mojin/mojinJ.hdy"}, op)
}
func tTr() {
	defer sdk.SdkClose()
	op := map[string]string{}
	op["saveto"] = "/d"
	op["download"] = "true"
	err := sdk.RunShareTransfer([]string{"https://pan.baidu.com/s/1owdUpAEUq8rPJr4TdTUZtA?pwd=ugg4"}, op)
	fmt.Println("err #", err)
}
func main() {
	tTr()
}
