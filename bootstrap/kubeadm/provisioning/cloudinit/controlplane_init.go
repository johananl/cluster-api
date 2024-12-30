/*
Copyright 2019 The Kubernetes Authors.

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

package cloudinit

import (
	"sigs.k8s.io/cluster-api/util/secret"
)

const (
	controlPlaneCloudInit = `{{.Header}}
{{template "files" .WriteFiles}}
-   path: /run/kubeadm/kubeadm.yaml
    owner: root:root
    permissions: '0640'
    content: |
      ---
{{.ClusterConfiguration | Indent 6}}
      ---
{{.InitConfiguration | Indent 6}}
-   path: /run/cluster-api/placeholder
    owner: root:root
    permissions: '0640'
    content: "This placeholder file is used to create the /run/cluster-api sub directory in a way that is compatible with both Linux and Windows (mkdir -p /run/cluster-api does not work with Windows)"
{{- template "boot_commands" .BootCommands }}
runcmd:
{{- template "commands" .PreKubeadmCommands }}
  - 'kubeadm init --config /run/kubeadm/kubeadm.yaml {{.KubeadmVerbosity}} && {{ .SentinelFileCommand }}'
{{- template "commands" .PostKubeadmCommands }}
{{- template "ntp" .NTP }}
{{- template "users" .Users }}
{{- template "disk_setup" .DiskSetup}}
{{- template "fs_setup" .DiskSetup}}
{{- template "mounts" .Mounts}}
`
)

// ControlPlaneInitInput defines the context for generating provisioning data for initializing a
// control plane node.
type ControlPlaneInitInput struct {
	BaseUserData
	secret.Certificates

	ClusterConfiguration string
	InitConfiguration    string
}

func controlPlaneInitData(input *ControlPlaneInitInput) ([]byte, error) {
	input.Header = cloudConfigHeader
	input.WriteFiles = input.AsFiles()
	input.WriteFiles = append(input.WriteFiles, input.AdditionalFiles...)
	input.SentinelFileCommand = sentinelFileCommand
	userData, err := generate("InitControlplane", controlPlaneCloudInit, input)
	if err != nil {
		return nil, err
	}

	return userData, nil
}
