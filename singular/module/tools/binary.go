/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tools

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/Huawei/containerops/common/utils"
)

const (
	DefaultSSHUser = "root"
	DefaultSSHPort = 22
)

//DownloadComponents is download component binary file to the host.
//If the src is URL, execute curl command in the host.
//If the src is local file, execute scp command upload to the host.
func DownloadComponent(src, dest, host, private, user string, stdout io.Writer) (string, error) {
	var cmd string

	if u, err := url.Parse(src); err != nil {
		return "", fmt.Errorf("Invalid src format, neither URL or local path.")
	} else {
		if u.Scheme == "" {
			if utils.IsFileExist(src) == true {
				if cmd, err = uploadBinary(src, dest, host, private, user); err != nil {
					return "", err
				}
			} else {
				return "", fmt.Errorf("The file not exist")
			}
		} else {
			if cmd, err = downloadBinary(src, dest, host, private, user, stdout); err != nil {
				return "", err
			}
		}

	}

	return cmd, nil
}

//downloadBinary exec curl command download binary in the node.
func downloadBinary(src, dest, host, private, user string, stdout io.Writer) (string, error) {
	cmd := fmt.Sprintf("curl %s -o %s", src, dest)
	if err := utils.SSHCommand(user, private, host, DefaultSSHPort, cmd, stdout, os.Stderr); err != nil {
		return cmd, err
	}

	return cmd, nil
}

//uploadBinary exec scp command copy local file to the node.
func uploadBinary(file, dest, host, private, user string) (string, error) {
	cmd := fmt.Sprintf("scp -i %s %s %s", private, file, dest)

	if err := utils.SSHScp(user, private, host, DefaultSSHPort, file, dest); err != nil {
		return cmd, err
	}

	return cmd, nil
}
