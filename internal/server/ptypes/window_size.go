package ptypes

import (
	"github.com/creack/pty"

	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

// Winsize turns a proto WindowSize to a pty.Winsize
func Winsize(ws *pb.ExecStreamRequest_WindowSize) *pty.Winsize {
	if ws == nil {
		return nil
	}

	return &pty.Winsize{
		Rows: uint16(ws.Rows),
		Cols: uint16(ws.Cols),
		X:    uint16(ws.Width),
		Y:    uint16(ws.Height),
	}
}

func WinsizeProto(ws *pty.Winsize) *pb.ExecStreamRequest_WindowSize {
	if ws == nil {
		return nil
	}

	return &pb.ExecStreamRequest_WindowSize{
		Rows:   int32(ws.Rows),
		Cols:   int32(ws.Cols),
		Width:  int32(ws.X),
		Height: int32(ws.Y),
	}
}
