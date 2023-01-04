package telescope

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type entries struct {
	Uuid                 string  `json:"uuid,omitempty"`
	BatchId              string  `json:"batch_id,omitempty"`
	FamilyHash           *string `json:"family_hash,omitempty"`
	ShouldDisplayOnIndex int     `json:"should_display_on_index,omitempty"`
	Type                 string  `json:"type,omitempty"`
	Content              string  `json:"content,omitempty"`
	CreatedAt            string  `json:"created_at,omitempty"`
}

type tag struct {
	EntryUuid string `json:"entry_uuid"`
	Tag       string `json:"tag"`
}

type Type interface {
	BindType() string
	Handler(entry *logrus.Entry) (*entries, []tag)
}

func ToContent(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}
