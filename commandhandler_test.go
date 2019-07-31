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
	"reflect"
	"testing"

	eh "github.com/looplab/eventhorizon"
)

func Test_CommandHandlerMiddleware(t *testing.T) {
	order := []string{}
	middleware := func(s string) eh.CommandHandlerMiddleware {
		return eh.CommandHandlerMiddleware(func(h eh.CommandHandler) eh.CommandHandler {
			return eh.CommandHandlerFunc(func(ctx context.Context, cmd eh.Command) error {
				order = append(order, s)
				return h.HandleCommand(ctx, cmd)
			})
		})
	}
	handler := func(ctx context.Context, cmd eh.Command) error {
		return nil
	}
	h := eh.UseCommandHandlerMiddleware(eh.CommandHandlerFunc(handler),
		middleware("first"),
		middleware("second"),
		middleware("third"),
	)
	h.HandleCommand(context.Background(), TestCommand{})
	if !reflect.DeepEqual(order, []string{"first", "second", "third"}) {
		t.Error("the order of middleware should be correct")
		t.Log(order)
	}
}

type TestCommand struct{}

var _ = eh.Command(TestCommand{})

func (a TestCommand) AggregateID() eh.ID              { return eh.NilID }
func (a TestCommand) AggregateType() eh.AggregateType { return "test" }
func (a TestCommand) CommandType() eh.CommandType     { return "tes" }
