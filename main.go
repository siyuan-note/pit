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
	"time"

	"github.com/88250/gulu"
)

var (
	HomeDir, _       = gulu.OS.Home()
	WorkingDir, _    = os.Getwd()
	ConfDir          = filepath.Join(HomeDir, ".siyuan")
	LogPath          = filepath.Join(ConfDir, "siyuan-update.log")
	SiYuanTempFolder = filepath.Join(os.TempDir(), "siyuan")
	UpdateZipPath    = filepath.Join(SiYuanTempFolder, "update.zip")
)

var (
	Logger  *gulu.Logger
	logFile *os.File
)

func main() {
	defer logFile.Close()

	confPath := flag.String("conf", "", "dir path of conf dir (.siyuan/), default to ~/siyuan/")
	wdPath := flag.String("wd", WorkingDir, "working directory of SiYuan")
	packPath := flag.String("pack", UpdateZipPath, "update package path, default to $temp/siyuan/update.zip")
	flag.Parse()

	if "" != *confPath {
		ConfDir = *confPath
	}
	if "" != *wdPath {
		WorkingDir = *wdPath
	}
	if "" != *packPath {
		UpdateZipPath = *packPath
	}

	initConfLog()

	if !gulu.File.IsExist(UpdateZipPath) {
		return
	}

	time.Sleep(500 * time.Millisecond)
	kernel := filepath.Join(WorkingDir, "kernel")
	if gulu.OS.IsWindows() {
		kernel += ".exe"
	} else {
		if gulu.OS.IsLinux() {
			kernel += "-linux"
		} else if gulu.OS.IsDarwin() {
			kernel += "-darwin"
		}
	}
	appearance := filepath.Join(WorkingDir, "appearance")
	stage := filepath.Join(WorkingDir, "stage")
	guide := filepath.Join(WorkingDir, "guide")
	asar := filepath.Join(WorkingDir, "app.asar")

	for cnt := 0; cnt < 7; cnt++ {
		if err := os.RemoveAll(kernel); nil != err {
			Logger.Errorf("remove [kernel] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if err := os.RemoveAll(asar); nil != err {
			Logger.Errorf("remove [app.asar] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		}

		if err := os.RemoveAll(appearance); nil != err {
			Logger.Errorf("remove [appearance] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if err := os.RemoveAll(stage); nil != err {
			Logger.Errorf("remove [stage] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		}
		if err := os.RemoveAll(guide); nil != err {
			Logger.Errorf("remove [guide] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		}

		break
	}

	Logger.Infof("unzipping update pack [from=%s, to=%s]", UpdateZipPath, WorkingDir)
	if err := gulu.Zip.Unzip(UpdateZipPath, WorkingDir); nil != err {
		Logger.Errorf("unzip update pack failed: %s", err)
		return
	}
	Logger.Infof("unzipped update pack")
	if gulu.OS.IsWindows() {
		unzippedKernel := filepath.Join(WorkingDir, "kernel-win.exe")
		if err := os.Rename(unzippedKernel, kernel); nil != err {
			Logger.Errorf("rename kernel [from=%s, to=%s] failed: %s", unzippedKernel, kernel, err)
			return
		}
	} else {
		exec.Command("chmod", "+x", kernel).CombinedOutput()
	}
	if err := os.RemoveAll(UpdateZipPath); nil != err {
		Logger.Errorf("remove update pack [%s] failed: %s", UpdateZipPath, err)
		return
	}
}

func initConfLog() {
	if !gulu.File.IsExist(ConfDir) {
		if err := os.Mkdir(ConfDir, 0755); nil != err && !os.IsExist(err) {
			Logger.Fatalf("create conf folder [%s] failed: %s", ConfDir, err)
		}
	}

	if gulu.File.IsExist(LogPath) {
		if size := gulu.File.GetFileSize(LogPath); 1024*1024*1 <= size {
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
