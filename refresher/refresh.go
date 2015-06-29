package refresher


import (
	_ "net/http"
	_ "fmt"
	"sync"
	"time"
	_ "encoding/json"
	"log"

	// My stuff
	"github.com/furio/widserve/db"
	"github.com/furio/widserve/cache"

	//
	"github.com/parnurzeal/gorequest"
)

var _ = log.Flags()


// Channels
var requestChan chan string
var killChan chan int

// Cache & Db
var cacheIstance cache.CacheGeneric = nil
var dbIstance db.DataSource = nil

func doRefresh(wId string) {
	p1, err := dbIstance.GetWidget(wId)

	if (err != nil) {
		return
	}

	request := gorequest.New().Timeout(30 * time.Second)
	resp, body, errs := request.Get(p1.ApiPath).
		Set(p1.ApiHeader, p1.ApiKey).
		End()

	// Update db
	dbIstance.UpdateNextCheckWidget(p1)

	if len(errs) != 0 {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	cacheIstance.Set(wId, body, time.Duration(p1.CacheElapse) * time.Second)
}

func concurrentRefresh(id int, work *sync.WaitGroup) {
	defer work.Done()
	doWork := true

	for doWork {
		select {
			case msg := <-requestChan:
				doRefresh(msg)
			case _ = <-killChan:
				doWork = false
		}
	}
}

/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */

func addCacheRefresher(id int, wg *sync.WaitGroup) {
	wg.Add(id)
	for i := 0; i < id; i++ {
		go concurrentRefresh(i, wg)
	}
}


func initCache() {
	cacheIstance = cache.GetCacheClient(cache.Local, nil)
}

func initDb() {
	dbIstance = db.GetDataSource(db.Local, nil)
}

func initRefresh() {
	var wg sync.WaitGroup

	// Decide some buffer
	requestChan = make(chan string, 50)
	killChan = make(chan int, 50)

	// A go routine here that fetch the data and send to chan
	// =========


	maximumWorker := 50
	addCacheRefresher(maximumWorker, &wg)

	// A go routine that increase the number if necessary here
	// =========

	wg.Wait()

	// Kill
	close(requestChan)
	close(killChan)
}

func Main() {
	initDb()
	initCache()
	// initRefresh()
}