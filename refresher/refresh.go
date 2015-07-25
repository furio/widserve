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
	"math"
)

var _ = log.Flags()

const (
	bufferChanSize int = 50
)


// Channels
var (
	refreshRequestChan chan string
	refreshKillChan chan int

	serviceVariationChan chan int
	serviceKillChan chan int

	currentWorkers int = 0


	tickStuff time.Ticker
	wg sync.WaitGroup
)

// Cache & Db
var (
	cacheIstance cache.CacheGeneric = nil
	dbIstance db.DataSource = nil
)

/* ===================================================================== */

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

func concurrentRefresh(work *sync.WaitGroup) {
	defer work.Done()
	doWork := true

	for doWork {
		select {
			case msg := <-refreshRequestChan:
				doRefresh(msg)
			case _ = <-refreshKillChan:
				doWork = false
		}
	}
}

func howManyToRefresh(nowTime uint64) int64 {
	count,err := dbIstance.ExpiredWidgetCount(nowTime)

	if (err != nil) {
		return 0;
	}

	return count;
}

func refreshWidgets() {
	nowTime := uint64( time.Now().Unix() )

	countWidgets := howManyToRefresh(nowTime);

	if (countWidgets <= 0) {
		return
	}

	workerRefreshQty := math.Ceil((float64)(countWidgets) / (float64)(bufferChanSize))

	for i := 0; i < (int)(workerRefreshQty); i++ {
		go func(timeExp uint64, start int, qty int) {
			widgets, err := dbIstance.GetWidgets(timeExp, start, qty)
			if err == nil && len(widgets) != 0 {
				for j := 0; j < len(widgets); j++ {
					refreshRequestChan <- widgets[i].WidgetID
				}
			}
		} (nowTime, i * bufferChanSize, bufferChanSize)
	}
}

func tickerChecker(secToRefresh int) {
	doWork := true
	secondToRefresh := secToRefresh;
	nextRefresh := secondToRefresh;

	for doWork {
		select {
			case _ <- tickStuff.C:
				{
					nextRefresh -= 1
					if (nextRefresh == 0) {
						refreshWidgets()
					}
				}

			case newRefresh := <- serviceVariationChan:
				secondToRefresh = newRefresh

			case _ <- serviceKillChan:
				doWork = false
		}
	}
}

/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */

func addCacheRefresher(count int, wg *sync.WaitGroup) {
	wg.Add(count)
	for i := 0; i < count; i++ {
		go concurrentRefresh(wg)
	}

	// This to see currently active
	currentWorkers = currentWorkers + count
}

func removeCacheRefresher(count int, wg *sync.WaitGroup) {
	for i := 0; i < count; i++ {
		refreshKillChan <- 0
	}

	// This to see currently active
	currentWorkers = currentWorkers - count
}

func initCache() {
	cacheIstance = cache.GetCacheClient(cache.Local, nil)
}

func initDb() {
	dbIstance = db.GetDataSource(db.Local, nil)
}

func initRefresh() {
	// Decide some buffer
	refreshRequestChan = make(chan string, bufferChanSize)
	refreshKillChan = make(chan int, bufferChanSize)
	serviceVariationChan = make(chan int)
	serviceKillChan = make(chan int, bufferChanSize)

	tickStuff = time.NewTicker(time.Second)


	// A go routine here that fetch the data and send to chan
	// =========
	minutes,errMinutes := dbIstance.MinimumWidgetElapseMinutes()
	if (errMinutes != nil) {
		minutes = 1
	}
	go tickerChecker((int64)(minutes * 60))


	addCacheRefresher(bufferChanSize, &wg)


	// Non blocking wait
	go func () {
		wg.Wait()

		log.Print("WaitGroup ended, remaining go routines: %d", currentWorkers)

		// Kill
		close(refreshRequestChan)
		close(refreshKillChan)
	} ()
}

func Main() {
	initDb()
	initCache()
	// initRefresh()
}