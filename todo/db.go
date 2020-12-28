package todo

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

// Item struct.
type Item struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	URL       string `json:"url"`
	Order     int    `json:"order"`
}

type database struct {
	db *buntdb.DB
}

func newDatabase(db *buntdb.DB) *database {
	return &database{db}
}

func (d *database) getAll() ([]Item, error) {
	items := make([]Item, 0)

	err := d.db.View(func(tx *buntdb.Tx) error {
		_ = tx.Ascend("order", func(k, v string) bool {
			item := Item{}
			_ = json.Unmarshal([]byte(v), &item)
			items = append(items, item)
			return true
		})
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get items")
	}

	return items, nil
}

func (d *database) save(item Item) (Item, error) {
	err := d.db.Update(func(tx *buntdb.Tx) error {
		k := item.Key
		itemBytes, _ := json.Marshal(item)
		v := string(itemBytes)
		_, _, _ = tx.Set(k, v, nil)
		return nil
	})
	if err != nil {
		return Item{}, errors.Wrap(err, "failed to save item")
	}

	return item, nil
}

func (d *database) delete(key string) error {
	err := d.db.Update(func(tx *buntdb.Tx) error {
		delkeys := make([]string, 0)

		_ = tx.Ascend("", func(k, v string) bool {
			if key == k {
				delkeys = append(delkeys, k)
				return false
			}
			if key == "" {
				delkeys = append(delkeys, k)
			}
			return true
		})

		for _, k := range delkeys {
			if _, err := tx.Delete(k); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete item(s)")
	}

	return nil
}

func (d *database) getOne(key string) (Item, error) {
	var item Item

	err := d.db.View(func(tx *buntdb.Tx) error {
		v, _ := tx.Get(key)
		val := Item{}
		_ = json.Unmarshal([]byte(v), &val)
		item = val
		return nil
	})
	if err != nil {
		return Item{}, errors.Wrap(err, "failed to get item")
	}

	return item, nil
}
