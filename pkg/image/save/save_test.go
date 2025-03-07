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

package save

import (
	"context"
	"testing"

	v1 "github.com/sealerio/sealer/types/api/v1"
)

func TestSaveImages(t *testing.T) {
	tests := []string{"ubuntu", "ubuntu:18.04", "registry.aliyuncs.com/google_containers/coredns:1.6.5", "fanux/lvscare", "kubernetesui/dashboard:v2.2.0", "multiarch/ubuntu-core:arm64-focal"}
	is := NewImageSaver(context.Background())
	err := is.SaveImages(tests, "/var/lib/registry", v1.Platform{OS: "linux", Architecture: "amd64"})
	if err != nil {
		t.Error(err)
	}
}

func Test_splitDockerDomain(t *testing.T) {
	tests := []struct {
		name       string
		imageName  string
		wantDomain string
		wantRemain string
	}{
		{
			name:       "test1",
			imageName:  "docker.io/library/alpine:latest",
			wantDomain: defaultDomain,
			wantRemain: "library/alpine:latest",
		},
		{
			name:       "test2",
			imageName:  "ubuntu",
			wantDomain: defaultDomain,
			wantRemain: "library/ubuntu",
		},
		{
			name:       "test3",
			imageName:  "k8s.gcr.io/kube-apiserver",
			wantDomain: "k8s.gcr.io",
			wantRemain: "kube-apiserver",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if domain, remainer := splitDockerDomain(tt.imageName, ""); domain != tt.wantDomain || remainer != tt.wantRemain {
				t.Errorf("split image %s error", tt.name)
			}
		})
	}
}

func Test_parseNormalizedNamed(t *testing.T) {
	tests := []struct {
		name       string
		imageName  string
		wantDomain string
		wantRepo   string
		wantTag    string
	}{
		{
			name:       "test1",
			imageName:  "docker.io/library/alpine:latest",
			wantDomain: defaultDomain,
			wantRepo:   "library/alpine",
			wantTag:    defaultTag,
		},
		{
			name:       "test2",
			imageName:  "ubuntu",
			wantDomain: defaultDomain,
			wantRepo:   "library/ubuntu",
			wantTag:    defaultTag,
		},
		{
			name:       "test3",
			imageName:  "k8s.gcr.io/kube-apiserver",
			wantDomain: "k8s.gcr.io",
			wantRepo:   "kube-apiserver",
			wantTag:    defaultTag,
		},
		{
			name:       "test4",
			imageName:  "fanux/lvscare",
			wantDomain: defaultDomain,
			wantRepo:   "fanux/lvscare",
			wantTag:    defaultTag,
		},
		{
			name:       "test5",
			imageName:  "alpine",
			wantDomain: defaultDomain,
			wantRepo:   "library/alpine",
			wantTag:    defaultTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if named, err := ParseNormalizedNamed(tt.imageName, ""); err != nil || named.Domain() != tt.wantDomain || named.Repo() != tt.wantRepo || named.tag != tt.wantTag {
				t.Errorf("parse image %s error", tt.name)
			}
		})
	}
}
