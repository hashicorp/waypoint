package handlers

import pb "github.com/hashicorp/waypoint/pkg/server/gen"

// varContainsDynamic returns true if there are any dynamic values in the list.
func varContainsDynamic(vars []*pb.ConfigVar) bool {
	for _, v := range vars {
		if _, ok := v.Value.(*pb.ConfigVar_Dynamic); ok {
			return true
		}
	}

	return false
}
