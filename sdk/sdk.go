package sdk

import (
	"fmt"
	"os"

	"github.com/milin2436/BaiduPCS-Go/internal/pcscommand"
	"github.com/milin2436/BaiduPCS-Go/internal/pcsconfig"
	"github.com/milin2436/BaiduPCS-Go/pcsutil"
)

func SdkInit() {
	pcsutil.ChWorkDir()
	err := pcsconfig.Config.Init()
	switch err {
	case nil:
	case pcsconfig.ErrConfigFileNoPermission, pcsconfig.ErrConfigContentsParseError:
		fmt.Fprintf(os.Stderr, "FATAL ERROR: config file error: %s\n", err)
		os.Exit(1)
	default:
		fmt.Printf("WARNING: config init error: %s\n", err)
	}
}

func SdkClose() {
	//TODO
	pcsconfig.Config.Close()
}
func RunDownload(paths []string) {
	pcscommand.RunDownload(paths, nil)
}
