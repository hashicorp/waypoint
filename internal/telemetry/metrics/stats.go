package metrics

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func init() {
	// Register OpenCensus views.
	if err := view.Register(statsViews...); err != nil {
		log.Printf("========================")
		log.Printf("=== error registering view ===")
		log.Printf("=== %v ===", err)
		log.Printf("========================")
		fmt.Fprintf(os.Stderr, "error registering OpenCensus views: %v", err)
	} else {
		log.Printf("========================")
		log.Printf("=== no registering view ===")
		log.Printf("========================")
	}
}

var (
	// TagMethod is a tag for capturing the method.
	TagOperation = tag.MustNewKey("operation")

	operationDurationMeasure = stats.Float64(
		"waypoint_operation",
		"The number of seconds duration for this operation",
		stats.UnitSeconds,
	)

	operationCountMeasure = stats.Int64(
		"waypoint_operation_count",
		"count of operations",
		stats.UnitDimensionless,
	)

	waypointOperationCounts = &view.View{
		Name:        operationCountMeasure.Name(),
		Description: operationCountMeasure.Description(),
		TagKeys:     []tag.Key{TagOperation},
		Measure:     operationCountMeasure,
		Aggregation: view.Count(),
	}

	waypointOperationDurations = &view.View{
		Name:        operationDurationMeasure.Name(),
		Description: operationDurationMeasure.Description(),
		TagKeys:     []tag.Key{TagOperation},
		Measure:     operationDurationMeasure,
		// add a custom distribution bucket of 10 second intervals between
		// 0 and 240 seconds, then 10 minute intervals up to 60 minutes
		Aggregation: view.Distribution(
			(10 * time.Second).Seconds(),
			(20 * time.Second).Seconds(),
			(30 * time.Second).Seconds(),
			(40 * time.Second).Seconds(),
			(50 * time.Second).Seconds(),
			(60 * time.Second).Seconds(),
			(70 * time.Second).Seconds(),
			(80 * time.Second).Seconds(),
			(90 * time.Second).Seconds(),
			(100 * time.Second).Seconds(),
			(110 * time.Second).Seconds(),
			(120 * time.Second).Seconds(),
			(130 * time.Second).Seconds(),
			(140 * time.Second).Seconds(),
			(150 * time.Second).Seconds(),
			(160 * time.Second).Seconds(),
			(170 * time.Second).Seconds(),
			(180 * time.Second).Seconds(),
			(190 * time.Second).Seconds(),
			(200 * time.Second).Seconds(),
			(210 * time.Second).Seconds(),
			(220 * time.Second).Seconds(),
			(230 * time.Second).Seconds(),
			(240 * time.Second).Seconds(),
			(10 * time.Minute).Seconds(),
			(20 * time.Minute).Seconds(),
			(30 * time.Minute).Seconds(),
			(40 * time.Minute).Seconds(),
			(50 * time.Minute).Seconds(),
			(60 * time.Minute).Seconds(),
		),
	}

	// statsViews is a list of all stats views for
	// measurements emitted by this package.
	statsViews = []*view.View{
		waypointOperationDurations,
		waypointOperationCounts,
	}
)

func MeasureOperation(ctx context.Context, lastWriteAt time.Time, operationName string) {
	_ = stats.RecordWithTags(ctx, []tag.Mutator{
		tag.Upsert(TagOperation, operationName),
	}, operationDurationMeasure.M(time.Since(lastWriteAt).Seconds()))
}

func CountOperation(ctx context.Context, operationName string) {
	_ = stats.RecordWithTags(ctx, []tag.Mutator{
		tag.Upsert(TagOperation, operationName),
	}, operationCountMeasure.M(1))
}
