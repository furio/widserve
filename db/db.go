package db

import(
    "log"
    "os"
    _ "sync"
    "database/sql"
	_ "github.com/furio/widserve/db/uid"

	"gopkg.in/gorp.v1"
	"github.com/mattn/go-sqlite3"
	"github.com/go-sql-driver/mysql"
)

// Force godeps
var _ = mysql.ErrBusyBuffer
var _ = sqlite3.SQLiteConn{}


type DataSource struct {
    orm *gorp.DbMap
}

type DbType int
const (
    Local DbType = iota
    MySQL
)

const tableName string = "widgets"

func GetDataSource(dbType DbType, config map[string]string) DataSource {
    outDb := DataSource{}

    if (dbType == Local) {
        db, err := sql.Open("sqlite3", "/tmp/post_db.bin") // config["dbSource"]
        if (err != nil) {
            return nil
        }


        outDb.orm = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
    } else if (dbType == MySQL) {
        db, err := sql.Open("mysql", "user:password@/dbname") // config["dbSource"]
        if (err != nil) {
            return nil
        }

        outDb.orm = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
    }

    // Bind table
    outDb.orm.AddTableWithName(Widget{}, tableName).SetKeys(false, "WidgetId")

    // create DB
    createDb(outDb, dbType)

    // If from config
    outDb.orm.TraceOn("[gorp]", log.New(os.Stdout, "db:", log.Lmicroseconds))

    return outDb
}

func createDb(dbSource DataSource, dbType DbType) bool {
    err := dbSource.orm.CreateTablesIfNotExists()

    if (err == nil) {
        if (dbType == Local) {
            dbSource.orm.Exec("CREATE INDEX IF NOT EXIST nextcheckindex ON " + tableName + "(next_cache_check)")
            dbSource.orm.Exec("CREATE INDEX IF NOT EXIST cachedurationindex ON " + tableName + "(cache_elapse)")
            dbSource.orm.Exec("CREATE INDEX IF NOT EXIST apikeyindex ON " + tableName + "(api_key)")
        } else if (dbType == MySQL) {
            mySqlIndex := "SELECT COUNT(1) IndexIsThere FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema=DATABASE() AND table_name=? AND index_name=?"

            i64, err := dbSource.orm.SelectInt(mySqlIndex, tableName, "nextcheckindex")
            if (err == nil && i64 == 1) {
                dbSource.orm.Exec("CREATE INDEX nextcheckindex ON " + tableName + "(next_cache_check)")
            }

            i64, err = dbSource.orm.SelectInt(mySqlIndex, tableName, "cachedurationindex")
            if (err == nil && i64 == 1) {
                dbSource.orm.Exec("CREATE INDEX cachedurationindex ON " + tableName + "(cache_elapse)")
            }

            i64, err = dbSource.orm.SelectInt(mySqlIndex, tableName, "apikeyindex")
            if (err == nil && i64 == 1) {
                dbSource.orm.Exec("CREATE INDEX apikeyindex ON " + tableName + "(api_key)")
            }
        }
    }

    return err == nil
}

func (this DataSource) GetWidget(id string) Widget {
    p1, _ := this.orm.Get(Widget{}, id)

    return p1
}