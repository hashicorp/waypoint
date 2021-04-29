package state

import (
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// MetadataSet updates metadata on a project or application.
func (s *State) MetadataSet(req *pb.MetadataSetRequest) error {
	switch scope := req.Scope.(type) {
	case *pb.MetadataSetRequest_Application:
		app, err := s.AppGet(scope.Application)
		if err != nil {
			return err
		}

		switch val := req.Value.(type) {
		case *pb.MetadataSetRequest_FileChangeSignal:
			app.FileChangeSignal = val.FileChangeSignal
		}

		_, err = s.AppPut(app)
		return err
	case *pb.MetadataSetRequest_Project:
		project, err := s.ProjectGet(scope.Project)
		if err != nil {
			return err
		}

		switch val := req.Value.(type) {
		case *pb.MetadataSetRequest_FileChangeSignal:
			project.FileChangeSignal = val.FileChangeSignal
		}

		return s.ProjectPut(project)
	}

	return nil
}

// MetadataGetFileChangeSignal checks the metadata for the given application and it's
// project, returning the value of FileChangeSignal that is most relevent.
func (s *State) MetadataGetFileChangeSignal(scope *pb.Ref_Application) (string, error) {
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
