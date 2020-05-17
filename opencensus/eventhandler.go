package opencensus

import (
	"context"

	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/internal/oc"
)

type EventHandler struct {
	tracer *oc.Tracer
	eh.EventHandler
}

func NewEventHandler(handler eh.EventHandler) *EventHandler {
	return &EventHandler{
		tracer:       newTracer(handler),
		EventHandler: handler,
	}
}

func (h *EventHandler) HandleEvent(ctx context.Context, event eh.Event) (err error) {
	ctx = h.tracer.Start(ctx, "EventHandler.HandleEvent")
	ctx, err = tag.New(ctx,
		tag.Upsert(oc.AggregateTypeKey, (string)(event.AggregateType())),
		tag.Upsert(oc.HandlerTypeKey, (string)(h.EventHandler.HandlerType())),
		tag.Upsert(oc.EventTypeKey, (string)(event.EventType())),
	)
	if err != nil {
		panic(err)
	}
	span := trace.FromContext(ctx)
	defer func() {
		span.AddAttributes(
			trace.StringAttribute("aggregateID", (string)(event.AggregateID().String())),
			trace.StringAttribute("aggregateType", (string)(event.AggregateType())),
			trace.StringAttribute("eventType", (string)(event.EventType())),
			trace.Int64Attribute("version", (int64)(event.Version())),
		)
		if err == nil {
			h.tracer.End(ctx, err, messageMeasure.M(1))
		}
		h.tracer.End(ctx, err)
	}()

	return h.EventHandler.HandleEvent(ctx, event)
}
