package oc

import (
	"fmt"

	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

func LatencyMeasure(pkg string) *stats.Float64Measure {
	return stats.Float64(
		pkg+"/latency",
		"Latency of method call",
		stats.UnitMilliseconds)
}

func MessageMeasure(pkg string) *stats.Int64Measure {
	return stats.Int64(
		pkg+"/message",
		"Total of messages processed",
		stats.UnitNone)
}

var (
	MethodKey   = MustKey("eh_method")
	StatusKey   = MustKey("eh_status")
	ProviderKey = MustKey("eh_provider")
)

func MustKey(name string) tag.Key {
	k, err := tag.NewKey(name)
	if err != nil {
		panic(fmt.Sprintf("tag.NewKey(%q): %v", name, err))
	}
	return k
}

// Views returns the views supported by Go CDK APIs.
func Views(pkg string, latencyMeasure *stats.Float64Measure) []*view.View {
	return []*view.View{
		{
			Name:        pkg + "/completed_calls",
			Measure:     latencyMeasure,
			Description: "Count of method calls by provider, method and status.",
			TagKeys:     []tag.Key{ProviderKey, MethodKey, StatusKey},
			Aggregation: view.Count(),
		},
		{
			Name:        pkg + "/latency",
			Measure:     latencyMeasure,
			Description: "Distribution of method latency, by provider and method.",
			TagKeys:     []tag.Key{ProviderKey, MethodKey},
			Aggregation: ocgrpc.DefaultMillisecondsDistribution,
		},
	}
}
