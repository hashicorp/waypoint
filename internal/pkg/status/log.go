package status

import "github.com/hashicorp/go-hclog"

type HCLog struct {
	L     hclog.Logger
	Level hclog.Level
}

func (h *HCLog) Update(str string) {
	if h.L == nil {
		h.L = hclog.L()
	}

	if h.Level == hclog.NoLevel {
		h.Level = hclog.Info
	}

	h.L.Log(h.Level, str)
}

func (h *HCLog) Close() {}

var _ Updater = &HCLog{}
