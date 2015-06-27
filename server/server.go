package server

import (
    "net/http"
    "fmt"
    "time"
    "encoding/json"
    _ "log"

    // My stuff
    "github.com/furio/widserve/cache"

    // Routing
    "github.com/gorilla/mux"

    // Middleware
    "github.com/codegangsta/negroni"
    "github.com/phyber/negroni-gzip/gzip"
    "github.com/rs/cors"
    "gopkg.in/tylerb/graceful.v1"
    "github.com/thoas/stats"

    // https://github.com/mholt/binding
    // https://github.com/unrolled/secure
)

var corsConfig = cors.New(cors.Options{
    AllowedMethods: []string{"GET","POST","OPTIONS"},
})

var statsMiddle = stats.New()

var cacheIstance cache.CacheGeneric = nil


func newWidget(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome admin!")
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
        // Chack a validity
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

func Main() {
    initCache()
    initServer()
}

func initCache() {
    cacheIstance = cache.GetCacheClient(cache.Local, nil)
}

func initServer() {
    router := mux.NewRouter()

    // Admin stuff
    adminRoutes := router.PathPrefix("/admin").Subrouter()
    adminRoutes.HandleFunc("/stats", adminStats).Methods("GET")

    adminRoutes.HandleFunc("/widgets", newWidget).Methods("GET,POST")
    adminRoutes.HandleFunc("/widget/{key}", newWidget).Methods("GET,DELETE")
    adminRoutes.HandleFunc("/widget/{key}/force", newWidget).Methods("POST")


    // Client stuff
    clientRoutes := router.PathPrefix("/widgets").Subrouter()
    clientRoutes.HandleFunc("/{wkey}", getWidget).Methods("GET")

    // Make the server
    n := negroni.New()
    n.Use(corsConfig)
    n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.Use(statsMiddle)
    n.UseHandler(router)

    graceful.Run(":3000", 10*time.Second, n)
}
