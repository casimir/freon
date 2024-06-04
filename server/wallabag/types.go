package wallabag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type IntBool bool

func (b *IntBool) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*b = v == 1
	return nil
}

func (b IntBool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal(1)
	}
	return json.Marshal(0)
}

type MagicInt struct {
	IsNil bool
	Value int
}

func (i *MagicInt) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if reflect.TypeOf(v) == nil {
		i.IsNil = true
		return nil
	}
	switch v := v.(type) {
	case int:
		i.IsNil = false
		i.Value = v
	case float64:
		i.IsNil = false
		i.Value = int(v)
	case string:
		i.IsNil = false
		i.Value, _ = strconv.Atoi(v)
	default:
		return fmt.Errorf("invalid type: %T", v)
	}
	return nil
}

var WallabagTimeLayout = "2006-01-02T15:04:05-0700"

type Time struct{ time.Time }

func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := time.Parse(WallabagTimeLayout, s)
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(WallabagTimeLayout))
}
