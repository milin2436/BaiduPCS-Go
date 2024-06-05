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
	op["saveto"] = "/d/test001"
	op["dnsaveto"] = "/home/super/Downloads/papa"
	op["download"] = "true"

	sdk.RunMkdir(op["saveto"])
	paths, err := sdk.RunShareTransfer([]string{"https://pan.baidu.com/s/1UZFN-tIq-LUoBPc7hUCZYg?pwd=8962"}, op)
	fmt.Println("err #", err)
	fmt.Println("paths #", paths)
}
func tRm() {
	defer sdk.SdkClose()
	sdk.RunRemove("/d/UEi0SbEFkfE")
}
func main() {
	tTr()
}
