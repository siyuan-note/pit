// Pit - 思源笔记更新程序
// Copyright (c) 2020-present, ld246.com
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package pit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/88250/gulu"
)

func ApplyUpdate(updateZipPath, workingDir string) error {
	if !gulu.File.IsExist(updateZipPath) {
		return nil
	}

	fis, err := ioutil.ReadDir(workingDir)
	if nil != err {
		return errors.New(fmt.Sprintf("read working dir [%s] failed: %s", workingDir, err))
	}

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "kernel") {
			kernel := filepath.Join(workingDir, fi.Name())
			if err = os.Rename(kernel, kernel+".old"); nil != err {
				return errors.New(fmt.Sprintf("rename kernel [%s] failed: %s", kernel, err))
			}
		}
	}

	appearance := filepath.Join(workingDir, "appearance")
	if err = os.Rename(appearance, appearance+".old"); nil != err {
		return errors.New(fmt.Sprintf("rename [appearance] failed: %s", err))
	}

	stage := filepath.Join(workingDir, "stage")
	if err = os.Rename(stage, stage+".old"); nil != err {
		return errors.New(fmt.Sprintf("rename [stage] failed: %s", err))
	}

	guide := filepath.Join(workingDir, "guide")
	if err = os.Rename(guide, guide+".old"); nil != err {
		return errors.New(fmt.Sprintf("rename [guide] failed: %s", err))
	}

	if err = gulu.Zip.Unzip(updateZipPath, workingDir); nil != err {
		return errors.New(fmt.Sprintf("unzip update pack failed: %s", err))
	}

	if !gulu.OS.IsWindows() {
		fis, _ = ioutil.ReadDir(workingDir)
		for _, fi := range fis {
			if strings.HasPrefix(fi.Name(), "kernel") {
				kernel := filepath.Join(workingDir, fi.Name())
				exec.Command("chmod", "+x", kernel).CombinedOutput()
			}
		}
	}

	if err = os.RemoveAll(updateZipPath); nil != err {
		return errors.New(fmt.Sprintf("remove update pack [%s] failed: %s", updateZipPath, err))
	}
	return nil
}
