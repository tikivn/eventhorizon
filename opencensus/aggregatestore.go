package opencensus

import (
	"context"

	"github.com/google/uuid"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/internal/oc"
)

type TraceAggregateStore struct {
	tracer *oc.Tracer
	eh.AggregateStore
}

func NewAggregateStore(store eh.AggregateStore) *TraceAggregateStore {
	return &TraceAggregateStore{
		tracer:         newTracer(store),
		AggregateStore: store,
	}
}

// Load loads the most recent version of an aggregate with a type and id.
func (s *TraceAggregateStore) Load(ctx context.Context, aggregateType eh.AggregateType, id uuid.UUID) (aggregate eh.Aggregate, err error) {
	ctx = s.tracer.Start(ctx, "AggregateStore.Load")
	ctx, err = tag.New(ctx, tag.Upsert(AggregateTypeKey, (string)(aggregateType)))
	if err != nil {
		panic(err)
	}
	span := trace.FromContext(ctx)
	span.AddAttributes(
		trace.StringAttribute("aggregateID", id.String()),
		trace.StringAttribute("aggregateType", (string)(aggregateType)),
	)
	defer func() {
		s.tracer.End(ctx, err)
	}()

	return s.AggregateStore.Load(ctx, aggregateType, id)
}

// Save saves the uncommittend events for an aggregate.
func (s *TraceAggregateStore) Save(ctx context.Context, aggregate eh.Aggregate) (err error) {
	ctx = s.tracer.Start(ctx, "AggregateStore.Save")
	ctx, err = tag.New(ctx, tag.Upsert(AggregateTypeKey, (string)(aggregate.AggregateType())))
	if err != nil {
		panic(err)
	}
	span := trace.FromContext(ctx)
	span.AddAttributes(
		trace.StringAttribute("aggregateID", aggregate.EntityID().String()),
		trace.StringAttribute("aggregateType", (string)(aggregate.AggregateType())),
	)
	defer func() {
		s.tracer.End(ctx, err)
	}()

	return s.AggregateStore.Save(ctx, aggregate)
}
