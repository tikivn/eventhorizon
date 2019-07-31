// Copyright (c) 2014 - The Event Horizon authors.
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
	"testing"

	"github.com/google/uuid"

	eh "github.com/looplab/eventhorizon"
)

func Test_CreateAggregate(t *testing.T) {
	id := uuid.New().String()
	aggregate, err := eh.CreateAggregate(TestAggregateRegisterType, id)
	if err != eh.ErrAggregateNotRegistered {
		t.Error("there should be a aggregate not registered error:", err)
	}

	eh.RegisterAggregate(func(id eh.ID) eh.Aggregate {
		return &TestAggregateRegister{id: id}
	})

	aggregate, err = eh.CreateAggregate(TestAggregateRegisterType, id)
	if err != nil {
		t.Error("there should be no error:", err)
	}
	// NOTE: The aggregate type used to register with is another than the aggregate!
	if aggregate.AggregateType() != TestAggregateRegisterType {
		t.Error("the aggregate type should be correct:", aggregate.AggregateType())
	}
	if aggregate.EntityID() != id {
		t.Error("the ID should be correct:", aggregate.EntityID())
	}
}

func Test_RegisterAggregateEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: attempt to register empty aggregate type" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterAggregate(func(id eh.ID) eh.Aggregate {
		return &TestAggregateRegisterEmpty{id: id}
	})
}

func Test_RegisterAggregateNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: created aggregate is nil" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterAggregate(func(id eh.ID) eh.Aggregate { return nil })
}

func Test_RegisterAggregateTwice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: registering duplicate types for \"TestAggregateRegisterTwice\"" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterAggregate(func(id eh.ID) eh.Aggregate {
		return &TestAggregateRegisterTwice{id: id}
	})
	eh.RegisterAggregate(func(id eh.ID) eh.Aggregate {
		return &TestAggregateRegisterTwice{id: id}
	})
}

const (
	TestAggregateRegisterType      eh.AggregateType = "TestAggregateRegister"
	TestAggregateRegisterEmptyType eh.AggregateType = ""
	TestAggregateRegisterTwiceType eh.AggregateType = "TestAggregateRegisterTwice"
)

type TestAggregateRegister struct {
	id eh.ID
}

var _ = eh.Aggregate(&TestAggregateRegister{})

func (a *TestAggregateRegister) EntityID() eh.ID { return a.id }

func (a *TestAggregateRegister) AggregateType() eh.AggregateType {
	return TestAggregateRegisterType
}
func (a *TestAggregateRegister) HandleCommand(ctx context.Context, cmd eh.Command) error {
	return nil
}

type TestAggregateRegisterEmpty struct {
	id eh.ID
}

var _ = eh.Aggregate(&TestAggregateRegisterEmpty{})

func (a *TestAggregateRegisterEmpty) EntityID() eh.ID { return a.id }

func (a *TestAggregateRegisterEmpty) AggregateType() eh.AggregateType {
	return TestAggregateRegisterEmptyType
}
func (a *TestAggregateRegisterEmpty) HandleCommand(ctx context.Context, cmd eh.Command) error {
	return nil
}

type TestAggregateRegisterTwice struct {
	id eh.ID
}

var _ = eh.Aggregate(&TestAggregateRegisterTwice{})

func (a *TestAggregateRegisterTwice) EntityID() eh.ID { return a.id }

func (a *TestAggregateRegisterTwice) AggregateType() eh.AggregateType {
	return TestAggregateRegisterTwiceType
}
func (a *TestAggregateRegisterTwice) HandleCommand(ctx context.Context, cmd eh.Command) error {
	return nil
}
