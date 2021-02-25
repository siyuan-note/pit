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
	"path/filepath"
	"strings"
	"sync"

	"github.com/88250/gulu"
)

var applyUpdateLock = sync.Mutex{}
var applied bool

func ApplyUpdate(updateZipPath, workingDir string) error {
	applyUpdateLock.Lock()
	defer applyUpdateLock.Unlock()

	if applied || !gulu.File.IsExist(updateZipPath) {
		return nil
	}
	applied = true

	fis, err := ioutil.ReadDir(workingDir)
	if nil != err {
		return errors.New(fmt.Sprintf("read working dir [%s] failed: %s", workingDir, err))
	}

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "kernel") && !strings.HasSuffix(fi.Name(), ".old") {
			kernel := filepath.Join(workingDir, fi.Name())
			os.Rename(kernel, kernel+".old")
		}
	}

	appearance := filepath.Join(workingDir, "appearance")
	if err = os.Rename(appearance, appearance+".old"); nil != err {
		return err
	}
	stage := filepath.Join(workingDir, "stage")
	if err = os.Rename(stage, stage+".old"); nil != err {
		return err
	}
	guide := filepath.Join(workingDir, "guide")
	if err = os.Rename(guide, guide+".old"); nil != err {
		return err
	}

	if err = gulu.Zip.Unzip(updateZipPath, workingDir); nil != err {
		return errors.New(fmt.Sprintf("unzip update pack failed: %s", err))
	}
	return nil
}

func Rollback(workingDir string) error {
	fis, err := ioutil.ReadDir(workingDir)
	if nil != err {
		return errors.New(fmt.Sprintf("read working dir [%s] failed: %s", workingDir, err))
	}

	for _, fi := range fis {
		if strings.HasSuffix(fi.Name(), ".old") {
			file := filepath.Join(workingDir, fi.Name())
			newFile := strings.TrimSuffix(file, ".old")
			if err = os.RemoveAll(newFile); nil != err {
				return errors.New(fmt.Sprintf("remove [%s] failed: %s", newFile, err))
			}
			if err = os.Rename(file, newFile); nil != err {
				return errors.New(fmt.Sprintf("rename [%s] to [%s] failed: %s", file, newFile, err))
			}
		}
	}

	os.RemoveAll(filepath.Join(workingDir, "update.asar"))
	return nil
}
