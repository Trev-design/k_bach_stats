package sqldb

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DBAdapter struct {
	*sql.DB
	waitgroup *sync.WaitGroup
	isclosed  bool
}

func NewDatabaseAdapter(setupOpt string) (*DBAdapter, error) {
	db, err := computeSetup(setupOpt)
	if err != nil {
		return nil, err
	}

	return &DBAdapter{
		DB:        db,
		waitgroup: &sync.WaitGroup{},
		isclosed:  false,
	}, nil
}

func (dba *DBAdapter) CloseDB() error {
	return dba.Close()
}

func (dba *DBAdapter) Wait() {
	dba.waitgroup.Wait()
}

func (dba *DBAdapter) Closed() bool {
	return dba.isclosed
}

func (dba *DBAdapter) Cancel() {
	dba.isclosed = true
}
