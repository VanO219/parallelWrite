package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"net/url"
)

type DB struct {
	DB1 *sql.DB
}

func (d *DB) Connect(user, password string) (err error)  {
	defer func() {err = errors.Wrap(err, "main.Connect")}()

	u := url.URL{
		Scheme: "postgres",
		Host: "127.0.0.1:5000",
		User: url.UserPassword(user, password),
		Path: "testDB",
		RawQuery: "sslmode=disable",
	}
	d.DB1, err = sql.Open("postgres", u.String())
	if err != nil {
		err = errors.Wrap(err, "failed open connection")
		return err
	}
	return err
}

func (d *DB) Close() (err error)  {
	defer func() {err = errors.Wrap(err, "main.Close")}()

	if err = d.DB1.Close(); err != nil {
		err = errors.Wrap(err, "failed close connection")
		return err
	}
	return err
}
