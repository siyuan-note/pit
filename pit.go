// Pit - 思源笔记更新函数
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
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/88250/gulu"
)

func ApplyUpdate(logger *gulu.Logger, updateZipPath, workingDir string) {
	if !gulu.File.IsExist(updateZipPath) {
		return
	}

	fis, err := ioutil.ReadDir(workingDir)
	if nil != err {
		logger.Fatalf("read working dir [%s] failed: %s", workingDir, err)
	}

	// 重命名为 .old 后缀

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "kernel") {
			kernel := filepath.Join(workingDir, fi.Name())
			if err = os.Rename(kernel, kernel+".old"); nil != err {
				logger.Errorf("rename kernel [%s] failed: %s", kernel, err)
			}
		}
	}

	asar := filepath.Join(workingDir, "app.asar")
	if err = os.Rename(asar, asar+".old"); nil != err {
		logger.Errorf("rename [app.asar] failed: %s", err)
	}

	appearance := filepath.Join(workingDir, "appearance")
	if err = os.Rename(appearance, appearance+".old"); nil != err {
		logger.Errorf("rename [appearance] failed: %s", err)
	}

	stage := filepath.Join(workingDir, "stage")
	if err = os.Rename(stage, stage+".old"); nil != err {
		logger.Errorf("rename [stage] failed: %s", err)
	}

	guide := filepath.Join(workingDir, "guide")
	if err = os.Rename(guide, guide+".old"); nil != err {
		logger.Errorf("rename [guide] failed: %s", err)
	}

	logger.Infof("unzipping update pack [from=%s, to=%s]", updateZipPath, workingDir)
	if err = gulu.Zip.Unzip(updateZipPath, workingDir); nil != err {
		logger.Errorf("unzip update pack failed: %s", err)
		return
	}
	logger.Infof("unzipped update pack")

	if !gulu.OS.IsWindows() {
		fis, _ = ioutil.ReadDir(workingDir)
		for _, fi := range fis {
			if strings.HasPrefix(fi.Name(), "kernel") {
				kernel := filepath.Join(workingDir, fi.Name())
				exec.Command("chmod", "+x", kernel).CombinedOutput()
				logger.Infof("chmod +x " + kernel)
			}
		}
	}

	if err = os.RemoveAll(updateZipPath); nil != err {
		logger.Errorf("remove update pack [%s] failed: %s", updateZipPath, err)
		return
	}
}
