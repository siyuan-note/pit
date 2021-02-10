// Pit - 思源笔记更新程序
// Copyright (c) 2020-present, ld246.com
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	time.Sleep(500 * time.Millisecond) // 稍微等待一下思源内核进程退出

	fis, err := ioutil.ReadDir(WorkingDir)
	if nil != err {
		Logger.Fatalf("read working dir [%s] failed: %s", WorkingDir, err)
	}

	// 重命名为 .old 后缀

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "kernel") {
			kernel := filepath.Join(WorkingDir, fi.Name())
			if err = os.Rename(kernel, kernel+".old"); nil != err {
				Logger.Errorf("rename kernel [%s] failed: %s", kernel, err)
			}
		}
	}

	asar := filepath.Join(WorkingDir, "app.asar")
	if err = os.Rename(asar, asar+".old"); nil != err {
		Logger.Errorf("rename [app.asar] failed: %s", err)
	}

	appearance := filepath.Join(WorkingDir, "appearance")
	if err = os.Rename(appearance, appearance+".old"); nil != err {
		Logger.Errorf("rename [appearance] failed: %s", err)
	}

	stage := filepath.Join(WorkingDir, "stage")
	if err = os.Rename(stage, stage+".old"); nil != err {
		Logger.Errorf("rename [stage] failed: %s", err)
	}

	guide := filepath.Join(WorkingDir, "guide")
	if err = os.Rename(guide, guide+".old"); nil != err {
		Logger.Errorf("rename [guide] failed: %s", err)
	}

	Logger.Infof("unzipping update pack [from=%s, to=%s]", UpdateZipPath, WorkingDir)
	if err = gulu.Zip.Unzip(UpdateZipPath, WorkingDir); nil != err {
		Logger.Errorf("unzip update pack failed: %s", err)
		return
	}
	Logger.Infof("unzipped update pack")

	if !gulu.OS.IsWindows() {
		fis, _ = ioutil.ReadDir(WorkingDir)
		for _, fi := range fis {
			if strings.HasPrefix(fi.Name(), "kernel") {
				kernel := filepath.Join(WorkingDir, fi.Name())
				exec.Command("chmod", "+x", kernel).CombinedOutput()
				Logger.Infof("chmod +x " + kernel)
			}
		}
	}

	if err = os.RemoveAll(UpdateZipPath); nil != err {
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
