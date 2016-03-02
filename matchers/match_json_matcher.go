package matchers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
)

type MatchJSONMatcher struct {
	JSONToMatch interface{}
}

func (matcher *MatchJSONMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, expectedString, err := matcher.prettyPrint(actual)
	if err != nil {
		return false, err
	}

	var aval interface{}
	var eval interface{}

	// this is guarded by prettyPrint
	json.Unmarshal([]byte(actualString), &aval)
	json.Unmarshal([]byte(expectedString), &eval)

	return reflect.DeepEqual(aval, eval), nil
}

func (matcher *MatchJSONMatcher) FailureMessage(actual interface{}) (message string) {
	var aval interface{}
	var eval interface{}

	actualString, expectedString, _ := matcher.prettyPrint(actual)
	json.Unmarshal([]byte(actualString), &aval)
	json.Unmarshal([]byte(expectedString), &eval)

	actualMap, actualOK := aval.(map[string]interface{})
	expectedMap, expectedOK := eval.(map[string]interface{})
	if actualOK && expectedOK {
		for key, _ := range actualMap {
			println("looking for", key, "in expectedMap")
			_, ok := expectedMap[key]
			if !ok {
				println("HEYOOOOO")
				return "zoinks!"
			}
		}

		for key, value := range expectedMap {
			maybe, ok := actualMap[key]
			fmt.Printf("looking for %s in actual map, we found %#v %v\n", key, maybe, ok)
			if !ok || maybe == nil {
				println("WOOOOOOOOOOAH")
				return format.Message("Expected", actualString, "to match json", expectedString, "but it lacks the field:", key, ":", value)
			}
		}
	}

	/* TODO : arrays, nested objects, types? edge cases with strings
	iterate through fields on actual
	for each field name
	verify that it exists within expected

	iterate through fields on expected
	for each field name
	verify that it exists within actual
	*/

	return format.Message(actualString, "to match JSON of", expectedString)
}

func (matcher *MatchJSONMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, expectedString, _ := matcher.prettyPrint(actual)
	return format.Message(actualString, "not to match JSON of", expectedString)
}

func (matcher *MatchJSONMatcher) prettyPrint(actual interface{}) (actualFormatted, expectedFormatted string, err error) {
	actualString, ok := toString(actual)
	if !ok {
		return "", "", fmt.Errorf("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got actual:\n%s", format.Object(actual, 1))
	}
	expectedString, ok := toString(matcher.JSONToMatch)
	if !ok {
		return "", "", fmt.Errorf("MatchJSONMatcher matcher requires a string, stringer, or []byte.  Got expected:\n%s", format.Object(matcher.JSONToMatch, 1))
	}

	abuf := new(bytes.Buffer)
	ebuf := new(bytes.Buffer)

	if err := json.Indent(abuf, []byte(actualString), "", "  "); err != nil {
		return "", "", fmt.Errorf("Actual '%s' should be valid JSON, but it is not.\nUnderlying error:%s", actualString, err)
	}

	if err := json.Indent(ebuf, []byte(expectedString), "", "  "); err != nil {
		return "", "", fmt.Errorf("Expected '%s' should be valid JSON, but it is not.\nUnderlying error:%s", expectedString, err)
	}

	return abuf.String(), ebuf.String(), nil
}
