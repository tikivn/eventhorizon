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
	"testing"
	"time"

	eh "github.com/looplab/eventhorizon"
)

func Test_MatchAny(t *testing.T) {
	m := eh.MatchAny()

	if !m(nil) {
		t.Error("match any should always match")
	}

	e := eh.NewEvent("test", nil, time.Now())
	if !m(e) {
		t.Error("match any should always match")
	}
}
func Test_MatchEvent(t *testing.T) {
	et := eh.EventType("test")
	m := eh.MatchEvent(et)

	if m(nil) {
		t.Error("match event should not match nil event")
	}

	e := eh.NewEvent(et, nil, time.Now())
	if !m(e) {
		t.Error("match event should match the event")
	}

	e = eh.NewEvent("other", nil, time.Now())
	if m(e) {
		t.Error("match event should not match the event")
	}
}

func Test_MatchAggregate(t *testing.T) {
	at := eh.AggregateType("test")
	m := eh.MatchAggregate(at)

	if m(nil) {
		t.Error("match aggregate should not match nil event")
	}

	e := eh.NewEventForAggregate("test", nil, time.Now(), at, eh.NilID, 0)
	if !m(e) {
		t.Error("match aggregate should match the event")
	}

	e = eh.NewEventForAggregate("test", nil, time.Now(), "other", eh.NilID, 0)
	if m(e) {
		t.Error("match aggregate should not match the event")
	}
}

func Test_MatchAnyOf(t *testing.T) {
	et1 := eh.EventType("et1")
	et2 := eh.EventType("et2")
	m := eh.MatchAnyOf(
		eh.MatchEvent(et1),
		eh.MatchEvent(et2),
	)

	e := eh.NewEvent(et1, nil, time.Now())
	if !m(e) {
		t.Error("match any of should match the first event")
	}
	e = eh.NewEvent(et2, nil, time.Now())
	if !m(e) {
		t.Error("match any of should match the last event")
	}
}

func Test_MatchAnyEventOf(t *testing.T) {
	et1 := eh.EventType("test")
	et2 := eh.EventType("test")
	m := eh.MatchAnyEventOf(et1, et2)

	if m(nil) {
		t.Error("match any event of should not match nil event")
	}

	e1 := eh.NewEvent(et1, nil, time.Now())
	if !m(e1) {
		t.Error("match any event of should match the first event")
	}
	e2 := eh.NewEvent(et2, nil, time.Now())
	if !m(e2) {
		t.Error("match any event of should match the second event")
	}
}
