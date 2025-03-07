// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package collector

import "fmt"

type Collector interface {
	// Collect git package;download common file work as wget or curl;copy local file to dst.
	Collect(buildContext, src, savePath string) error
}

func NewCollector(src string) (Collector, error) {
	// if src is detected as remote context,will new different Collector via src type.
	switch {
	case src == "":
		return nil, fmt.Errorf("src can not be nil")
	case IsGitURL(src):
		// remote git context
		return NewGitCollector(), nil
	case IsURL(src):
		// remote web context
		return NewWebFileCollector(), nil
	default:
		//local context
		return NewLocalCollector(), nil
	}
}
