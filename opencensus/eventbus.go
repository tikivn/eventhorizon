package opencensus

import (
	"context"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/internal/oc"
)

type TraceEventBus struct {
	tracer *oc.Tracer
	eh.EventBus
}

func NewEventBus(bus eh.EventBus) *TraceEventBus {
	return &TraceEventBus{
		tracer:   newTracer(bus),
		EventBus: bus,
	}
}

// PublishEvent publishes the event on the bus.
func (b *TraceEventBus) PublishEvent(ctx context.Context, event eh.Event) (err error) {
	ctx = b.tracer.Start(ctx, "EventBus.PublishEvent")
	ctx, err = tag.New(ctx,
		tag.Upsert(AggregateTypeKey, (string)(event.AggregateType())),
		tag.Upsert(EventTypeKey, (string)(event.EventType())),
	)
	if err != nil {
		panic(err)
	}
	span := trace.FromContext(ctx)
	span.AddAttributes(
		trace.StringAttribute("aggregateID", event.AggregateID().String()),
		trace.StringAttribute("aggregateType", (string)(event.AggregateType())),
		trace.Int64Attribute("eventVersion", (int64)(event.Version())),
	)
	defer func() {
		b.tracer.End(ctx, err, messageMeasure.M(1))
	}()
	return b.EventBus.PublishEvent(ctx, event)
}

// AddHandler adds a handler for an event. Panics if either the matcher
// or handler is nil or the handler is already added.
func (b *TraceEventBus) AddHandler(m eh.EventMatcher, h eh.EventHandler) {
	b.EventBus.AddHandler(m, NewEventHandler(h))
}

// AddObserver adds an observer. Panics if the observer is nil or the observer
// is already added.
func (b *TraceEventBus) AddObserver(m eh.EventMatcher, h eh.EventHandler) {
	b.EventBus.AddObserver(m, NewEventHandler(h))
}
