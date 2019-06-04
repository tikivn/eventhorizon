package opencensus

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.opencensus.io/trace"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/internal/oc"
)

type TraceStore struct {
	tracer *oc.Tracer
	store  eh.EventStore
}

func NewEventStore(store eh.EventStore) *TraceStore {
	return &TraceStore{
		tracer: newTracer(store),
		store:  store,
	}
}

func (s *TraceStore) Load(
	ctx context.Context, id uuid.UUID, aggregateType eh.AggregateType,
) (events []eh.Event, err error) {
	ctx = s.tracer.Start(ctx, "EventStore.Load")
	span := trace.FromContext(ctx)
	span.AddAttributes(
		trace.StringAttribute("aggregateID", id.String()),
		trace.StringAttribute("aggregateType", (string)(aggregateType)),
	)
	defer func() {
		lenEvents := (int64)(len(events))
		span.AddAttributes(
			trace.Int64Attribute("events", lenEvents),
		)
		if err == nil {
			s.tracer.End(ctx, err, messageMeasure.M(lenEvents))
		} else {
			s.tracer.End(ctx, err)
		}
	}()

	return s.store.Load(ctx, id, aggregateType)
}

func (s *TraceStore) Save(
	ctx context.Context, events []eh.Event, originalVersion int,
) (err error) {
	ctx = s.tracer.Start(ctx, "EventStore.Save")
	span := trace.FromContext(ctx)
	lenEvents := (int64)(len(events))
	span.AddAttributes(
		trace.Int64Attribute("events", lenEvents),
		trace.Int64Attribute("originalVersion", (int64)(originalVersion)),
	)
	defer func() {
		if err == nil {
			s.tracer.End(ctx, err, messageMeasure.M(lenEvents))
		} else {
			s.tracer.End(ctx, err)
		}
	}()

	return s.store.Save(ctx, events, originalVersion)
}

func (s *TraceStore) LoadEvent(ctx context.Context, eventIDs []uuid.UUID) ([]eh.Event, error) {
	type es interface {
		LoadEvent(ctx context.Context, eventIDs []uuid.UUID) ([]eh.Event, error)
	}

	if newES, ok := s.store.(es); ok {
		return newES.LoadEvent(ctx, eventIDs)
	}

	return nil, errors.New("Nil Interface")
}
