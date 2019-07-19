package webhook

type Webhook struct {
	ID         string `storm:"unique" json:"id,omitempty"`
	LastOffset int    `json:"last_offset,omitempty"`
	ProxyURL   string `json:"proxyUrl,omitempty`
}

type Repo interface {
	UpdateLastOffset(id string, offset int) error
	Create(id string) error
	Get(id string) (*Webhook, error)
	SetProxy(id string, url string) error
}
