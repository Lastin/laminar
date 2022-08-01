package api

import (
	"github.com/digtux/laminar/pkg/shared"
	"github.com/tidwall/buntdb"
)

type Client struct {
	db          *buntdb.DB
	resourceMap map[string]shared.Cachable
}

func New(db *buntdb.DB) *Client {
	return &Client{
		db: db,
	}
}

func (c Client) GetDockerImages() {

}
