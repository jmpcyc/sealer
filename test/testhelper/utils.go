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

package testhelper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sealerio/sealer/utils/os/fs"

	"github.com/sealerio/sealer/utils/exec"
	"github.com/sealerio/sealer/utils/net"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"sigs.k8s.io/yaml"

	"github.com/sealerio/sealer/common"
	"github.com/sealerio/sealer/logger"
	"github.com/sealerio/sealer/test/testhelper/settings"
	v1 "github.com/sealerio/sealer/types/api/v1"
	"github.com/sealerio/sealer/utils/ssh"
)

func GetPwd() string {
	pwd, err := os.Getwd()
	CheckErr(err)
	return pwd
}

func CreateTempFile() string {
	dir := os.TempDir()
	file, err := ioutil.TempFile(dir, "tmpfile")
	CheckErr(err)
	defer CheckErr(file.Close())
	return file.Name()
}

func RemoveTempFile(file string) {
	CheckErr(os.Remove(file))
}

func WriteFile(fileName string, content []byte) error {
	dir := filepath.Dir(fileName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, settings.FileMode0755); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(fileName, content, settings.FileMode0644); err != nil {
		return err
	}
	return nil
}

type SSHClient struct {
	RemoteHostIP string
	SSH          ssh.Interface
}

func NewSSHByCluster(cluster *v1.Cluster) ssh.Interface {
	if cluster.Spec.SSH.User == "" {
		cluster.Spec.SSH.User = common.ROOT
	}
	address, err := net.GetLocalHostAddresses()
	if err != nil {
		logger.Warn("failed to get local address, %v", err)
	}
	return &ssh.SSH{
		Encrypted:    cluster.Spec.SSH.Encrypted,
		User:         cluster.Spec.SSH.User,
		Password:     cluster.Spec.SSH.Passwd,
		Port:         cluster.Spec.SSH.Port,
		PkFile:       cluster.Spec.SSH.Pk,
		PkPassword:   cluster.Spec.SSH.PkPasswd,
		LocalAddress: address,
		IsStdout:     true,
		Fs:           fs.NewFilesystem(),
	}
}

func NewSSHClientByCluster(cluster *v1.Cluster) *SSHClient {
	var (
		ipList []string
		host   string
	)
	sshClient := NewSSHByCluster(cluster)
	if cluster.Spec.Provider == common.AliCloud {
		host = cluster.GetAnnotationsByKey(common.Eip)
		CheckNotEqual(host, "")
		ipList = append(ipList, host)
	} else {
		host = cluster.Spec.Masters.IPList[0]
		ipList = append(ipList, append(cluster.Spec.Masters.IPList, cluster.Spec.Nodes.IPList...)...)
	}
	err := ssh.WaitSSHReady(sshClient, 6, ipList...)
	CheckErr(err)
	CheckNotNil(sshClient)

	return &SSHClient{
		SSH:          sshClient,
		RemoteHostIP: host,
	}
}

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func UnmarshalYamlFile(file string, obj interface{}) error {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, obj)
	return err
}

func MarshalYamlToFile(file string, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	if err = WriteFile(file, data); err != nil {
		return err
	}
	return nil
}

// GetFileDataLocally get file data for cloud apply
func GetFileDataLocally(filePath string) string {
	cmd := fmt.Sprintf("sudo -E cat %s", filePath)
	result, err := exec.RunSimpleCmd(cmd)
	CheckErr(err)
	return result
}

// DeleteFileLocally delete file for cloud apply
func DeleteFileLocally(filePath string) {
	cmd := fmt.Sprintf("sudo -E rm -rf %s", filePath)
	_, err := exec.RunSimpleCmd(cmd)
	CheckErr(err)
}

func CheckErr(err error) {
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

func CheckNotNil(obj interface{}) {
	gomega.Expect(obj).NotTo(gomega.BeNil())
}

func CheckEqual(obj1 interface{}, obj2 interface{}) {
	gomega.Expect(obj1).To(gomega.Equal(obj2))
}

func CheckNotEqual(obj1 interface{}, obj2 interface{}) {
	gomega.Expect(obj1).NotTo(gomega.Equal(obj2))
}

func CheckExit0(sess *gexec.Session, waitTime time.Duration) {
	gomega.Eventually(sess, waitTime).Should(gexec.Exit(0))
}
func CheckNotExit0(sess *gexec.Session, waitTime time.Duration) {
	gomega.Eventually(sess, waitTime).ShouldNot(gexec.Exit(0))
}

func CheckFuncBeTrue(f func() bool, t time.Duration) {
	gomega.Eventually(f(), t).Should(gomega.BeTrue())
}

func CheckBeTrue(b bool) {
	gomega.Eventually(b).Should(gomega.BeTrue())
}
func CheckNotBeTrue(b bool) {
	gomega.Eventually(b).ShouldNot(gomega.BeTrue())
}
