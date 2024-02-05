package postgresql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewClient(dns string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dns)
	if err != nil {
		return nil, err
	}

	return db, nil
}
