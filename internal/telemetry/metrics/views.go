package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	MetricJobsView = &view.View{
		Name:        "catsbyjobs",
		Measure:     MJobs,
		Description: "The distribution of the job build latencies",

		// Latency in buckets:
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		// Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000),
		TagKeys:     []tag.Key{KeyJobType},
	}
	// stat names
	Latency = "repl/latency"

	Jobs = "waypoint-metrics"

	JobOperation = "operation"
	JobBuild     = "build"
	JobDeploy    = "deploy"
	JobReport    = "report"
	JobRelease   = "release"
	JobInit      = "init"
	JobUp        = "up"
	JobPush      = "push"

	KeyMethod, _  = tag.NewKey("method")
	KeyService, _ = tag.NewKey("service")
	KeyStatus, _  = tag.NewKey("status")
	KeyError, _   = tag.NewKey("error")
	KeyJobType, _ = tag.NewKey("job_type")

	// The latency in milliseconds
	MLatencyMs = stats.Float64(Latency, "The latency in milliseconds per REPL loop", "ms")

	// Counts/groups the lengths of lines read in.
	MJobs = stats.Int64(Jobs, "Waypoint Job durations", "ms")

	// StatsViews is a list of all stats views for
	// measurements emitted by this package.
	StatsViews = []*view.View{
		MetricJobsView,
	}

// 	clusterIDTag = tag.MustNewKey("cluster_id")

// 	clusterTierTag = tag.MustNewKey("cluster_tier")

// 	clusterInconsistencyDetectedMeasure = stats.Int64(
// 		"hcp_vault_cluster_inconsistency_detected",
// 		"Sealed/Unsealed cluster inconsistency detected",
// 		stats.UnitDimensionless,
// 	)

// 	failedToGetClusterHealthMesaure = stats.Int64(
// 		"hcp_vault_failed_to_getclusterhealth",
// 		"Failed to get cluster health",
// 		stats.UnitDimensionless,
// 	)

// 	failedHealthCheckMeasure = stats.Int64(
// 		"hcp_vault_health_check_failed",
// 		"Failed to check health",
// 		stats.UnitDimensionless,
// 	)

// 	succeededHealthCheckMeasure = stats.Int64(
// 		"hcp_vault_health_check_succeeded",
// 		"Failed to check health",
// 		stats.UnitDimensionless,
// 	)

// 	performedHealthCheckMeasure = stats.Int64(
// 		"hcp_vault_health_check_performed",
// 		"Health Check Performed",
// 		stats.UnitDimensionless,
// 	)

// 	clusterHealthyMesaure = stats.Int64(
// 		"hcp_vault_cluster_healthy",
// 		"Cluster healthy",
// 		stats.UnitDimensionless,
// 	)

// 	clusterHealthTimeSinceLastWrite = stats.Float64(
// 		"hcp_vault_cluster_health_check_time_since_last_write",
// 		"The number of seconds since the last write from the previous health check",
// 		stats.UnitSeconds,
// 	)

// 	clusterInconsistencyDetectedView = &view.View{
// 		Name:        clusterInconsistencyDetectedMeasure.Name(),
// 		Description: clusterInconsistencyDetectedMeasure.Description(),
// 		TagKeys:     []tag.Key{clusterIDTag},
// 		Measure:     clusterInconsistencyDetectedMeasure,
// 		Aggregation: view.Count(),
// 	}

// 	failedToGetClusterHealthInformationView = &view.View{
// 		Name:        failedToGetClusterHealthMesaure.Name(),
// 		Description: failedToGetClusterHealthMesaure.Description(),
// 		TagKeys:     []tag.Key{clusterIDTag},
// 		Measure:     failedToGetClusterHealthMesaure,
// 		Aggregation: view.Count(),
// 	}

// 	failedHealthCheckView = &view.View{
// 		Name:        failedHealthCheckMeasure.Name(),
// 		Description: failedHealthCheckMeasure.Description(),
// 		TagKeys:     []tag.Key{clusterTierTag},
// 		Measure:     failedHealthCheckMeasure,
// 		Aggregation: view.Count(),
// 	}

// 	succeededHealthCheckView = &view.View{
// 		Name:        succeededHealthCheckMeasure.Name(),
// 		Description: succeededHealthCheckMeasure.Description(),
// 		TagKeys:     []tag.Key{clusterTierTag},
// 		Measure:     succeededHealthCheckMeasure,
// 		Aggregation: view.Count(),
// 	}

// 	performedHealthCheckView = &view.View{
// 		Name:        performedHealthCheckMeasure.Name(),
// 		Description: performedHealthCheckMeasure.Description(),
// 		TagKeys:     []tag.Key{clusterTierTag},
// 		Measure:     performedHealthCheckMeasure,
// 		Aggregation: view.Count(),
// 	}

// 	clusterHealthTimeSinceLastWriteView = &view.View{
// 		Name:        clusterHealthTimeSinceLastWrite.Name(),
// 		Description: clusterHealthTimeSinceLastWrite.Description(),
// 		TagKeys:     []tag.Key{clusterTierTag},
// 		Measure:     clusterHealthTimeSinceLastWrite,
// 		// add a custom distribution bucket of 10 second intervals between
// 		// 0 and 240 seconds, then 10 minute intervals up to 60 minutes
// 		Aggregation: view.Distribution(
// 			(10 * time.Second).Seconds(),
// 			(20 * time.Second).Seconds(),
// 			(30 * time.Second).Seconds(),
// 			(40 * time.Second).Seconds(),
// 			(50 * time.Second).Seconds(),
// 			(60 * time.Second).Seconds(),
// 			(70 * time.Second).Seconds(),
// 			(80 * time.Second).Seconds(),
// 			(90 * time.Second).Seconds(),
// 			(100 * time.Second).Seconds(),
// 			(110 * time.Second).Seconds(),
// 			(120 * time.Second).Seconds(),
// 			(130 * time.Second).Seconds(),
// 			(140 * time.Second).Seconds(),
// 			(150 * time.Second).Seconds(),
// 			(160 * time.Second).Seconds(),
// 			(170 * time.Second).Seconds(),
// 			(180 * time.Second).Seconds(),
// 			(190 * time.Second).Seconds(),
// 			(200 * time.Second).Seconds(),
// 			(210 * time.Second).Seconds(),
// 			(220 * time.Second).Seconds(),
// 			(230 * time.Second).Seconds(),
// 			(240 * time.Second).Seconds(),
// 			(10 * time.Minute).Seconds(),
// 			(20 * time.Minute).Seconds(),
// 			(30 * time.Minute).Seconds(),
// 			(40 * time.Minute).Seconds(),
// 			(50 * time.Minute).Seconds(),
// 			(60 * time.Minute).Seconds(),
// 		),
// 	}

// 	clusterHealthyView = &view.View{
// 		Name:        clusterHealthyMesaure.Name(),
// 		Description: clusterHealthyMesaure.Description(),
// 		TagKeys:     []tag.Key{clusterIDTag},
// 		Measure:     clusterHealthyMesaure,
// 		Aggregation: view.Count(),
// 	}
)
