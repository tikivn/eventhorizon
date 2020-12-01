package events

import (
	"context"

	"github.com/google/uuid"
	eh "github.com/tikivn/eventhorizon"
)

type AggregateQuerier struct {
	store eh.EventStore
}

// Load implements the Load method of the eventhorizon.AggregateQuerier interface.
// It loads an aggregate from the event store by creating a new aggregate of the
// type with the ID and then applies all events to it, thus making it the most
// current version of the aggregate.
func (r *AggregateQuerier) Load(ctx context.Context, aggregateType eh.AggregateType, id uuid.UUID) (eh.Aggregate, error) {
	agg, err := eh.CreateAggregate(aggregateType, id)
	if err != nil {
		return nil, err
	}
	a, ok := agg.(Aggregate)
	if !ok {
		return nil, ErrInvalidAggregateType
	}

	events, err := r.store.Load(ctx, a.EntityID())
	if err != nil {
		return nil, err
	}

	if err := r.applyEvents(ctx, a, events); err != nil {
		return nil, err
	}

	return a, nil
}

func (r *AggregateQuerier) applyEvents(ctx context.Context, a Aggregate, events []eh.Event) error {
	for _, event := range events {
		if event.AggregateType() != a.AggregateType() {
			return ErrMismatchedEventType
		}

		if err := a.ApplyEvent(ctx, event); err != nil {
			return ApplyEventError{
				Event: event,
				Err:   err,
			}
		}
		a.IncrementVersion()
	}

	return nil
}
