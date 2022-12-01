package container

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint-plugin-sdk/component"
	"github.com/hashicorp/waypoint-plugin-sdk/framework/resource"
	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/builtin/docker"
)

// Platform is the Platform implementation for Scaleway Container.
type Platform struct {
	config PlatformConfig
}

var (
	_ component.Configurable       = (*Platform)(nil)
	_ component.ConfigurableNotify = (*Platform)(nil)
	_ component.Platform           = (*Platform)(nil)
	_ component.Status             = (*Platform)(nil)
)

// PlatformConfig is the config for the Scaleway Container Platform
type PlatformConfig struct {
	Region    string `hcl:"region,optional"`
	Port      int    `hcl:"port,optional"`
	Namespace string `hcl:"namespace"`
}

func (p *Platform) ConfigSet(i interface{}) error {
	//TODO implement me
	return nil
}

func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

func (p *Platform) StatusFunc() interface{} {
	return p.status
}

func (p *Platform) DeployFunc() interface{} {
	return p.deploy
}

func (p *Platform) DestroyFunc() interface{} {
	return p.destroy
}

func (p *Platform) resourceManager(log hclog.Logger, dcr *component.DeclaredResourcesResp) *resource.Manager {
	return resource.NewManager(
		resource.WithLogger(log.Named("resource_manager")),
		resource.WithValueProvider(p.scalewayContainerAPI),
		resource.WithDeclaredResourcesResp(dcr),
		resource.WithResource(resource.NewResource(
			resource.WithName("container"),
			resource.WithPlatform("scaleway"),
			resource.WithCategoryDisplayHint(sdk.ResourceCategoryDisplayHint_INSTANCE_MANAGER),
			resource.WithState(&Resource_Container{}),
			resource.WithCreate(p.resourceContainerCreate),
			resource.WithStatus(p.resourceContainerStatus),
			resource.WithDestroy(p.resourceContainerDestroy),
		)),
	)
}

func (p *Platform) status(
	ctx context.Context,
	ji *component.JobInfo,
	ui terminal.UI,
	log hclog.Logger,
	container *Container,
) (*sdk.StatusReport, error) {
	sg := ui.StepGroup()
	s := sg.Add("Checking the status of the container deployment...")

	rm := p.resourceManager(log, nil)

	if container.ResourceState == nil {
		s.Update("Creating state")
		err := rm.Resource("container").SetState(&Resource_Container{
			Id: container.Id,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set resource container state: %w", err)
		}
	} else {
		s.Update("Loading state: %v", container.ResourceState)
		if err := rm.LoadState(container.ResourceState); err != nil {
			return nil, err
		}
	}

	report, err := rm.StatusReport(ctx, log, sg, ui)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "resource manager failed to generate resource statuses: %s", err)
	}

	//report.Health = sdk.StatusReport_READY
	//s.Update("Deployment no implemented: " + container.Image)
	s.Done()
	return report, nil
}

func (p *Platform) deploy(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	dcr *component.DeclaredResourcesResp,
	src *component.Source,
	img *docker.Image,
	deployConfig *component.DeploymentConfig,
) (*Container, error) {
	st := ui.Status()
	defer st.Close()
	st.Update("Deploying container")

	container := &Container{
		Region: p.config.Region,
	}

	id, err := component.Id()
	if err != nil {
		return nil, err
	}

	container.DeploymentId = id
	container.Name = strings.ToLower(fmt.Sprintf("%s-v%v", src.App, deployConfig.Sequence))

	rm := p.resourceManager(log, dcr)

	err = rm.CreateAll(ctx, container, log, st, deployConfig, img)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	st.Step(terminal.StatusOK, "Created resource")

	container.ResourceState = rm.State()

	servState := rm.Resource("container").State().(*Resource_Container)
	if servState == nil {
		return nil, status.Errorf(codes.Internal, "service state is nil")
	}

	return container, nil
}

func (p *Platform) destroy(
	ctx context.Context,
	ui terminal.UI,
	log hclog.Logger,
	container *Container,
) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	rm := p.resourceManager(log, nil)

	// If we don't have resource state, this state is from an older version
	// and we need to manually recreate it.
	if container.ResourceState == nil {
		rm.Resource("deployment").SetState(&Resource_Container{
			Id:     container.Id,
			Region: container.Region,
		})
	} else {
		// Load our set state
		if err := rm.LoadState(container.ResourceState); err != nil {
			return err
		}
	}

	// Destroy
	return rm.DestroyAll(ctx, log, sg, ui)
}
