package server

import (
    "net/http"
    "fmt"
    "time"
    "encoding/json"
    "log"

    // My stuff
    "github.com/furio/widserve/db"
    "github.com/furio/widserve/cache"

    // Routing
    "github.com/gorilla/mux"

    // Middleware
    "github.com/codegangsta/negroni"
    "github.com/phyber/negroni-gzip/gzip"
    "github.com/rs/cors"
    "gopkg.in/tylerb/graceful.v1"
    "github.com/thoas/stats"
	"github.com/mholt/binding"

    // https://github.com/unrolled/secure
	// https://github.com/martini-contrib/throttle
)

var _ = log.Flags()

// Midleware
var corsConfig = cors.New(cors.Options{
    AllowedMethods: []string{"GET","POST","DELETE","OPTIONS"},
})
var statsMiddle = stats.New()


// Cache & Db
var cacheIstance cache.CacheGeneric = nil
var dbIstance db.DataSource = nil

/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */

type WidgetForm struct {
	WidgetID	string
	ApiKey		string
	ApiPath		string
	CacheElapse	uint32
}

func (cf *WidgetForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&cf.WidgetID: "WidgetID",
		&cf.ApiKey:  binding.Field{
			Form:     "ApiKey",
			Required: true,
		},
		&cf.ApiPath:  binding.Field{
			Form:     "ApiPath",
			Required: true,
		},
		&cf.CacheElapse: binding.Field{
			Form:     "CacheElapse",
			Required: true,
		},
	}
}

func adminStats(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    stats := statsMiddle.Data()
    b, _ := json.Marshal(stats)

    w.Write(b)
}

func getWidget(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get vars
	vars := mux.Vars(req)
	val, key := vars["wkey"]

	if (key) {
		widget, err := dbIstance.GetWidget(val)

		if (err == nil) {
			b, _ := json.Marshal(widget)

			w.Write(b)
		} else {
			http.NotFound(w, req);
		}
	} else {
		http.NotFound(w, req);
	}
}

func createWidget(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	wForm := new(WidgetForm)
	errs := binding.Bind(req, wForm)

	if errs.Handle(w) {
		return
	}

	wData, err := dbIstance.NewWidget(wForm.ApiKey, wForm.ApiPath, wForm.CacheElapse)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	} else {
		b, _ := json.Marshal(wData)
		w.Write(b)
	}
}

func deleteWidget(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get vars
	vars := mux.Vars(req)
	val, key := vars["wkey"]

	if (key) {
		res, _ := dbIstance.DeleteWidgetByKey(val)
		b, _ := json.Marshal(map[string]bool {"result": res})

		w.Write(b)
	} else {
		http.NotFound(w, req);
	}
}

func getCachedWidget(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // get vars
    vars := mux.Vars(req)
    val, key := vars["wkey"]

    if (key) {
        data, found := cacheIstance.Get(val)

        if (found) {
            fmt.Fprintf(w, data.(string))
        } else {
            http.NotFound(w, req);
        }
    } else {
        http.NotFound(w, req);
    }
}

/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */
/* ===================================================================== */

func initCache() {
    cacheIstance = cache.GetCacheClient(cache.Local, nil)
}

func initDb() {
    dbIstance = db.GetDataSource(db.Local, nil)
}

func initServer() {
    router := mux.NewRouter()

    // Admin stuff
    adminRoutes := router.PathPrefix("/admin").Subrouter()
    adminRoutes.HandleFunc("/stats", adminStats).Methods("GET")

//  adminRoutes.HandleFunc("/widgets", listWidgets).Methods("GET")
	adminRoutes.HandleFunc("/widgets", createWidget).Methods("POST")
    adminRoutes.HandleFunc("/widgets/{wkey}", getWidget).Methods("GET")
	adminRoutes.HandleFunc("/widgets/{wkey}", deleteWidget).Methods("DELETE")
//  adminRoutes.HandleFunc("/widgets/{wkey}/force", newWidget).Methods("POST")


    // Client stuff
    clientRoutes := router.PathPrefix("/widgets").Subrouter()
    clientRoutes.HandleFunc("/{wkey}", getCachedWidget).Methods("GET")

    // Make the server
    n := negroni.New()
    n.Use(corsConfig)
    n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.Use(statsMiddle)
    n.UseHandler(router)

    graceful.Run(":3000", 10*time.Second, n)
}

func Main() {
	initDb()
	initCache()
	initServer()
}