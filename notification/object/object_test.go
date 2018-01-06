package object

import (
	"testing"
	. "github.com/onsi/gomega"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func TestNotificationObject_Deserialization(t *testing.T) {
	t.Run("JSON", jsonDeserialization)
	t.Run("Mapstructure", mapstructureDeserialization)
}

func TestNotificationObject_Serialization(t *testing.T) {
	RegisterTestingT(t)

	t.Run("JSON", jsonSerialization)
}

func jsonSerialization(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		Object interface{}
		Json   string
	}{
		{Object: NewDataObject("module_name", "value"), Json: `{"module":"module_name","value":"value"}`},
		{Object: NewDataObject("module_name", ""), Json: `{"module":"module_name","value":""}`},
		{Object: NewDataObject("", "value"), Json: `{"module":"","value":"value"}`},
		{Object: DataObject{}, Json: `{"module":"","value":null}`},

		{Object: NewErrorObject("origin", fmt.Errorf("why ?")), Json: `{"origin":"origin","reason":"why ?"}`},
		{Object: NewErrorObject("origin", nil), Json: `{"origin":"origin","reason":""}`},
		{Object: NewErrorObject("", fmt.Errorf("why ?")), Json: `{"origin":"","reason":"why ?"}`},
		{Object: ErrorObject{}, Json: `{"origin":"","reason":""}`},

		{Object: NewCommandObject("command_name", "one", "two"), Json: `{"command":"command_name","args":["one","two"]}`},
		{Object: NewCommandObject("command_name", "one"), Json: `{"command":"command_name","args":["one"]}`},
		{Object: NewCommandObject("command_name"), Json: `{"command":"command_name","args":[]}`},
		{Object: NewCommandObject("", "one", "two"), Json: `{"command":"","args":["one","two"]}`},
		{Object: CommandObject{}, Json: `{"command":"","args":null}`},
	}

	for _, tc := range testCases {
		res, err := json.Marshal(tc.Object)

		Expect(err).Should(Succeed())
		Expect(string(res)).Should(Equal(tc.Json))
	}
}

func jsonDeserialization(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		Json   string
		Object interface{}
		Result interface{}
	}{
		{Json: `{"module": "module_name", "value": "value"}`, Object: &DataObject{}, Result: NewDataObject("module_name", "value")},
		{Json: `{"module": "module_name"}`, Object: &DataObject{}, Result: NewDataObject("module_name", nil)},
		{Json: `{"value": "value"}`, Object: &DataObject{}, Result: NewDataObject("", "value")},
		{Json: `{}`, Object: &DataObject{}, Result: &DataObject{}},

		{Json: `{"origin": "origin", "reason": "why ?"}`, Object: &ErrorObject{}, Result: NewErrorObject("origin", fmt.Errorf("why ?"))},
		{Json: `{"origin": "origin"}`, Object: &ErrorObject{}, Result: NewErrorObject("origin", nil)},
		{Json: `{"reason": "why ?"}`, Object: &ErrorObject{}, Result: NewErrorObject("", fmt.Errorf("why ?"))},
		{Json: `{}`, Object: &ErrorObject{}, Result: &ErrorObject{}},

		{Json: `{"command": "command_name", "args": ["one", "two"]}`, Object: &CommandObject{}, Result: NewCommandObject("command_name", "one", "two")},
		{Json: `{"command": "command_name", "args": ["one"]}`, Object: &CommandObject{}, Result: NewCommandObject("command_name", "one")},
		{Json: `{"command": "command_name", "args": []}`, Object: &CommandObject{}, Result: NewCommandObject("command_name")},
		{Json: `{"command": "command_name"}`, Object: &CommandObject{}, Result: &CommandObject{Command: "command_name"}},
		{Json: `{"args": ["one", "two"]}`, Object: &CommandObject{}, Result: NewCommandObject("", "one", "two")},
		{Json: `{}`, Object: &CommandObject{}, Result: &CommandObject{}},
	}

	for _, tc := range testCases {
		json.Unmarshal([]byte(tc.Json), tc.Object)
		Expect(tc.Object).Should(Equal(tc.Result))
	}
}

func mapstructureDeserialization(t *testing.T) {
	RegisterTestingT(t)

	testCases := []struct {
		Mapstructure map[string]interface{}
		Object       interface{}
		Result       interface{}
	}{
		{Mapstructure: map[string]interface{}{"module": "module_name", "value": "value"}, Object: &DataObject{}, Result: NewDataObject("module_name", "value")},
		{Mapstructure: map[string]interface{}{"module": "module_name"}, Object: &DataObject{}, Result: NewDataObject("module_name", nil)},
		{Mapstructure: map[string]interface{}{"value": "value"}, Object: &DataObject{}, Result: NewDataObject("", "value")},
		{Mapstructure: map[string]interface{}{}, Object: &DataObject{}, Result: &DataObject{}},

		{Mapstructure: map[string]interface{}{"origin": "origin", "reason": "why ?"}, Object: &ErrorObject{}, Result: NewErrorObject("origin", fmt.Errorf("why ?"))},
		{Mapstructure: map[string]interface{}{"origin": "origin"}, Object: &ErrorObject{}, Result: NewErrorObject("origin", nil)},
		{Mapstructure: map[string]interface{}{"reason": "why ?"}, Object: &ErrorObject{}, Result: NewErrorObject("", fmt.Errorf("why ?"))},
		{Mapstructure: map[string]interface{}{}, Object: &ErrorObject{}, Result: &ErrorObject{}},

		{Mapstructure: map[string]interface{}{"command": "command_name", "args": []string{"one", "two"}}, Object: &CommandObject{}, Result: NewCommandObject("command_name", "one", "two")},
		{Mapstructure: map[string]interface{}{"command": "command_name", "args": []string{"one"}}, Object: &CommandObject{}, Result: NewCommandObject("command_name", "one")},
		{Mapstructure: map[string]interface{}{"command": "command_name", "args": []string{}}, Object: &CommandObject{}, Result: NewCommandObject("command_name")},
		{Mapstructure: map[string]interface{}{"command": "command_name"}, Object: &CommandObject{}, Result: &CommandObject{Command: "command_name"}},
		{Mapstructure: map[string]interface{}{"args": []string{"one", "two"}}, Object: &CommandObject{}, Result: NewCommandObject("", "one", "two")},
		{Mapstructure: map[string]interface{}{}, Object: &CommandObject{}, Result: &CommandObject{}},
	}

	for _, tc := range testCases {
		mapstructure.Decode(tc.Mapstructure, tc.Object)
		Expect(tc.Object).Should(Equal(tc.Result))
	}
}

func TestErrorObject_Why(t *testing.T) {
	RegisterTestingT(t)

	obj := NewErrorObject("origin").Why(fmt.Errorf("reason"))
	Expect(obj).Should(Equal(&ErrorObject{Origin: "origin", Reason: "reason"}))
}

func TestCommandObject_AddArgument(t *testing.T) {
	RegisterTestingT(t)

	obj := NewCommandObject("command_name").AddArgument("one", "two").AddArgument("three")
	Expect(obj).Should(Equal(&CommandObject{Command: "command_name", Args: []string{"one", "two", "three"}}))
}
