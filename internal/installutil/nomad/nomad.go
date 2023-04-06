// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"context"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"time"
)

const (
	// bytes
	DefaultCSIVolumeCapacityMin = int64(1073741824)
	DefaultCSIVolumeCapacityMax = int64(2147483648)
	DefaultCSIVolumeMountFS     = "xfs"
	DefaultResourcesCPU         = 200
	DefaultResourcesMemory      = 600
)

type NomadConfig struct {
	authSoftFail       bool              `hcl:"auth_soft_fail,optional"`
	image              string            `hcl:"server_image,optional"`
	namespace          string            `hcl:"namespace,optional"`
	serviceAnnotations map[string]string `hcl:"service_annotations,optional"`

	consulService            bool     `hcl:"consul_service,optional"`
	consulServiceUITags      []string `hcl:"consul_service_ui_tags:optional"`
	consulServiceBackendTags []string `hcl:"consul_service_backend_tags:optional"`
	consulDatacenter         string   `hcl:"consul_datacenter,optional"`
	consulDomain             string   `hcl:"consul_datacenter,optional"`

	// If set along with consul, will use this hostname instead of
	// making a consul DNS hostname for the server address in its context
	consulServiceHostname string `hcl:"consul_service_hostname,optional"`

	odrImage string `hcl:"odr_image,optional"`

	region         string   `hcl:"namespace,optional"`
	datacenters    []string `hcl:"datacenters,optional"`
	policyOverride bool     `hcl:"policy_override,optional"`

	serverResourcesCPU    string `hcl:"server_resources_cpu,optional"`
	serverResourcesMemory string `hcl:"server_resources_memory,optional"`
	runnerResourcesCPU    string `hcl:"runner_resources_cpu,optional"`
	runnerResourcesMemory string `hcl:"runner_resources_memory,optional"`

	hostVolume           string `hcl:"host_volume,optional"`
	csiVolumeProvider    string `hcl:"csi_volume_provider,optional"`
	csiVolumeCapacityMin int64  `hcl:"csi_volume_capacity_min,optional"`
	csiVolumeCapacityMax int64  `hcl:"csi_volume_capacity_max,optional"`
	csiFS                string `hcl:"csi_fs,optional"`

	nomadHost     string            `hcl:"nomad_host,optional"`
	csiTopologies map[string]string `hcl:"nomad_csi_topologies,optional"`
	csiExternalId string            `hcl:"nomad_csi_external_id,optional"`
	csiPluginId   string            `hcl:"nomad_csi_plugin_id,optional"`
	csiSecrets    map[string]string `hcl:"nomad_csi_secrets,optional"`
}

func RunJob(
	ctx context.Context,
	s terminal.Step,
	client *api.Client,
	job *api.Job,
	policyOverride bool,
) (string, error) {
	jobOpts := &api.RegisterOptions{
		PolicyOverride: policyOverride,
	}

	resp, _, err := client.Jobs().RegisterOpts(job, jobOpts, nil)
	if err != nil {
		return "", err
	}

	s.Update("Waiting for allocation to be scheduled")
	qopts := &api.QueryOptions{
		WaitIndex: resp.EvalCreateIndex,
		WaitTime:  time.Duration(500 * time.Millisecond),
	}

	eval, meta, err := waitForEvaluation(ctx, s, client, resp, qopts)
	if err != nil {
		return "", err
	}
	if eval == nil {
		return "", fmt.Errorf("evaluation status could not be determined")
	}
	qopts.WaitIndex = meta.LastIndex

	var allocID string
	retries := 0
	maxRetries := 3
	for {
		allocs, qmeta, err := client.Evaluations().Allocations(eval.ID, qopts)
		if err != nil {
			return "", err
		}
		qopts.WaitIndex = qmeta.LastIndex
		if len(allocs) == 0 {
			return "", fmt.Errorf("no allocations found after evaluation completed")
		}

		switch allocs[0].ClientStatus {
		case "running":
			allocID = allocs[0].ID
			s.Update("Nomad allocation running")
			retries++
		case "pending":
			s.Update(fmt.Sprintf("Waiting for allocation %q to start", allocs[0].ID))
			// retry
		default:
			return "", fmt.Errorf("allocation failed")
		}

		if allocID != "" {
			if retries == maxRetries {
				return allocID, nil
			} else {
				s.Update("Ensuring allocation %q has properly started up...", allocs[0].ID)
				time.Sleep(1 * time.Second)
			}
		}

		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func waitForEvaluation(
	ctx context.Context,
	s terminal.Step,
	client *api.Client,
	resp *api.JobRegisterResponse,
	qopts *api.QueryOptions,
) (*api.Evaluation, *api.QueryMeta, error) {

	for {
		eval, meta, err := client.Evaluations().Info(resp.EvalID, qopts)
		if err != nil {
			return nil, nil, err
		}

		qopts.WaitIndex = meta.LastIndex

		switch eval.Status {
		case "pending":
			s.Update("Nomad allocation pending...")
		case "complete":
			s.Update("Nomad allocation created")

			return eval, meta, nil
		case "failed", "canceled", "blocked":
			s.Update("Nomad failed to schedule the job")
			s.Status(terminal.StatusError)
			return nil, nil, fmt.Errorf("Nomad evaluation did not transition to 'complete'")
		default:
			return nil, nil, fmt.Errorf("receieved unknown eval status from Nomad: %q", eval.Status)
		}
	}
}

func CreatePersistentVolume(
	ctx context.Context,
	client *api.Client,
	id, name, csiPluginId, csiVolumeProvider, csiFS, csiExternalId string,
	csiVolumeCapacityMin, csiVolumeCapacityMax int64,
	csiTopologies, csiSecrets, csiParams map[string]string,
	csiMountFlags []string,
) error {
	vol := api.CSIVolume{
		ID:         id,
		Name:       name,
		ExternalID: csiExternalId,
		RequestedCapabilities: []*api.CSIVolumeCapability{
			{
				AccessMode:     "single-node-writer",
				AttachmentMode: "file-system",
			},
		},
		MountOptions: &api.CSIMountOptions{
			FSType:     DefaultCSIVolumeMountFS,
			MountFlags: csiMountFlags,
		},
		RequestedCapacityMin: DefaultCSIVolumeCapacityMin,
		RequestedCapacityMax: DefaultCSIVolumeCapacityMax,
		PluginID:             csiPluginId,
		Parameters:           csiParams,
		Provider:             csiVolumeProvider,
		RequestedTopologies: &api.CSITopologyRequest{
			Required:  nil,
			Preferred: nil,
		},
		Secrets: nil,
	}
	if csiVolumeCapacityMin != 0 {
		vol.RequestedCapacityMin = csiVolumeCapacityMin
	}
	if csiVolumeCapacityMax != 0 {
		vol.RequestedCapacityMax = csiVolumeCapacityMax
	}
	if csiFS != "" {
		vol.MountOptions.FSType = csiFS
	}
	if len(csiTopologies) != 0 {
		vol.RequestedTopologies.Required = []*api.CSITopology{
			{
				Segments: csiTopologies,
			},
		}
	}
	if len(csiSecrets) != 0 {
		vol.Secrets = csiSecrets
	}

	_, _, err := client.CSIVolumes().Create(&vol, &api.WriteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func SetupPretask(volumeMounts []*api.VolumeMount) *api.Task {
	preTask := api.NewTask("pre_task", "docker")
	// Observed WP user and group IDs in the published container, update if those ever change
	waypointUserID := 100
	waypointGroupID := 1000
	cpu := DefaultResourcesCPU
	mem := DefaultResourcesMemory
	preTask.Config = map[string]interface{}{
		// Doing this because this is the only way https://github.com/hashicorp/nomad/issues/8892
		"image":   "busybox:latest",
		"command": "sh",
		"args":    []string{"-c", fmt.Sprintf("chown -R %d:%d /data/", waypointUserID, waypointGroupID)},
	}
	preTask.VolumeMounts = volumeMounts
	preTask.Resources = &api.Resources{
		CPU:      &cpu,
		MemoryMB: &mem,
	}
	preTask.Lifecycle = &api.TaskLifecycle{
		Hook:    "prestart",
		Sidecar: false,
	}
	return preTask
}
