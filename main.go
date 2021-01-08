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
	// 删除内核程序
	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "kernel") {
			kernel := filepath.Join(WorkingDir, fi.Name())
			if err = os.RemoveAll(kernel); nil != err {
				Logger.Errorf("remove kernel [%s] failed: %s", kernel, err)
			} else {
				Logger.Infof("removed old kernel [%s]", kernel)
			}
		}
	}

	appearance := filepath.Join(WorkingDir, "appearance")
	stage := filepath.Join(WorkingDir, "stage")
	guide := filepath.Join(WorkingDir, "guide")
	asar := filepath.Join(WorkingDir, "app.asar")

	// TODO: 考虑回滚机制（比如先整体移动到临时目录，后续解压等操作如果失败再移动回来）

	for cnt := 0; cnt < 7; cnt++ { // 重试执行
		if err = os.RemoveAll(asar); nil != err { // 删除 app.asar
			Logger.Errorf("remove [app.asar] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			Logger.Infof("removed [app.asar]", asar)
		}

		if err = os.RemoveAll(appearance); nil != err { // 删除 appearance 文件夹
			Logger.Errorf("remove [appearance] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			Logger.Infof("removed [appearance]")
		}
		if err = os.RemoveAll(stage); nil != err { // 删除 stage 文件夹
			Logger.Errorf("remove [stage] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			Logger.Infof("removed [stage]")
		}
		if err = os.RemoveAll(guide); nil != err { // 删除 guide 文件夹
			Logger.Errorf("remove [guide] failed: %s", err)
			time.Sleep(50 * time.Millisecond)
			continue
		} else {
			Logger.Infof("removed [guide]")
		}

		break
	}

	Logger.Infof("unzipping update pack [from=%s, to=%s]", UpdateZipPath, WorkingDir)
	if err = gulu.Zip.Unzip(UpdateZipPath, WorkingDir); nil != err {
		Logger.Errorf("unzip update pack failed: %s", err)
		return
	}
	Logger.Infof("unzipped update pack")

	if !gulu.OS.IsWindows() {
		// 赋予执行权限

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
