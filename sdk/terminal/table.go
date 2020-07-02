package terminal

import (
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var colorMapping = map[string]int{
	Green:  tablewriter.FgGreenColor,
	Yellow: tablewriter.FgYellowColor,
	Red:    tablewriter.FgRedColor,
}

func (u *BasicUI) Table(tbl *Table, opts ...Option) {
	// Build our config and set our options
	cfg := &config{Writer: color.Output}
	for _, opt := range opts {
		opt(cfg)
	}

	table := tablewriter.NewWriter(cfg.Writer)
	table.SetHeader(tbl.Headers)
	table.SetBorder(false)

	for _, row := range tbl.Rows {
		colors := make([]tablewriter.Colors, len(row))
		entries := make([]string, len(row))

		for i, ent := range row {
			entries[i] = ent.Value

			color, ok := colorMapping[ent.Color]
			if ok {
				colors[i] = tablewriter.Colors{color}
			}
		}

		table.Rich(entries, colors)
	}

	table.Render()
}
