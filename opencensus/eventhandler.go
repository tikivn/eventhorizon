package opencensus

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/internal/oc"
)

type EventHandler struct {
	tracer      *oc.Tracer
	handleEvent string
	eh.EventHandler
}

func NewEventHandler(handler eh.EventHandler) *EventHandler {
	return &EventHandler{
		tracer:       newTracer(handler),
		handleEvent:  fmt.Sprintf("EventHandler(%s).HandleEvent", handler.HandlerType()),
		EventHandler: handler,
	}
}

func (h *EventHandler) HandleEvent(ctx context.Context, event eh.Event) (err error) {
	ctx = h.tracer.Start(ctx, h.handleEvent)
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
