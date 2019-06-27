package update

import (
	"encoding/json"

	"github.com/pkg/errors"
	whPkg "github.com/zhuharev/whu/domain/webhook"

	"github.com/bloom42/rz-go/v2"
	"github.com/bloom42/rz-go/v2/log"
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

func (uc *uc) Get(id string, initialOffset int) ([]Update, error) {
	offset := initialOffset
	wh, err := uc.whRepo.Get(id)
	if err != nil {
		return nil, errors.Wrap(err, "get repo by id")
	}

	count, err := uc.repo.GetUpdatesCount(id)
	if err != nil {
		return nil, errors.Wrap(err, "get updates count")
	}

	if offset == 0 {
		offset = wh.LastOffset
	}

	if offset > count {
		log.Debug("offset greater then count", rz.Int("offset", offset), rz.Int("count", count))
		offset = count
	}

	if wh.LastOffset < offset {
		log.Debug("wh.lastoffset least then passed offset", rz.Int("lastoffset", wh.LastOffset), rz.Int("offset", offset))
		wh.LastOffset = offset
		err = uc.whRepo.UpdateLastOffset(id, offset)
		if err != nil {
			return nil, errors.Wrap(err, "update last offset")
		}
	}

	log.Debug("stat",
		rz.Int("offset", offset),
		rz.Int("last_offset", wh.LastOffset),
		rz.Int("count", count),
		rz.Int("initial_offset", initialOffset),
	)

	updates, err := uc.repo.Get(id, offset)
	if err != nil {
		return nil, err
	}

	for i := range updates {
		updates[i].ID = offset + i + 1
	}
	return updates, nil
}
