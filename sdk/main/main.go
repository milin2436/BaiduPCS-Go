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
	op["saveto"] = "/d/fsf2"
	op["dnsaveto"] = "/home/super/Downloads/papa"
	op["download"] = "true"

	sdk.RunMkdir(op["saveto"])
	paths, err := sdk.RunShareTransfer([]string{"https://pan.baidu.com/s/1GjVLYb-VIEShjouLllEnbw?pwd=36pt"}, op)
	fmt.Println("err #", err)
	fmt.Println("paths #", paths)
}
func tRm() {
	defer sdk.SdkClose()
	sdk.RunRemove("/d/UEi0SbEFkfE")
}
func main() {
	tRm()
}
