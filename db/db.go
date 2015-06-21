package db

import(

	_ "github.com/furio/widserve/db/uid"

	_ "gopkg.in/gorp.v1"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)