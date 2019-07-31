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

package memory_test

import (
	"context"
	"testing"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/mocks"
	"github.com/looplab/eventhorizon/repo"
	"github.com/looplab/eventhorizon/repo/memory"
)

func Test_ReadRepo(t *testing.T) {
	r := memory.NewRepo()
	if r == nil {
		t.Error("there should be a repository")
	}
	if r.Parent() != nil {
		t.Error("the parent repo should be nil")
	}

	// Repo with default namespace.
	repo.AcceptanceTest(t, context.Background(), r)

	// Repo with other namespace
	ctx := eh.NewContextWithNamespace(context.Background(), "ns")
	repo.AcceptanceTest(t, ctx, r)

}

func Test_Repository(t *testing.T) {
	if r := memory.Repository(nil); r != nil {
		t.Error("the parent repository should be nil:", r)
	}

	inner := &mocks.Repo{}
	if r := memory.Repository(inner); r != nil {
		t.Error("the parent repository should be nil:", r)
	}

	repo := memory.NewRepo()
	outer := &mocks.Repo{ParentRepo: repo}
	if r := memory.Repository(outer); r != repo {
		t.Error("the parent repository should be correct:", r)
	}
}
