package repo

import (
	"github.com/asdine/storm"
	"github.com/zhuharev/whu/domain/errors"
	"github.com/zhuharev/whu/domain/update"
)

type DB struct {
	db *storm.DB
}

func New(path string) (*DB, error) {
	sdb, err := storm.Open(path)
	//defer sdb.Close()

	if err != nil {
		return nil, errors.ErrCannotOpenDB
	}
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
	err = db.db.Find("WH", wh, &updates, storm.Skip(offset))
	if err == storm.ErrNotFound {
		return nil, nil
	}
	return
}
