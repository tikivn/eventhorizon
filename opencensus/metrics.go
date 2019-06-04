package opencensus

import (
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/looplab/eventhorizon/internal/oc"
)

const pkgName = "github.com/tikivn/eventhorizon"

var (
	latencyMeasure = oc.LatencyMeasure(pkgName)
	messageMeasure = oc.MessageMeasure(pkgName)

	// OpenCensusViews are predefined views for OpenCensus metrics.
	// The views include counts and latency distributions for API method calls.
	// See the example at https://godoc.org/go.opencensus.io/stats/view for usage.
	OpenCensusViews = append(oc.Views(pkgName, latencyMeasure), &view.View{
		Name:        pkgName + "/message",
		Measure:     messageMeasure,
		Description: "Distribution of method latency, by provider and method.",
		TagKeys:     []tag.Key{oc.ProviderKey, oc.MethodKey, oc.AggregateTypeKey, oc.HandlerTypeKey, oc.EventTypeKey},
		Aggregation: ocgrpc.DefaultMessageCountDistribution,
	})
)

func newTracer(driver interface{}) *oc.Tracer {
	return &oc.Tracer{
		Package:        pkgName,
		Provider:       oc.ProviderName(driver),
		LatencyMeasure: latencyMeasure,
	}
}
