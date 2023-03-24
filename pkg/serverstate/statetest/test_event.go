package statetest

//
//import (
//	"context"
//	"fmt"
//	"math/rand"
//	"testing"
//	"time"
//
//	pb "github.com/hashicorp/waypoint/pkg/server/gen"
//	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"
//	"github.com/stretchr/testify/require"
//)
//
//func init() {
//	tests["events"] = []testFunc{
//		TestEvent,
//		TestEventListPagination,
//	}
//}
//
//func TestEvent(t *testing.T, factory Factory, restartF RestartFactory) {
//	//TODO: Update this for events, not projects
//	ctx := context.Background()
//
//	t.Run("Basic put all types", func(t *testing.T) {
//		require := require.New(t)
//
//		s := factory(t)
//		defer s.Close()
//
//		ws := &pb.Ref_Workspace{
//			Workspace: "default",
//		}
//
//		// Write project
//		ref := &pb.Ref_Project{Project: "foo"}
//		require.NoError(s.ProjectPut(ctx, serverptypes.TestProject(t, &pb.Project{
//			Name: ref.Project,
//		})))
//
//		// Put Build
//		err := s.EventPut(ctx, &pb.UI_EventBuild{
//			BuildId:   "a1",
//			Sequence:  0,
//			Component: &pb.Component{Name: "test"},
//			Workspace: ws,
//		})
//		require.NoError(err)
//
//		// Put Deployment
//		err = s.EventPut(ctx, &pb.UI_EventDeployment{
//			DeploymentId: "a1",
//			Sequence:     0,
//			Component:    &pb.Component{Name: "test"},
//			Workspace:    ws,
//		})
//		require.NoError(err)
//
//		// Put Release
//		err = s.EventPut(ctx, &pb.UI_EventRelease{
//			ReleaseId: "a1",
//			Sequence:  0,
//			Component: &pb.Component{Name: "test"},
//			Workspace: ws,
//		})
//		require.NoError(err)
//
//		// Put Pipeline_Run
//		err = s.EventPut(ctx, &pb.UI_EventPipelineRun{
//			PipelineRunId: "a1",
//			Sequence:      0,
//			Component:     &pb.Component{Name: "test"},
//			Workspace:     ws,
//		})
//		require.NoError(err)
//
//		//TODO:put another one
//
//		// Create another one so we're sure that List can see more than one.
//		// Put Build
//		err = s.EventPut(ctx, &pb.UI_EventBuild{
//			BuildId:   "a2",
//			Sequence:  0,
//			Component: &pb.Component{Name: "test"},
//			Workspace: ws,
//		})
//		require.NoError(err)
//
//		// Put Deployment
//		err = s.EventPut(ctx, &pb.UI_EventDeployment{
//			DeploymentId: "a2",
//			Sequence:     0,
//			Component:    &pb.Component{Name: "test"},
//			Workspace:    ws,
//		})
//		require.NoError(err)
//
//		// Put Release
//		err = s.EventPut(ctx, &pb.UI_EventRelease{
//			ReleaseId: "a2",
//			Sequence:  0,
//			Component: &pb.Component{Name: "test"},
//			Workspace: ws,
//		})
//		require.NoError(err)
//
//		// Put Pipeline_Run
//		err = s.EventPut(ctx, &pb.UI_EventPipelineRun{
//			PipelineRunId: "a2",
//			Sequence:      0,
//			Component:     nil,
//			Workspace:     ws,
//		})
//		require.NoError(err)
//
//		//
//		//resp, _, err := s.EventListBundles(ctx, &pb.UI_ListEventsRequest{
//		//	Application: &pb.Ref_Application{
//		//		Application: "app",
//		//		Project:     "project",
//		//	},
//		//	Workspace:   &pb.Ref_Workspace{Workspace: "default"},
//		//	Pagination:  &pb.PaginationRequest{
//		//		PageSize:          5,
//		//		NextPageToken:     "",
//		//		PreviousPageToken: "",
//		//	},
//		//	//Sorting:     &pb.SortingRequest{OrderBy: []string{"name","created_at asc"}},
//		//})
//		//require.NoError(err)
//		//require.Len(resp, int(eventCount))
//		//
//		//// ListBundles
//		//{
//		//	resp, _, err := s.ProjectListBundles(ctx, &pb.PaginationRequest{})
//		//	require.NoError(err)
//		//	require.Len(resp, 2)
//		//}
//	})
//}
//
//func TestEventListPagination(t *testing.T, factory Factory, rf RestartFactory) {
//	ctx := context.Background()
//	require := require.New(t)
//	s := factory(t)
//	defer s.Close()
//	// a b c d e
//	// f g h i j
//	// k l m n o
//	// p q r s t
//	// u v w x y
//	// z
//	startChar := 'a'
//	endChar := 'm'
//	eventCount := endChar - startChar + 1
//	var chars []string
//
//	// Generate randomized events
//	for char := startChar; char <= endChar; char++ {
//		chars = append(chars, fmt.Sprintf("%c", char))
//	}
//	rand.Seed(time.Now().UnixNano())
//	rand.Shuffle(len(chars), func(i, j int) {
//		chars[i], chars[j] = chars[j], chars[i]
//	})
//	for _, char := range chars {
//		err := s.BuildPut(ctx, false, serverptypes.TestBuild(t, &pb.Build{
//			Id:       char,
//			Sequence: 1,
//			Application: &pb.Ref_Application{
//				Application: "app",
//				Project:     "project",
//			},
//			Workspace: &pb.Ref_Workspace{Workspace: "default"},
//		}))
//		require.NoError(err)
//
//		err = s.DeploymentPut(ctx, false, serverptypes.TestDeployment(t, &pb.Deployment{
//			Id: char,
//			Application: &pb.Ref_Application{
//				Application: "app",
//				Project:     "project",
//			},
//			Workspace: &pb.Ref_Workspace{Workspace: "default"},
//		}))
//		require.NoError(err)
//
//		err = s.ReleasePut(ctx, false, serverptypes.TestRelease(t, &pb.Release{
//			Id: char,
//			Application: &pb.Ref_Application{
//				Application: "app",
//				Project:     "project",
//			},
//			Workspace: &pb.Ref_Workspace{Workspace: "default"},
//		}))
//		require.NoError(err)
//
//		err = s.PipelinePut(ctx, &pb.Pipeline{
//			Id:   "testPipeline",
//			Name: "testPipeline",
//			Owner: &pb.Pipeline_Project{
//				Project: &pb.Ref_Project{
//					Project: "project",
//				},
//			},
//			Steps: map[string]*pb.Pipeline_Step{
//				"testStep": {
//					Name: "testStep",
//					Kind: &pb.Pipeline_Step_Up_{
//						Up: &pb.Pipeline_Step_Up{},
//					},
//				},
//			},
//		})
//		require.NoError(err)
//
//		err = s.PipelineRunPut(ctx, &pb.PipelineRun{
//			Id: char,
//			Pipeline: &pb.Ref_Pipeline{
//				Ref: &pb.Ref_Pipeline_Id{
//					Id: "test",
//				},
//			},
//			State: pb.PipelineRun_PENDING,
//		})
//		require.NoError(err)
//
//	}
//
//	t.Run("EventList", func(t *testing.T) {
//		t.Run(fmt.Sprintf("works with nil for compatibility and returns all %d results", eventCount), func(t *testing.T) {
//			{
//				resp, _, err := s.EventListBundles(ctx, &pb.UI_ListEventsRequest{
//					Application: &pb.Ref_Application{
//						Application: "app",
//						Project:     "project",
//					},
//					Workspace:   &pb.Ref_Workspace{Workspace: "default"},
//					Pagination:  &pb.PaginationRequest{
//						PageSize:          5,
//						NextPageToken:     "",
//						PreviousPageToken: "",
//					},
//					//Sorting:     &pb.SortingRequest{OrderBy: []string{"name","created_at asc"}},
//				})
//				require.NoError(err)
//				require.Len(resp, int(eventCount))
//
//			}
//		})
//
//	})
//
//}
//
////TODO: JUST test that the length returned from the eventlistbundle is correct since pagination is already
////thoroughly tested
