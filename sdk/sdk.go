package sdk

import (
	"fmt"
	"os"

	"github.com/milin2436/BaiduPCS-Go/baidupcs"
	"github.com/milin2436/BaiduPCS-Go/baidupcs/pcserror"
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
func RunDownload(paths []string, options map[string]string) {
	op := new(pcscommand.DownloadOptions)
	if options != nil {
		op.SaveTo = options["saveto"]
	}
	pcscommand.RunDownload(paths, op)
}

func RunShareTransfer(params []string, opt map[string]string) ([]string, error) {
	op := new(baidupcs.TransferOption)
	if opt != nil {
		op.SaveTo = opt["saveto"]
		op.DnSaveTo = opt["dnsaveto"]
		if opt["download"] == "true" {
			op.Download = true
		}
	}
	li, err := pcscommand.RunShareTransferForSdk(params, op)
	if err != nil && op.SaveTo != "" {
		RunRemove(op.SaveTo)
	}
	return li, err
}

func RunMkdir(path string) pcserror.Error {
	pcs := pcscommand.GetBaiduPCS()
	return pcs.Mkdir(path)
}
func RunRemove(paths ...string) {
	pcscommand.RunRemove(paths...)
}
