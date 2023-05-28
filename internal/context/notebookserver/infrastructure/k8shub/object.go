package k8shub

import (
	corev1 "k8s.io/api/core/v1"
)

type kubespawnerOverride struct {
	Image                   string               `json:"image"`
	CPULimit                int64                `json:"cpu_limit"`
	CPUGuarantee            int64                `json:"cpu_guarantee"`
	MemLimit                int64                `json:"mem_limit"`
	MemGuarantee            int64                `json:"mem_guarantee"`
	ExtraResourceGuarantees map[string]string    `json:"extra_resource_guarantees"`
	ExtraResourceLimits     map[string]string    `json:"extra_resource_limits"`
	NodeSelector            map[string]string    `json:"node_selector"`
	Volumes                 []corev1.Volume      `json:"volumes"`
	VolumeMounts            []corev1.VolumeMount `json:"volume_mounts"`
	// StorageCapacity         string               `json:"storage_capacity"`
}

// define from item in c.KubeSpawner.profile_list
type profile struct {
	DisplayName         string              `json:"display_name"`
	Slug                string              `json:"slug"`
	KubespawnerOverride kubespawnerOverride `json:"kubespawner_override"`
}

type userOptions struct {
	Profile profile `json:"profile"`
}
