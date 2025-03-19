package models

import (
	"encoding/json"
	"errors"

	"github.com/futugyou/yomawari/kernel-memory/abstractions"
)

type TagCollection struct {
	datas map[string][]string
}

func (tc TagCollection) MarshalJSON() ([]byte, error) {
	if tc.datas == nil {
		tc.datas = make(map[string][]string)
	}
	return json.Marshal(tc.datas)
}

func (tc *TagCollection) UnmarshalJSON(data []byte) error {
	if tc == nil {
		return errors.New("TagCollection: UnmarshalJSON on nil pointer")
	}
	return json.Unmarshal(data, &tc.datas)
}

func NewTagCollection() *TagCollection {
	return &TagCollection{datas: make(map[string][]string)}
}

func (t *TagCollection) Add(tag string, value []string) {
	if (t) == nil {
		t = NewTagCollection()
	}
	t.datas[tag] = value
}

func (t *TagCollection) AddOrAppend(tag string, value string) {
	if (t) == nil {
		t = NewTagCollection()
	}

	if v, exists := t.datas[tag]; !exists {
		t.datas[tag] = []string{value}
	} else {
		v = append(v, value)
		t.datas[tag] = v
	}
}

func (t *TagCollection) Remove(tag string) {
	if t == nil {
		return
	}
	delete(t.datas, tag)
}

func (t *TagCollection) Keys() []string {
	if t == nil {
		return []string{}
	}
	keys := []string{}
	for k := range t.datas {
		keys = append(keys, k)
	}
	return keys
}

func (t *TagCollection) Values() [][]string {
	if t == nil {
		return [][]string{}
	}
	values := [][]string{}
	for k := range t.datas {
		values = append(values, t.datas[k])
	}
	return values
}

func (m *TagCollection) Get(key string) ([]string, bool) {
	if m == nil {
		return nil, false
	}
	value, exists := m.datas[key]
	return value, exists
}

func (t *TagCollection) Count() int {
	if t == nil {
		return 0
	}
	return len(t.datas)
}

func (t *TagCollection) Clear() {
	if t == nil {
		return
	}
	t.datas = make(map[string][]string)
}

func (t *TagCollection) AddSyntheticTag(value string) {
	t.AddOrAppend(abstractions.ReservedSyntheticTypeTag, value)
}
