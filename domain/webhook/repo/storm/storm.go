package storm

import (
	"net/url"

	"github.com/asdine/storm"
	"github.com/zhuharev/whu/domain/webhook"
)

type repo struct {
	db *storm.DB
}

func New(db *storm.DB) webhook.Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) UpdateLastOffset(id string, offset int) error {
	return r.db.UpdateField(&webhook.Webhook{ID: id}, "LastOffset", offset)
}

func (r *repo) Create(id string) error {
	wh := webhook.Webhook{
		ID:         id,
		LastOffset: 0,
	}
	return r.db.Save(&wh)
}

func (r *repo) Get(id string) (*webhook.Webhook, error) {
	wh := webhook.Webhook{}
	err := r.db.One("ID", id, &wh)
	return &wh, err
}

func (r *repo) SetProxy(id string, u string) error {
	_, err := url.Parse(u)
	if err != nil {
		return err
	}
	wh, err := r.Get(id)
	if err != nil {
		return err
	}
	wh.ProxyURL = u
	return r.db.UpdateField(wh, "ProxyURL", u)
}
