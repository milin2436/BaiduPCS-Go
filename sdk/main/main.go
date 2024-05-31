package main

import (
	"github.com/milin2436/BaiduPCS-Go/sdk"
)

func init() {
	sdk.SdkInit()
}

func main() {
	defer sdk.SdkClose()
	sdk.RunDownload([]string{"mojinJ.hdy"})
}
