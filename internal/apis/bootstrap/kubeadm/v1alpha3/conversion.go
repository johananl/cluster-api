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

package v1alpha3

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	v1beta1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/upstreamv1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func (src *KubeadmConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1.KubeadmConfig)

	if err := Convert_v1alpha3_KubeadmConfig_To_v1beta1_KubeadmConfig(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &bootstrapv1.KubeadmConfig{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	MergeRestoredKubeadmConfigSpec(&dst.Spec, &restored.Spec)
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

func MergeRestoredKubeadmConfigSpec(dst *bootstrapv1.KubeadmConfigSpec, restored *bootstrapv1.KubeadmConfigSpec) {
	dst.Files = restored.Files

	dst.Users = restored.Users
	if restored.Users != nil {
		for i := range restored.Users {
			if restored.Users[i].PasswdFrom != nil {
				dst.Users[i].PasswdFrom = restored.Users[i].PasswdFrom
			}
		}
	}

	dst.BootCommands = restored.BootCommands
	dst.Ignition = restored.Ignition

	if restored.ClusterConfiguration != nil {
		if dst.ClusterConfiguration == nil {
			dst.ClusterConfiguration = &bootstrapv1.ClusterConfiguration{}
		}
		dst.ClusterConfiguration.APIServer.ExtraEnvs = restored.ClusterConfiguration.APIServer.ExtraEnvs
		dst.ClusterConfiguration.ControllerManager.ExtraEnvs = restored.ClusterConfiguration.ControllerManager.ExtraEnvs
		dst.ClusterConfiguration.Scheduler.ExtraEnvs = restored.ClusterConfiguration.Scheduler.ExtraEnvs

		if restored.ClusterConfiguration.Etcd.Local != nil {
			if dst.ClusterConfiguration.Etcd.Local == nil {
				dst.ClusterConfiguration.Etcd.Local = &bootstrapv1.LocalEtcd{}
			}
			dst.ClusterConfiguration.Etcd.Local.ExtraEnvs = restored.ClusterConfiguration.Etcd.Local.ExtraEnvs
		}
	}

	if restored.InitConfiguration != nil {
		if dst.InitConfiguration == nil {
			dst.InitConfiguration = &bootstrapv1.InitConfiguration{}
		}
		dst.InitConfiguration.Patches = restored.InitConfiguration.Patches
		dst.InitConfiguration.SkipPhases = restored.InitConfiguration.SkipPhases

		// Important! whenever adding fields to NodeRegistration, same fields must be added to hub.NodeRegistration's custom serialization func
		// otherwise those field won't exist in restored.

		dst.InitConfiguration.NodeRegistration.IgnorePreflightErrors = restored.InitConfiguration.NodeRegistration.IgnorePreflightErrors
		dst.InitConfiguration.NodeRegistration.ImagePullPolicy = restored.InitConfiguration.NodeRegistration.ImagePullPolicy
		dst.InitConfiguration.NodeRegistration.ImagePullSerial = restored.InitConfiguration.NodeRegistration.ImagePullSerial
	}

	if restored.JoinConfiguration != nil {
		if dst.JoinConfiguration == nil {
			dst.JoinConfiguration = &bootstrapv1.JoinConfiguration{}
		}
		dst.JoinConfiguration.Patches = restored.JoinConfiguration.Patches
		dst.JoinConfiguration.SkipPhases = restored.JoinConfiguration.SkipPhases

		if restored.JoinConfiguration.Discovery.File != nil && restored.JoinConfiguration.Discovery.File.KubeConfig != nil {
			if dst.JoinConfiguration.Discovery.File == nil {
				dst.JoinConfiguration.Discovery.File = &bootstrapv1.FileDiscovery{}
			}
			dst.JoinConfiguration.Discovery.File.KubeConfig = restored.JoinConfiguration.Discovery.File.KubeConfig
		}

		// Important! whenever adding fields to NodeRegistration, same fields must be added to hub.NodeRegistration's custom serialization func
		// otherwise those field won't exist in restored.

		dst.JoinConfiguration.NodeRegistration.IgnorePreflightErrors = restored.JoinConfiguration.NodeRegistration.IgnorePreflightErrors
		dst.JoinConfiguration.NodeRegistration.ImagePullPolicy = restored.JoinConfiguration.NodeRegistration.ImagePullPolicy
		dst.JoinConfiguration.NodeRegistration.ImagePullSerial = restored.JoinConfiguration.NodeRegistration.ImagePullSerial
	}
}

func (dst *KubeadmConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1.KubeadmConfig)

	if err := Convert_v1beta1_KubeadmConfig_To_v1alpha3_KubeadmConfig(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *KubeadmConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1.KubeadmConfigList)

	return Convert_v1alpha3_KubeadmConfigList_To_v1beta1_KubeadmConfigList(src, dst, nil)
}

func (dst *KubeadmConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1.KubeadmConfigList)

	return Convert_v1beta1_KubeadmConfigList_To_v1alpha3_KubeadmConfigList(src, dst, nil)
}

func (src *KubeadmConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1.KubeadmConfigTemplate)

	if err := Convert_v1alpha3_KubeadmConfigTemplate_To_v1beta1_KubeadmConfigTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &bootstrapv1.KubeadmConfigTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta

	MergeRestoredKubeadmConfigSpec(&dst.Spec.Template.Spec, &restored.Spec.Template.Spec)

	return nil
}

func (dst *KubeadmConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1.KubeadmConfigTemplate)

	if err := Convert_v1beta1_KubeadmConfigTemplate_To_v1alpha3_KubeadmConfigTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *KubeadmConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*bootstrapv1.KubeadmConfigTemplateList)

	return Convert_v1alpha3_KubeadmConfigTemplateList_To_v1beta1_KubeadmConfigTemplateList(src, dst, nil)
}

func (dst *KubeadmConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*bootstrapv1.KubeadmConfigTemplateList)

	return Convert_v1beta1_KubeadmConfigTemplateList_To_v1alpha3_KubeadmConfigTemplateList(src, dst, nil)
}

func Convert_v1alpha3_KubeadmConfigStatus_To_v1beta1_KubeadmConfigStatus(in *KubeadmConfigStatus, out *bootstrapv1.KubeadmConfigStatus, s apiconversion.Scope) error {
	// KubeadmConfigStatus.BootstrapData has been removed in v1alpha4 because its content has been moved to the bootstrap data secret, value will be lost during conversion.
	return autoConvert_v1alpha3_KubeadmConfigStatus_To_v1beta1_KubeadmConfigStatus(in, out, s)
}

func Convert_v1beta1_ClusterConfiguration_To_upstreamv1beta1_ClusterConfiguration(in *bootstrapv1.ClusterConfiguration, out *upstreamv1beta1.ClusterConfiguration, s apiconversion.Scope) error {
	// DNS.Type was removed in v1alpha4 because only CoreDNS is supported; the information will be left to empty (kubeadm defaults it to CoredDNS);
	// Existing clusters using kube-dns or other DNS solutions will continue to be managed/supported via the skip-coredns annotation.

	// ClusterConfiguration.UseHyperKubeImage was removed in kubeadm v1alpha4 API
	return upstreamv1beta1.Convert_v1beta1_ClusterConfiguration_To_upstreamv1beta1_ClusterConfiguration(in, out, s)
}

func Convert_upstreamv1beta1_ClusterConfiguration_To_v1beta1_ClusterConfiguration(in *upstreamv1beta1.ClusterConfiguration, out *bootstrapv1.ClusterConfiguration, s apiconversion.Scope) error {
	// DNS.Type was removed in v1alpha4 because only CoreDNS is supported; the information will be left to empty (kubeadm defaults it to CoredDNS);
	// ClusterConfiguration.UseHyperKubeImage was removed in kubeadm v1alpha4 API
	return upstreamv1beta1.Convert_upstreamv1beta1_ClusterConfiguration_To_v1beta1_ClusterConfiguration(in, out, s)
}

func Convert_upstreamv1beta1_InitConfiguration_To_v1beta1_InitConfiguration(in *upstreamv1beta1.InitConfiguration, out *bootstrapv1.InitConfiguration, s apiconversion.Scope) error {
	// NodeRegistrationOptions.IgnorePreflightErrors does not exist in kubeadm v1beta1 API
	return upstreamv1beta1.Convert_upstreamv1beta1_InitConfiguration_To_v1beta1_InitConfiguration(in, out, s)
}

func Convert_v1beta1_InitConfiguration_To_upstreamv1beta1_InitConfiguration(in *bootstrapv1.InitConfiguration, out *upstreamv1beta1.InitConfiguration, s apiconversion.Scope) error {
	// NodeRegistrationOptions.IgnorePreflightErrors does not exist in kubeadm v1beta1 API
	return upstreamv1beta1.Convert_v1beta1_InitConfiguration_To_upstreamv1beta1_InitConfiguration(in, out, s)
}

func Convert_upstreamv1beta1_JoinConfiguration_To_v1beta1_JoinConfiguration(in *upstreamv1beta1.JoinConfiguration, out *bootstrapv1.JoinConfiguration, s apiconversion.Scope) error {
	// NodeRegistrationOptions.IgnorePreflightErrors does not exist in kubeadm v1beta1 API
	return upstreamv1beta1.Convert_upstreamv1beta1_JoinConfiguration_To_v1beta1_JoinConfiguration(in, out, s)
}

func Convert_v1beta1_JoinConfiguration_To_upstreamv1beta1_JoinConfiguration(in *bootstrapv1.JoinConfiguration, out *upstreamv1beta1.JoinConfiguration, s apiconversion.Scope) error {
	// NodeRegistrationOptions.IgnorePreflightErrors does not exist in kubeadm v1beta1 API
	return upstreamv1beta1.Convert_v1beta1_JoinConfiguration_To_upstreamv1beta1_JoinConfiguration(in, out, s)
}

// Convert_v1beta1_KubeadmConfigSpec_To_v1alpha3_KubeadmConfigSpec is an autogenerated conversion function.
func Convert_v1beta1_KubeadmConfigSpec_To_v1alpha3_KubeadmConfigSpec(in *bootstrapv1.KubeadmConfigSpec, out *KubeadmConfigSpec, s apiconversion.Scope) error {
	// KubeadmConfigSpec.Ignition does not exist in kubeadm v1alpha3 API.

	if in.ClusterConfiguration != nil {
		out.ClusterConfiguration = &upstreamv1beta1.ClusterConfiguration{}
		if err := Convert_v1beta1_ClusterConfiguration_To_upstreamv1beta1_ClusterConfiguration(in.ClusterConfiguration, out.ClusterConfiguration, s); err != nil {
			return err
		}
	} else {
		out.ClusterConfiguration = nil
	}

	if in.InitConfiguration != nil {
		out.InitConfiguration = &upstreamv1beta1.InitConfiguration{}
		if err := Convert_v1beta1_InitConfiguration_To_upstreamv1beta1_InitConfiguration(in.InitConfiguration, out.InitConfiguration, s); err != nil {
			return err
		}
	} else {
		out.InitConfiguration = nil
	}

	if in.JoinConfiguration != nil {
		out.JoinConfiguration = &upstreamv1beta1.JoinConfiguration{}
		if err := Convert_v1beta1_JoinConfiguration_To_upstreamv1beta1_JoinConfiguration(in.JoinConfiguration, out.JoinConfiguration, s); err != nil {
			return err
		}
	} else {
		out.JoinConfiguration = nil
	}

	out.PreKubeadmCommands = in.PreKubeadmCommands
	out.PostKubeadmCommands = in.PostKubeadmCommands
	out.Format = Format(in.Format)
	out.Verbosity = in.Verbosity

	return autoConvert_v1beta1_KubeadmConfigSpec_To_v1alpha3_KubeadmConfigSpec(in, out, s)
}

func Convert_v1beta1_File_To_v1alpha3_File(in *bootstrapv1.File, out *File, s apiconversion.Scope) error {
	// File.Append does not exist in kubeadm v1alpha3 API.
	return autoConvert_v1beta1_File_To_v1alpha3_File(in, out, s)
}

func Convert_v1beta1_User_To_v1alpha3_User(in *bootstrapv1.User, out *User, s apiconversion.Scope) error {
	// User.PasswdFrom does not exist in kubeadm v1alpha3 API.
	return autoConvert_v1beta1_User_To_v1alpha3_User(in, out, s)
}

func Convert_v1beta1_KubeadmConfigTemplateResource_To_v1alpha3_KubeadmConfigTemplateResource(in *bootstrapv1.KubeadmConfigTemplateResource, out *KubeadmConfigTemplateResource, s apiconversion.Scope) error {
	// KubeadmConfigTemplateResource.metadata does not exist in kubeadm v1alpha3.
	return autoConvert_v1beta1_KubeadmConfigTemplateResource_To_v1alpha3_KubeadmConfigTemplateResource(in, out, s)
}

func Convert_v1beta1_KubeadmConfigStatus_To_v1alpha3_KubeadmConfigStatus(in *bootstrapv1.KubeadmConfigStatus, out *KubeadmConfigStatus, s apiconversion.Scope) error {
	// V1Beta2 was added in v1beta1.
	return autoConvert_v1beta1_KubeadmConfigStatus_To_v1alpha3_KubeadmConfigStatus(in, out, s)
}

func Convert_v1alpha3_KubeadmConfigSpec_To_v1beta1_KubeadmConfigSpec(in *KubeadmConfigSpec, out *bootstrapv1.KubeadmConfigSpec, s apiconversion.Scope) error {
	if in.ClusterConfiguration != nil {
		out.ClusterConfiguration = &v1beta1.ClusterConfiguration{}
		if err := Convert_upstreamv1beta1_ClusterConfiguration_To_v1beta1_ClusterConfiguration(in.ClusterConfiguration, out.ClusterConfiguration, s); err != nil {
			return err
		}
	} else {
		out.ClusterConfiguration = nil
	}

	if in.InitConfiguration != nil {
		out.InitConfiguration = &v1beta1.InitConfiguration{}
		if err := Convert_upstreamv1beta1_InitConfiguration_To_v1beta1_InitConfiguration(in.InitConfiguration, out.InitConfiguration, s); err != nil {
			return err
		}
	} else {
		out.InitConfiguration = nil
	}

	if in.JoinConfiguration != nil {
		out.JoinConfiguration = &v1beta1.JoinConfiguration{}
		if err := Convert_upstreamv1beta1_JoinConfiguration_To_v1beta1_JoinConfiguration(in.JoinConfiguration, out.JoinConfiguration, s); err != nil {
			return err
		}
	} else {
		out.JoinConfiguration = nil
	}

	out.PreKubeadmCommands = in.PreKubeadmCommands
	out.PostKubeadmCommands = in.PostKubeadmCommands
	out.Format = v1beta1.Format(in.Format)
	out.Verbosity = in.Verbosity

	return autoConvert_v1alpha3_KubeadmConfigSpec_To_v1beta1_KubeadmConfigSpec(in, out, s)
}
