// Copyright (c) 2018 - The Event Horizon authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventhorizon_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	eh "github.com/looplab/eventhorizon"
)

func Test_EventHandlerFunc(t *testing.T) {
	events := []eh.Event{}
	h := eh.EventHandlerFunc(func(ctx context.Context, e eh.Event) error {
		events = append(events, e)
		return nil
	})
	if h.HandlerType() != eh.EventHandlerType(fmt.Sprintf("handler-func-%v", h)) {
		t.Error("the handler type should be correct:", h.HandlerType())
	}

	e := eh.NewEvent("test", nil, time.Now())
	h.HandleEvent(context.Background(), e)
	if !reflect.DeepEqual(events, []eh.Event{e}) {
		t.Error("the events should be correct")
		t.Log(events)
	}
}

func Test_EventHandlerMiddleware(t *testing.T) {
	order := []string{}
	middleware := func(s string) eh.EventHandlerMiddleware {
		return eh.EventHandlerMiddleware(func(h eh.EventHandler) eh.EventHandler {
			return eh.EventHandlerFunc(func(ctx context.Context, e eh.Event) error {
				order = append(order, s)
				return h.HandleEvent(ctx, e)
			})
		})
	}
	handler := func(ctx context.Context, e eh.Event) error {
		return nil
	}
	h := eh.UseEventHandlerMiddleware(eh.EventHandlerFunc(handler),
		middleware("first"),
		middleware("second"),
		middleware("third"),
	)
	h.HandleEvent(context.Background(), eh.NewEvent("test", nil, time.Now()))
	if !reflect.DeepEqual(order, []string{"first", "second", "third"}) {
		t.Error("the order of middleware should be correct")
		t.Log(order)
	}
}
