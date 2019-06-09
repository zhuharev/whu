package update

import "encoding/json"

type Update struct {
	ID      int             `storm:"id,increment" json:"id"`
	WH      string          `storm:"index" json:"-"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Repo interface {
	Save(string, []byte) error
	Get(wh string, offset int) ([]Update, error)
}
