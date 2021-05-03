package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// GetFileChangeSignal checks the metadata for the given application and its
// project, returning the value of FileChangeSignal that is most relevent.
func (s *State) GetFileChangeSignal(scope *pb.Ref_Application) (string, error) {
	app, err := s.AppGet(scope)
	if err != nil {
		return "", err
	}

	if app.FileChangeSignal != "" {
		return app.FileChangeSignal, nil
	}

	project, err := s.ProjectGet(&pb.Ref_Project{Project: scope.Project})
	if err != nil {
		return "", err
	}

	return project.FileChangeSignal, nil
}
