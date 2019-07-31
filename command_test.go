// Copyright (c) 2016 - The Event Horizon authors.
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

	"github.com/google/uuid"
	eh "github.com/looplab/eventhorizon"
)

func Test_CreateCommand(t *testing.T) {
	cmd, err := eh.CreateCommand(TestCommandRegisterType)
	if err != eh.ErrCommandNotRegistered {
		t.Error("there should be a command not registered error:", err)
	}

	eh.RegisterCommand(func() eh.Command { return &TestCommandRegister{} })

	cmd, err = eh.CreateCommand(TestCommandRegisterType)
	if err != nil {
		t.Error("there should be no error:", err)
	}
	if cmd.CommandType() != TestCommandRegisterType {
		t.Error("the command type should be correct:", cmd.CommandType())
	}
}

func Test_RegisterCommandEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: attempt to register empty command type" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterCommand(func() eh.Command { return &TestCommandRegisterEmpty{} })
}

func Test_RegisterCommandNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: created command is nil" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterCommand(func() eh.Command { return nil })
}

func Test_RegisterCommandTwice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: registering duplicate types for \"TestCommandRegisterTwice\"" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterCommand(func() eh.Command { return &TestCommandRegisterTwice{} })
	eh.RegisterCommand(func() eh.Command { return &TestCommandRegisterTwice{} })
}

func Test_UnregisterCommandEmptyName(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: attempt to unregister empty command type" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.UnregisterCommand(TestCommandUnregisterEmptyType)
}

func Test_UnregisterCommandTwice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil || r != "eventhorizon: unregister of non-registered type \"TestCommandUnregisterTwice\"" {
			t.Error("there should have been a panic:", r)
		}
	}()
	eh.RegisterCommand(func() eh.Command { return &TestCommandUnregisterTwice{} })
	eh.UnregisterCommand(TestCommandUnregisterTwiceType)
	eh.UnregisterCommand(TestCommandUnregisterTwiceType)
}

func Test_CheckCommand(t *testing.T) {
	// Check all fields.
	err := eh.CheckCommand(&TestCommandFields{uuid.New().String(), "command1"})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Missing required string value.
	err = eh.CheckCommand(&TestCommandStringValue{TestID: uuid.New().String()})
	if err == nil || err.Error() != "missing field: Content" {
		t.Error("there should be a missing field error:", err)
	}

	// Missing required int value.
	err = eh.CheckCommand(&TestCommandIntValue{TestID: uuid.New().String()})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Missing required float value.
	err = eh.CheckCommand(&TestCommandFloatValue{TestID: uuid.New().String()})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Missing required bool value.
	err = eh.CheckCommand(&TestCommandBoolValue{TestID: uuid.New().String()})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Missing required slice.
	err = eh.CheckCommand(&TestCommandSlice{TestID: uuid.New().String()})
	if err == nil || err.Error() != "missing field: Slice" {
		t.Error("there should be a missing field error:", err)
	}

	// Missing required map.
	err = eh.CheckCommand(&TestCommandMap{TestID: uuid.New().String()})
	if err == nil || err.Error() != "missing field: Map" {
		t.Error("there should be a missing field error:", err)
	}

	// Missing required struct.
	err = eh.CheckCommand(&TestCommandStruct{TestID: uuid.New().String()})
	if err == nil || err.Error() != "missing field: Struct" {
		t.Error("there should be a missing field error:", err)
	}

	// Missing required time.
	err = eh.CheckCommand(&TestCommandTime{TestID: uuid.New().String()})
	if err == nil || err.Error() != "missing field: Time" {
		t.Error("there should be a missing field error:", err)
	}

	// Missing optional field.
	err = eh.CheckCommand(&TestCommandOptional{TestID: uuid.New().String()})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Missing private field.
	err = eh.CheckCommand(&TestCommandPrivate{TestID: uuid.New().String()})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Check all array fields.
	err = eh.CheckCommand(&TestCommandArray{uuid.New().String(), [1]string{"string"}, [1]int{0}, [1]struct{ Test string }{struct{ Test string }{"struct"}}})
	if err != nil {
		t.Error("there should be no error:", err)
	}

	// Empty array field.
	err = eh.CheckCommand(&TestCommandArray{uuid.New().String(), [1]string{""}, [1]int{0}, [1]struct{ Test string }{struct{ Test string }{"struct"}}})
	if err == nil || err.Error() != "missing field: StringArray" {
		t.Error("there should be a missing field error:", err)
	}
}

// Mocks for Register/Unregister.

const (
	TestCommandRegisterType        eh.CommandType = "TestCommandRegister"
	TestCommandRegisterEmptyType   eh.CommandType = ""
	TestCommandRegisterTwiceType   eh.CommandType = "TestCommandRegisterTwice"
	TestCommandUnregisterEmptyType eh.CommandType = ""
	TestCommandUnregisterTwiceType eh.CommandType = "TestCommandUnregisterTwice"

	TestAggregateType eh.AggregateType = "TestAggregate"
)

type TestCommandRegister struct{}

var _ = eh.Command(TestCommandRegister{})

func (a TestCommandRegister) AggregateID() eh.ID              { return eh.NilID }
func (a TestCommandRegister) AggregateType() eh.AggregateType { return TestAggregateType }
func (a TestCommandRegister) CommandType() eh.CommandType     { return TestCommandRegisterType }

type TestCommandRegisterEmpty struct{}

var _ = eh.Command(TestCommandRegisterEmpty{})

func (a TestCommandRegisterEmpty) AggregateID() eh.ID              { return eh.NilID }
func (a TestCommandRegisterEmpty) AggregateType() eh.AggregateType { return TestAggregateType }
func (a TestCommandRegisterEmpty) CommandType() eh.CommandType     { return TestCommandRegisterEmptyType }

type TestCommandRegisterTwice struct{}

var _ = eh.Command(TestCommandRegisterTwice{})

func (a TestCommandRegisterTwice) AggregateID() eh.ID              { return eh.NilID }
func (a TestCommandRegisterTwice) AggregateType() eh.AggregateType { return TestAggregateType }
func (a TestCommandRegisterTwice) CommandType() eh.CommandType     { return TestCommandRegisterTwiceType }

type TestCommandUnregisterTwice struct{}

var _ = eh.Command(TestCommandUnregisterTwice{})

func (a TestCommandUnregisterTwice) AggregateID() eh.ID              { return eh.NilID }
func (a TestCommandUnregisterTwice) AggregateType() eh.AggregateType { return TestAggregateType }
func (a TestCommandUnregisterTwice) CommandType() eh.CommandType {
	return TestCommandUnregisterTwiceType
}

// Mocks for CheckCommand.

type TestCommandFields struct {
	TestID  eh.ID
	Content string
}

var _ = eh.Command(TestCommandFields{})

func (t TestCommandFields) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandFields) AggregateType() eh.AggregateType { return TestAggregateType }
func (t TestCommandFields) CommandType() eh.CommandType {
	return eh.CommandType("TestCommandFields")
}

type TestCommandStringValue struct {
	TestID  eh.ID
	Content string
}

var _ = eh.Command(TestCommandStringValue{})

func (t TestCommandStringValue) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandStringValue) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandStringValue) CommandType() eh.CommandType {
	return eh.CommandType("TestCommandStringValue")
}

type TestCommandIntValue struct {
	TestID  eh.ID
	Content int
}

var _ = eh.Command(TestCommandIntValue{})

func (t TestCommandIntValue) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandIntValue) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandIntValue) CommandType() eh.CommandType {
	return eh.CommandType("TestCommandIntValue")
}

type TestCommandFloatValue struct {
	TestID  eh.ID
	Content float32
}

var _ = eh.Command(TestCommandFloatValue{})

func (t TestCommandFloatValue) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandFloatValue) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandFloatValue) CommandType() eh.CommandType {
	return eh.CommandType("TestCommandFloatValue")
}

type TestCommandBoolValue struct {
	TestID  eh.ID
	Content bool
}

var _ = eh.Command(TestCommandBoolValue{})

func (t TestCommandBoolValue) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandBoolValue) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandBoolValue) CommandType() eh.CommandType {
	return eh.CommandType("TestCommandBoolValue")
}

type TestCommandSlice struct {
	TestID eh.ID
	Slice  []string
}

var _ = eh.Command(TestCommandSlice{})

func (t TestCommandSlice) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandSlice) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandSlice) CommandType() eh.CommandType     { return eh.CommandType("TestCommandSlice") }

type TestCommandMap struct {
	TestID eh.ID
	Map    map[string]string
}

var _ = eh.Command(TestCommandMap{})

func (t TestCommandMap) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandMap) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandMap) CommandType() eh.CommandType     { return eh.CommandType("TestCommandMap") }

type TestCommandStruct struct {
	TestID eh.ID
	Struct struct {
		Test string
	}
}

var _ = eh.Command(TestCommandStruct{})

func (t TestCommandStruct) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandStruct) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandStruct) CommandType() eh.CommandType     { return eh.CommandType("TestCommandStruct") }

type TestCommandTime struct {
	TestID eh.ID
	Time   time.Time
}

var _ = eh.Command(TestCommandTime{})

func (t TestCommandTime) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandTime) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandTime) CommandType() eh.CommandType     { return eh.CommandType("TestCommandTime") }

type TestCommandOptional struct {
	TestID  eh.ID
	Content string `eh:"optional"`
}

var _ = eh.Command(TestCommandOptional{})

func (t TestCommandOptional) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandOptional) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandOptional) CommandType() eh.CommandType {
	return eh.CommandType("TestCommandOptional")
}

type TestCommandPrivate struct {
	TestID  eh.ID
	private string
}

var _ = eh.Command(TestCommandPrivate{})

func (t TestCommandPrivate) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandPrivate) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandPrivate) CommandType() eh.CommandType     { return eh.CommandType("TestCommandPrivate") }

type TestCommandArray struct {
	TestID      eh.ID
	StringArray [1]string
	IntArray    [1]int
	StructArray [1]struct {
		Test string
	}
}

var _ = eh.Command(TestCommandArray{})

func (t TestCommandArray) AggregateID() eh.ID              { return t.TestID }
func (t TestCommandArray) AggregateType() eh.AggregateType { return eh.AggregateType("Test") }
func (t TestCommandArray) CommandType() eh.CommandType     { return eh.CommandType("TestCommandArray") }
