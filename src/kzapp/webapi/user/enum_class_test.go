package user

import (
	"encoding/json"
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestEnumClass(t *testing.T) {
	var data = EnumData{}
	err := json.Unmarshal([]byte(`{"enum": "User"}`), &data)
	if err != nil {
		t.Errorf("json.Unmarshal() error: %v", err)
	}
	fmt.Println(data.Enum)
}

func TestEnumClass_invalid(t *testing.T) {
	var data = EnumData{}
	err := json.Unmarshal([]byte(`{"enum": "Happy"}`), &data)
	t.Log(err.Error())
	assert.Error(t, err, "invalid enum value")
}