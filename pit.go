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
	"os"
	"path/filepath"
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

	kernel := filepath.Join(workingDir, "kernel")
	if err := os.Rename(kernel, kernel+".old"); nil != err {
		return err
	}
	appearance := filepath.Join(workingDir, "appearance")
	if err := os.Rename(appearance, appearance+".old"); nil != err {
		return err
	}
	stage := filepath.Join(workingDir, "stage")
	if err := os.Rename(stage, stage+".old"); nil != err {
		return err
	}
	guide := filepath.Join(workingDir, "guide")
	if err := os.Rename(guide, guide+".old"); nil != err {
		return err
	}
	app := filepath.Join(workingDir, "app")
	if err := os.Rename(app, app+".old"); nil != err {
		return err
	}

	if err := gulu.Zip.Unzip(updateZipPath, workingDir); nil != err {
		return errors.New(fmt.Sprintf("unzip update pack failed: %s", err))
	}

	os.RemoveAll(updateZipPath)
	return nil
}
