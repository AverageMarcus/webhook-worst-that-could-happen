package patcher

import "encoding/json"

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}

type Patches struct {
	Entries []JSONPatchEntry
}

func New() *Patches {
	return &Patches{}
}

func (p *Patches) Add(patch JSONPatchEntry) *Patches {
	p.Entries = append(p.Entries, patch)
	return p
}

func (p *Patches) ToBytes() ([]byte, error) {
	return json.Marshal(&p.Entries)
}

func ToValueString(input string) json.RawMessage {
	return []byte(`"` + input + `"`)
}

func ToValueMap(input map[string]string) json.RawMessage {
	bytes, _ := json.Marshal(input)
	return bytes
}
