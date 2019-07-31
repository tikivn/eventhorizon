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
	"errors"
	"testing"
	"time"

	eh "github.com/looplab/eventhorizon"
)

func Test_EventBusError(t *testing.T) {
	var testCases = []struct {
		name              string
		err               error
		event             eh.Event
		expectedErrorText string
	}{
		{
			"both non-nil",
			errors.New("some error"),
			eh.NewEvent("some event type", nil, time.Time{}),
			"some error: (some event type@0)",
		},
		{
			"error nil",
			nil,
			eh.NewEvent("some event type", nil, time.Time{}),
			"%!s(<nil>): (some event type@0)",
		},
		{
			"event nil",
			errors.New("some error"),
			nil,
			"some error: (%!s(<nil>))",
		},

		{
			"both nil",
			nil,
			nil,
			"%!s(<nil>): (%!s(<nil>))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			busError := eh.EventBusError{
				Err:   tc.err,
				Event: tc.event,
			}

			if busError.Error() != tc.expectedErrorText {
				t.Errorf(
					"expected '%s', got '%s'",
					tc.expectedErrorText,
					busError.Error())
			}
		})
	}
}
