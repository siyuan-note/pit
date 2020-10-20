// SiYuan - 源于思考，饮水思源
// Copyright (c) 2020-present, ld246.com
//
// 本文件属于思源笔记源码的一部分，云南链滴科技有限公司版权所有。
//

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/88250/gulu"
)

var (
	HomeDir, _    = gulu.OS.Home()
	WorkingDir, _ = os.Getwd()
	ConfDir       = filepath.Join(HomeDir, ".siyuan")
	LogPath       = filepath.Join(ConfDir, "siyuan-update.log")
)

var (
	Logger  *gulu.Logger
	logFile *os.File
)

func init() {
	if !gulu.File.IsExist(ConfDir) {
		if err := os.Mkdir(ConfDir, 0755); nil != err && !os.IsExist(err) {
			Logger.Fatalf("create conf folder [%s] failed: %s", ConfDir, err)
		}
	}

	if gulu.File.IsExist(LogPath) {
		if size := gulu.File.GetFileSize(LogPath); 1024*1024*2 <= size {
			// 日志文件大于 2M 的话删了重建
			if err := os.Remove(LogPath); nil != err {
				fmt.Errorf("remove log file [%s] failed: %s", LogPath, err)
				os.Exit(-2)
			}
		}
	}

	var err error
	logFile, err = os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if nil != err {
		fmt.Errorf("create log file [%s] failed: %s", LogPath, err)
		os.Exit(-2)
	}

	gulu.Log.SetLevel("trace")
	Logger = gulu.Log.NewLogger(io.MultiWriter(logFile, os.Stdout))
}

func main() {
	defer logFile.Close()

	wd := flag.String("wd", WorkingDir, "working directory")
	flag.Parse()

	WorkingDir = *wd

	syTempFolder := filepath.Join(os.TempDir(), "siyuan")
	p := filepath.Join(syTempFolder, "update.zip")
	if !gulu.File.IsExist(p) {
		return
	}

	Logger.Infof("unzipping update pack [from=%s, to=%s]", p, WorkingDir)
	if err := gulu.Zip.Unzip(p, WorkingDir); nil != err {
		Logger.Errorf("unzip update pack failed: %s", err)
		return
	}
	Logger.Infof("unzipped update pack")

	kernel := filepath.Join(WorkingDir, "kernel")
	if gulu.OS.IsWindows() {
		kernel += ".exe"
	} else {
		if gulu.OS.IsLinux() {
			kernel += "-linux"
		} else if gulu.OS.IsDarwin() {
			kernel += "-darwin"
		}
		exec.Command("chmod", "+x", kernel)
	}

	if err := os.RemoveAll(p); nil != err {
		Logger.Errorf("remove update pack [%s] failed: %s", p, err)
		return
	}
}
