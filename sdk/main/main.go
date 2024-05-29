package main

import (
	"fmt"
	"os"

	"github.com/milin2436/BaiduPCS-Go/internal/pcscommand"
	"github.com/milin2436/BaiduPCS-Go/internal/pcsconfig"
	"github.com/milin2436/BaiduPCS-Go/pcsutil"
)

func init() {
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

func main() {
	defer pcsconfig.Config.Close()
	q, u, err := pcscommand.GetBaiduPCS().QuotaInfo()

	if err != nil {
		fmt.Println("err = ", err)
		return
	}
	fmt.Printf("%d\n", q)
	fmt.Printf("%d\n", u)
}
