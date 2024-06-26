package pcscommand

import (
	"encoding/base64"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/milin2436/BaiduPCS-Go/baidupcs"
)

// RunShareTransfer 执行分享链接转存到网盘
func RunShareTransfer(params []string, opt *baidupcs.TransferOption) {
	var link string
	var extracode string
	if len(params) == 1 {
		link = params[0]
		if strings.Contains(link, "bdlink=") || !strings.Contains(link, "pan.baidu.com/") {
			RunRapidTransfer(link, opt.Rname)
			//fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "秒传已不再被支持")
			return
		}
		extracode = "none"
		if strings.Contains(link, "?pwd=") {
			extracode = strings.Split(link, "?pwd=")[1]
			link = strings.Split(link, "?pwd=")[0]
		}
	} else if len(params) == 2 {
		link = params[0]
		extracode = params[1]
	}
	if link[len(link)-1:] == "/" {
		link = link[0 : len(link)-1]
	}
	featurestrs := strings.Split(link, "/")
	featurestr := featurestrs[len(featurestrs)-1]
	if strings.Contains(featurestr, "init?") {
		featurestr = "1" + strings.Split(featurestr, "=")[1]
	}
	if len(featurestr) > 23 || featurestr[0:1] != "1" || len(extracode) != 4 {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接地址或提取码非法")
		return
	}
	pcs := GetBaiduPCS()
	tokens := pcs.AccessSharePage(featurestr, true)
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}
	// pcs.UpdatePCSCookies(true)
	var vefiryurl string
	var randsk string
	featuremap := make(map[string]string)
	featuremap["bdstoken"] = tokens["bdstoken"]
	featuremap["surl"] = featurestr[1:len(featurestr)]
	if extracode != "none" {

		vefiryurl = pcs.GenerateShareQueryURL("verify", featuremap).String()
		res := pcs.PostShareQuery(vefiryurl, link, map[string]string{
			"pwd":       extracode,
			"vcode":     "",
			"vcode_str": "",
		})
		if res["ErrMsg"] != "0" {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"])
			return
		}
		randsk = res["randsk"]
	}
	pcs.UpdatePCSCookies(true)

	tokens = pcs.AccessSharePage(featurestr, false)
	tokens["randsk"] = randsk
	if tokens["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		return
	}
	metajsonstr := tokens["metajson"]
	trans_metas := pcs.ExtractShareInfo(metajsonstr)

	if trans_metas["ErrMsg"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, trans_metas["ErrMsg"])
		return
	}
	trans_metas["path"] = GetActiveUser().Workdir
	if trans_metas["item_num"] != "1" && opt.Collect {
		trans_metas["filename"] += "等文件"
		trans_metas["path"] = path.Join(GetActiveUser().Workdir, trans_metas["filename"])
		pcs.Mkdir(trans_metas["path"])
	}
	trans_metas["referer"] = "https://pan.baidu.com/s/" + featurestr
	pcs.UpdatePCSCookies(true)
	resp := pcs.GenerateRequestQuery("POST", trans_metas)
	if resp["ErrNo"] != "0" {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, resp["ErrMsg"])
		if resp["ErrNo"] == "4" {
			trans_metas["shorturl"] = featurestr
			pcs.SuperTransfer(trans_metas, resp["limit"]) // 试验性功能, 当前未启用
		}
		return
	}
	if opt.Collect {
		resp["filename"] = trans_metas["filename"]
	}
	fmt.Printf("%s成功, 保存了%s到当前目录\n", baidupcs.OperationShareFileSavetoLocal, resp["filename"])
	if opt.Download {
		fmt.Println("即将开始下载")
		paths := strings.Split(resp["filenames"], ",")
		paths = paths[0 : len(paths)-1]
		fmt.Println("paths = ", paths)
		RunDownload(paths, nil)
	}
}

// RunRapidTransfer 执行秒传链接解析及保存
func RunRapidTransfer(link string, rnameOpt ...bool) {
	if strings.Contains(link, "bdlink=") || strings.Contains(link, "bdpan://") {
		r, _ := regexp.Compile(`(bdlink=|bdpan://)([^\s]+)`)
		link1 := r.FindStringSubmatch(link)[2]
		decodeBytes, err := base64.StdEncoding.DecodeString(link1)
		if err != nil {
			fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
			return
		}
		link = string(decodeBytes)
	}
	rname := false
	if len(rnameOpt) > 0 {
		rname = rnameOpt[0]
	}
	link = strings.TrimSpace(link)
	substrs := strings.SplitN(link, "#", 4)
	if len(substrs) == 4 {
		md5, slicemd5 := substrs[0], substrs[1]
		size, _ := strconv.ParseInt(substrs[2], 10, 64)
		filename := path.Join(GetActiveUser().Workdir, randReplaceStr(substrs[3], rname))
		RunRapidUpload(filename, md5, slicemd5, size)
	} else if len(substrs) == 3 {
		md5 := substrs[0]
		size, _ := strconv.ParseInt(substrs[1], 10, 64)
		filename := path.Join(GetActiveUser().Workdir, randReplaceStr(substrs[2], rname))
		RunRapidUpload(filename, md5, "", size)
	} else {
		fmt.Printf("%s失败: %s\n", baidupcs.OperationRapidLinkSavetoLocal, "秒传链接格式错误")
	}
	return
}

// RunShareTransfer sdk 执行分享链接转存到网盘
func RunShareTransferForSdk(params []string, opt *baidupcs.TransferOption) ([]string, error) {
	var link string
	var extracode string
	var msg string
	if len(params) == 1 {
		link = params[0]
		if strings.Contains(link, "bdlink=") || !strings.Contains(link, "pan.baidu.com/") {
			//RunRapidTransfer(link, opt.Rname)
			msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "秒传已不再被支持")
			fmt.Print(msg)
			return nil, errors.New(msg)
		}
		extracode = "none"
		if strings.Contains(link, "?pwd=") {
			extracode = strings.Split(link, "?pwd=")[1]
			link = strings.Split(link, "?pwd=")[0]
		}
	} else if len(params) == 2 {
		link = params[0]
		extracode = params[1]
	}
	if link[len(link)-1:] == "/" {
		link = link[0 : len(link)-1]
	}
	featurestrs := strings.Split(link, "/")
	featurestr := featurestrs[len(featurestrs)-1]
	if strings.Contains(featurestr, "init?") {
		featurestr = "1" + strings.Split(featurestr, "=")[1]
	}
	if len(featurestr) > 23 || featurestr[0:1] != "1" || len(extracode) != 4 {
		msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, "链接地址或提取码非法")
		fmt.Print(msg)
		return nil, errors.New(msg)
	}
	pcs := GetBaiduPCS()
	tokens := pcs.AccessSharePage(featurestr, true)
	if tokens["ErrMsg"] != "0" {
		msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		fmt.Print(msg)
		return nil, errors.New(msg)
	}
	// pcs.UpdatePCSCookies(true)
	var vefiryurl string
	var randsk string
	featuremap := make(map[string]string)
	featuremap["bdstoken"] = tokens["bdstoken"]
	featuremap["surl"] = featurestr[1:len(featurestr)]
	if extracode != "none" {

		vefiryurl = pcs.GenerateShareQueryURL("verify", featuremap).String()
		res := pcs.PostShareQuery(vefiryurl, link, map[string]string{
			"pwd":       extracode,
			"vcode":     "",
			"vcode_str": "",
		})
		if res["ErrMsg"] != "0" {
			msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, res["ErrMsg"])
			fmt.Print(msg)
			return nil, errors.New(msg)
		}
		randsk = res["randsk"]
	}
	pcs.UpdatePCSCookies(true)

	tokens = pcs.AccessSharePage(featurestr, false)
	tokens["randsk"] = randsk
	if tokens["ErrMsg"] != "0" {
		msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, tokens["ErrMsg"])
		fmt.Print(msg)
		return nil, errors.New(msg)
	}
	metajsonstr := tokens["metajson"]
	trans_metas := pcs.ExtractShareInfo(metajsonstr)

	if trans_metas["ErrMsg"] != "0" {
		msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, trans_metas["ErrMsg"])
		fmt.Print(msg)
		return nil, errors.New(msg)
	}
	if opt.SaveTo != "" {
		trans_metas["path"] = opt.SaveTo
		/*
			perr := pcs.Mkdir(opt.SaveTo)
			if perr != nil {
				return perr
			}
		*/
	} else {
		trans_metas["path"] = GetActiveUser().Workdir
	}
	if trans_metas["item_num"] != "1" && opt.Collect {
		trans_metas["filename"] += "等文件"
		trans_metas["path"] = path.Join(GetActiveUser().Workdir, trans_metas["filename"])
		pcs.Mkdir(trans_metas["path"])
	}
	trans_metas["referer"] = "https://pan.baidu.com/s/" + featurestr
	pcs.UpdatePCSCookies(true)
	fileList := map[string][]string{}
	resp := pcs.TransferShareRequestQuery("POST", trans_metas, fileList)
	if resp["ErrNo"] != "0" {
		msg = fmt.Sprintf("%s失败: %s\n", baidupcs.OperationShareFileSavetoLocal, resp["ErrMsg"])
		if resp["ErrNo"] == "4" {
			trans_metas["shorturl"] = featurestr
			pcs.SuperTransfer(trans_metas, resp["limit"]) // 试验性功能, 当前未启用
		}
		fmt.Print(msg)
		return nil, errors.New(msg)
	}
	if opt.Collect {
		resp["filename"] = trans_metas["filename"]
	}
	fmt.Println("list = ", fileList["list"])
	fmt.Printf("%s成功, 保存了%s到当前目录\n", baidupcs.OperationShareFileSavetoLocal, resp["filename"])

	li := fileList["list"]
	if len(li) > 0 {
		return li, RunCheckDownloadTotalSize(li, sizeLimit)
	}
	return li, nil
}

const sizeLimit = int64(1024 * 1024 * 1024 * 10)
