package mysqlimpl

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vinhphuctadang/go-interface/db"
	"github.com/vinhphuctadang/go-interface/db/model"
)

type dbServiceMySqlBackend struct {
	db.DBService
	client *sql.DB
}

func NewDbServiceMySqlBackend(connectionUri string, maxConnection int) (db.DBService, error) {
	// connectionUri: id:password@tcp(your-amazonaws-uri.com:3306)/dbname

	db, err := sql.Open("mysql", connectionUri)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxIdleConns(maxConnection)
	db.SetMaxOpenConns(maxConnection)

	return &dbServiceMySqlBackend{
		client: db,
	}, nil
}

func (d *dbServiceMySqlBackend) CreateAccount(username, password string) error {
	queryTemplate := "INSERT INTO accounts(username,password) VALUES ('%s','%s')"

	// perform a db.Query insert
	result, err := d.client.Exec(queryTemplate, username, password)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("no row inserted")
	}

	return nil
}

func (d *dbServiceMySqlBackend) ListAccount(pageIndex, pageSize int) (accounts []model.Account, err error) {
	queryTemplate := "SELECT username,password FROM accounts LIMIT %d,%d"
	rows, err := d.client.Query(queryTemplate, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var acc model.Account
		err = rows.Scan(&acc.Username, &acc.Password)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (d *dbServiceMySqlBackend) Disconnect(ctx context.Context) error {
	return d.client.Close()
}
