package repo

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/zhuharev/whu/domain/update"
)

type DB struct {
	db *storm.DB
}

func New(sdb *storm.DB) (*DB, error) {
	return &DB{db: sdb}, nil
}

func (db *DB) Save(wh string, payload []byte) error {
	upd := update.Update{
		WH:      wh,
		Payload: payload,
	}
	err := db.db.Save(&upd)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Get(wh string, offset int) (updates []update.Update, err error) {
	err = db.db.Find("WH", wh, &updates, storm.Skip(offset), storm.Limit(10))
	if err == storm.ErrNotFound {
		return nil, nil
	}
	return
}

func (db *DB) GetUpdatesCount(id string) (count int, err error) {
	return db.db.Select(q.Eq("WH", id)).Count(&update.Update{})
}
