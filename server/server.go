package server

import (
    "net/http"
    "fmt"
    "time"
    "encoding/json"
    "log"

    // My stuff
    "github.com/furio/widserve/cache"

    // "gopkg.in/redis.v3"
    "github.com/gorilla/mux"

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

/*
var redisCli = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
})
*/

func NewWidget(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Welcome admin!")
}

func AdminStats(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    stats := statsMiddle.Data()
    b, _ := json.Marshal(stats)

    w.Write(b)
}

func GetWidget(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    val, key := vars["wkey"]

    if (key) {
        // Chack a validity
        // grab from redis

        log.Println(val)
        fmt.Fprintf(w, "Welcome widget!")
    } else {
        http.NotFound(w, req);
    }
}

func Main() {
    // cacheIstance = cache.GetCacheClient(cache.Local, nil)


    router := mux.NewRouter()

    // Admin stuff
    adminRoutes := router.PathPrefix("/admin").Subrouter()
    adminRoutes.HandleFunc("/stats", AdminStats).Methods("GET")

    adminRoutes.HandleFunc("/widget/new", NewWidget).Methods("POST")
    adminRoutes.HandleFunc("/widget/{key}", NewWidget).Methods("GET,DELETE")
    adminRoutes.HandleFunc("/widget/{key}/refresh", NewWidget).Methods("GET")


    // Client stuff
    clientRoutes := router.PathPrefix("/widgets").Subrouter()
    clientRoutes.HandleFunc("/{wkey}", GetWidget).Methods("GET")

    // Make the server
    n := negroni.New()
    n.Use(corsConfig)
    n.Use(gzip.Gzip(gzip.DefaultCompression))
    n.Use(statsMiddle)
    n.UseHandler(router)

    graceful.Run(":3000", 10*time.Second, n)
}
