package db

import(
    "log"
    "os"
    _ "sync"
    "database/sql"
	"github.com/furio/widserve/db/uid"

	"gopkg.in/gorp.v1"
	"github.com/mattn/go-sqlite3"
	"github.com/go-sql-driver/mysql"
)

// Force godeps
var _ = mysql.ErrBusyBuffer
var _ = sqlite3.SQLiteConn{}

type DataSource interface {
    GetWidget(id string) (Widget,error)
	NewWidget(apiKey string, apiPath string, cacheElapse uint32) (Widget,error)
	DeleteWidget(wObj Widget) (bool,error)
}

type DatabaseSource struct {
    orm *gorp.DbMap
}

type DbType int
const (
    Local DbType = iota
    MySQL
)

const tableName string = "widgets"

func GetDataSource(dbType DbType, config map[string]string) DataSource {
    outDb := DatabaseSource{}

    if (dbType == Local) {
        db, err := sql.Open("sqlite3", "/home/furio/tmp/widget_db.bin") // config["dbSource"]
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
    mapTable(outDb)

    // create DB
    createDb(outDb, dbType)

    // If from config
    outDb.orm.TraceOn("[gorp]", log.New(os.Stdout, "db:", log.Lmicroseconds))

    return outDb
}

func (this DatabaseSource) GetWidget(id string) (Widget,error) {
    p1, err := this.orm.Get(Widget{}, id)

    return *(p1.(*Widget)), err
}

func (this DatabaseSource) NewWidget(apiKey string, apiPath string, cacheElapse uint32) (Widget,error) {
	p1 := newWidget(uid.NewUid(apiKey + apiPath), apiKey, apiPath, cacheElapse)
	// _ = "breakpoint"
	err := this.orm.Insert(&p1)
	// _ = "breakpoint"

	return p1, err
}

func (this DatabaseSource) DeleteWidget(wObj Widget) (bool,error) {
	p1,err := this.orm.Delete(wObj)

	return p1==1, err
}

func mapTable(dbSource DatabaseSource) {
    // Had issue with inline mapping

    table := dbSource.orm.AddTableWithName(Widget{}, tableName)
    table.SetKeys(false, "WidgetID")
    table.ColMap("WidgetID").SetNotNull(true).SetMaxSize(255)
    table.ColMap("ApiKey").SetNotNull(true).SetMaxSize(255)
    table.ColMap("ApiPath").SetNotNull(true).SetMaxSize(1024)
    table.ColMap("Created")
    table.ColMap("CacheElapse")
    table.ColMap("NextCheck")
}

func createDb(dbSource DatabaseSource, dbType DbType) bool {
    err := dbSource.orm.CreateTablesIfNotExists()

    if (err == nil) {
        if (dbType == Local) {
            dbSource.orm.Db.Exec("CREATE INDEX IF NOT EXISTS nextcheckindex ON " + tableName + "(NextCheck)")
            dbSource.orm.Db.Exec("CREATE INDEX IF NOT EXISTS cachedurationindex ON " + tableName + "(CacheElapse)")
            dbSource.orm.Db.Exec("CREATE INDEX IF NOT EXISTS apikeyindex ON " + tableName + "(ApiKey)")
        } else if (dbType == MySQL) {
            mySqlIndex := "SELECT COUNT(1) IndexIsThere FROM INFORMATION_SCHEMA.STATISTICS WHERE table_schema=DATABASE() AND table_name=? AND index_name=?"

            i64, err := dbSource.orm.SelectInt(mySqlIndex, tableName, "nextcheckindex")
            if (err == nil && i64 == 1) {
                dbSource.orm.Db.Exec("CREATE INDEX nextcheckindex ON " + tableName + "(NextCheck)")
            }

            i64, err = dbSource.orm.SelectInt(mySqlIndex, tableName, "cachedurationindex")
            if (err == nil && i64 == 1) {
                dbSource.orm.Db.Exec("CREATE INDEX cachedurationindex ON " + tableName + "(CacheElapse)")
            }

            i64, err = dbSource.orm.SelectInt(mySqlIndex, tableName, "apikeyindex")
            if (err == nil && i64 == 1) {
                dbSource.orm.Db.Exec("CREATE INDEX apikeyindex ON " + tableName + "(ApiKey)")
            }
        }
    }

    return err == nil
}