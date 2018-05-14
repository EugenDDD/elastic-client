package cmd

import (
	"elastic-search/global"
	"sync"

	"github.com/olivere/elastic"
)

var (
	config       *global.Config
	configParams global.Config
	lock         sync.Mutex
	changePass   bool

	result   global.Result
	client   *elastic.Client
	getFlags global.LogFlags
	queries  []*elastic.Query
	results  []global.Result
)
