/*
 * Copyright 2019-present Ciena Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package commands

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
	yaml "gopkg.in/yaml.v2"
)

const copyrightNotice = `# Copyright 2021-present Ciena Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.`

type ConfigOptions struct {
}

func RegisterConfigCommands(parent *flags.Parser) {
	if command, err := parent.AddCommand("config", "generate bossctl configuration", "Commands to generate bossctl configuration", &ConfigOptions{}); err != nil {
		Error.Fatalf("Unexpected error while attempting to register config commands : %s", err)
	} else {
		command.SubcommandsOptional = true
	}
}

func (options *ConfigOptions) Execute(args []string) error {
	ReadConfig()
	ApplyOptionOverrides(nil)
	b, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		return err
	}
	fmt.Println(copyrightNotice)
	fmt.Println(string(b))
	return nil
}
