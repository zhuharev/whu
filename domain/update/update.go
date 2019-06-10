package update

import (
	"encoding/json"

	whPkg "github.com/zhuharev/whu/domain/webhook"
)

type Update struct {
	ID      int             `storm:"id,increment" json:"id"`
	WH      string          `storm:"index" json:"-"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type Repo interface {
	Save(string, []byte) error
	Get(wh string, offset int) ([]Update, error)
	GetUpdatesCount(id string) (int, error)
}

type UseCase interface {
	Save(string, []byte) error
	Get(wh string, offset int) ([]Update, error)
}

func New(whRepo whPkg.Repo, repo Repo) UseCase {
	return &uc{whRepo, repo}
}

type uc struct {
	whRepo whPkg.Repo
	repo   Repo
}

func (uc *uc) Save(id string, payload []byte) error {
	return uc.repo.Save(id, payload)
}

func (uc *uc) Get(id string, offset int) ([]Update, error) {
	wh, err := uc.whRepo.Get(id)
	if err != nil {
		return nil, err
	}

	count, err := uc.repo.GetUpdatesCount(id)
	if err != nil {
		return nil, err
	}

	if offset > count {
		offset = count
	}

	if wh.LastOffset < offset {
		wh.LastOffset = offset
		err = uc.whRepo.UpdateLastOffset(id, offset)
		if err != nil {
			return nil, err
		}
	}

	return uc.repo.Get(id, wh.LastOffset)
}
