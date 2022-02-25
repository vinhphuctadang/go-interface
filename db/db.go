package db

import (
	"context"

	"github.com/vinhphuctadang/go-interface/db/model"
)

type DBService interface {
	CreateAccount(username, password string) error
	ListAccount(pageIndex, pageSize int) ([]model.Account, error)
	Disconnect(context.Context) error
}
