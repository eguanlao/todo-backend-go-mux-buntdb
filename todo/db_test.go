package todo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/buntdb"
)

func Test_database_getAll(t *testing.T) {
	type fields struct {
		db *buntdb.DB
	}

	db, _ := buntdb.Open(":memory:")
	_ = db.CreateIndex("order", "*", buntdb.IndexJSON("order"))
	defer db.Close()

	stubDBWithItems := func(db *buntdb.DB) *buntdb.DB {
		_ = db.Update(func(tx *buntdb.Tx) error {
			_, _, _ = tx.Set("1", `{"title":"dummy title 1","completed":false,"order":0}`, nil)
			_, _, _ = tx.Set("2", `{"title":"dummy title 2","completed":true, "order":2}`, nil)
			return nil
		})
		return db
	}

	tests := []struct {
		name    string
		fields  fields
		want    []Item
		wantErr bool
	}{
		{
			name:   "ShouldReturnItems",
			fields: fields{stubDBWithItems(db)},
			want: []Item{
				{Title: "dummy title 1"},
				{Title: "dummy title 2", Completed: true, Order: 2},
			},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			d := database{
				db: test.fields.db,
			}

			got, err := d.getAll()

			if err != nil {
				assert.True(t, test.wantErr)
				assert.Nil(t, got)
			} else {
				assert.False(t, test.wantErr)
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func Test_database_save_getOne(t *testing.T) {
	type fields struct {
		db *buntdb.DB
	}

	type args struct {
		item Item
	}

	db, _ := buntdb.Open(":memory:")
	_ = db.CreateIndex("order", "*", buntdb.IndexJSON("order"))
	defer db.Close()

	stubDB := func(db *buntdb.DB) *buntdb.DB {
		return db
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Item
		wantErr bool
		saved   func() bool
	}{
		{
			name:   "ShouldSaveItem",
			fields: fields{stubDB(db)},
			args:   args{item: Item{Title: "dummy title 1"}},
			want:   Item{Title: "dummy title 1"},
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			d := database{
				db: test.fields.db,
			}

			got, err := d.save(test.args.item)

			if err != nil {
				assert.True(t, test.wantErr)
				assert.Nil(t, got)
			} else {
				assert.False(t, test.wantErr)
				assert.Equal(t, test.want, got)
				saved, _ := d.getOne(got.Key)
				assert.NotNil(t, saved)
			}
		})
	}
}

func Test_database_delete(t *testing.T) {
	type fields struct {
		db *buntdb.DB
	}

	type args struct {
		key string
	}

	db, _ := buntdb.Open(":memory:")
	_ = db.CreateIndex("order", "*", buntdb.IndexJSON("order"))
	defer db.Close()

	stubDB := func(db *buntdb.DB) *buntdb.DB {
		return db
	}

	populateDB := func(d database) {
		_, _ = d.save(Item{Key: "1", Title: "dummy title 1"})
		_, _ = d.save(Item{Key: "2", Title: "dummy title 2", Completed: true, Order: 2})
		_, _ = d.save(Item{Key: "3", Title: "dummy title 3", Completed: true})
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantLen int
	}{
		{
			name:    "ShouldDeleteItem",
			fields:  fields{stubDB(db)},
			args:    args{"2"},
			wantLen: 2,
		},
		{
			name:    "ShouldDeleteAllItems",
			fields:  fields{stubDB(db)},
			args:    args{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			d := database{
				db: test.fields.db,
			}

			populateDB(d)

			err := d.delete(test.args.key)

			if err != nil {
				assert.True(t, test.wantErr)
			} else {
				assert.False(t, test.wantErr)
				items, _ := d.getAll()
				assert.Equal(t, test.wantLen, len(items))
			}
		})
	}
}
