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

package objects

import (
	"encoding/json"

	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

// Infra is
type Infra struct {
	Name       string                 `json:"name" yaml:"name"`
	Version    string                 `json:"version" yaml:"version" `
	Master     int                    `json:"master" yaml:"master"`
	Minion     int                    `json:"minion" yaml:"minion"`
	Result     bool                   `json:"result,omitempty" yaml:"result,omitempty"`
	Logs       []string               `json:"logs,omitempty" yaml:"logs,omitempty"`
	Components []Component            `json:"components" yaml:"components"`
	Outputs    map[string]interface{} `json:"-" yaml:"-"`
}

//WriteLog implement Logger interface.
func (i *Infra) WriteLog(log string, writer io.Writer, output bool) error {
	i.Logs = append(i.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
	}

	return nil
}

// JSON export deployment data
func (i *Infra) JSON() ([]byte, error) {
	return json.Marshal(&i)
}

//
func (i *Infra) YAML() ([]byte, error) {
	return yaml.Marshal(&i)
}

func (i *Infra) Output(key, value string) {
	if i.Outputs == nil {
		i.Outputs = map[string]interface{}{}
	}

	i.Outputs[key] = value
}
