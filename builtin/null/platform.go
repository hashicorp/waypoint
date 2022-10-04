package null

import (
	"time"

	sdk "github.com/hashicorp/waypoint-plugin-sdk/proto/gen"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Platform struct {
	config Config
}

type Config struct{}

// Config implements Configurable
func (p *Platform) Config() (interface{}, error) {
	return &p.config, nil
}

// DeployFunc implements component.Platform
func (p *Platform) DeployFunc() interface{} {
	return func(ui terminal.UI) *emptypb.Empty {
		sg := ui.StepGroup()
		step := sg.Add("performing null deploy")

		time.Sleep(time.Second * 2)

		step.Update("null deploy complete")
		step.Done()
		return &emptypb.Empty{}
	}
}

// DestroyFunc implements component.Destroyer
func (p *Platform) DestroyFunc() interface{} {
	return func(ui terminal.UI) error {
		sg := ui.StepGroup()
		step := sg.Add("performing null destroy")

		time.Sleep(time.Second * 1)

		step.Update("null destroy complete")
		step.Done()
		return nil
	}
}

// ValidateAuthFunc implements component.Authenticator
func (p *Platform) ValidateAuthFunc() interface{} {
	return func() error {
		return nil
	}
}

// AuthFunc implements component.Authenticator
func (p *Platform) AuthFunc() interface{} {
	return func() error {
		return nil
	}
}

func (p *Platform) StatusFunc() interface{} {
	return func(ui terminal.UI) *sdk.StatusReport {
		sg := ui.StepGroup()
		step := sg.Add("performing null status")

		time.Sleep(time.Second * 1)

		step.Update("null status complete")
		step.Done()
		return &sdk.StatusReport{}
	}
}

// DefaultReleaserFunc implements component.PlatformReleaser
func (p *Platform) DefaultReleaserFunc() interface{} {
	return func() *Releaser {
		return &Releaser{
			config: ReleaserConfig{},
		}
	}
}

// ConfigSet is called after a configuration has been decoded
// we can use this to validate the config
func (p *Platform) ConfigSet(interface{}) error {
	return nil
}
