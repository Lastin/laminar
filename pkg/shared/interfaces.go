package shared

import "github.com/tidwall/buntdb"

type Cachable interface {
	WriteToCache(db *buntdb.DB) (err error)
	ReadFromCache(db *buntdb.DB, index string) (err error)
	GetIndex() string
}
