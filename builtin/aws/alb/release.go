package alb

import "github.com/hashicorp/waypoint/sdk/component"

func (r *Release) URL() string { return r.Url }

var _ component.Release = (*Release)(nil)
