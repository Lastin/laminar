package shared

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/buntdb"
	"time"
)

// TagInfo is the official data stored in buntDB
type TagInfo struct {
	Image            string            `json:"image"`
	Hash             string            `json:"hash"`
	Created          time.Time         `json:"created"`
	Tag              string            `json:"tag"`
	RegistryProvider *RegistryProvider `json:"-"`
}

func (ti *TagInfo) IsNewerThan(t2 *TagInfo) bool {
	return ti.Created.Sub(t2.Created) > 0
}

func (ti *TagInfo) WriteToCache(db *buntdb.DB) (err error) {
	var byteArray []byte
	byteArray, err = json.Marshal(ti)
	if err == nil {
		err = db.Update(func(tx *buntdb.Tx) error {
			_, _, err = tx.Set(
				fmt.Sprintf("TagInfo:%s:%s:%s", ti.Image, ti.Hash, ti.Tag),
				string(byteArray),
				// TTL on tag cache, https://github.com/tidwall/buntdb#data-expiration
				&buntdb.SetOptions{Expires: true, TTL: time.Second * 300},
			)
			return err
		})
	}
	return
}

func (ti *TagInfo) ReadFromCache(db *buntdb.DB, index string) (err error) {
	err = db.View(func(tx *buntdb.Tx) error {
		err = tx.Descend(index, func(_, val string) bool {
			err = json.Unmarshal([]byte(val), ti)
			return true
		})
		return nil
	})
	return
}
