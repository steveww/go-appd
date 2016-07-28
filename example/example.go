package main

import (
    "log"
    "fmt"
    "net/http"
    "encoding/json"
    "strconv"

    "github.com/gorilla/mux"

    "database/sql"
    _ "github.com/lib/pq"

    appd "github.com/steveww/go-appd"
)

var backendName string = "Postgres"

func main() {
    fmt.Println("Hello World starting")
    appd.Init("TestApp", "6ca84575-ffec-470d-8099-9e527ade5033")
    appd.SetTierName("Go Lang")
    appd.SetNodeName("go1")
    appd.SetControllerHost("localhost")
    appd.SetControllerPort(8090)
    appd.SetControllerAccount("customer1")
    appd.SetControllerUseSSL(0)
    appd.SetInitTimeout(5000)
    fmt.Println("Starting SDK")
    rc := appd.Sdk_init()
    if(rc != 0) {
        log.Fatal("SDK init ", rc)
    }

    // APPD Backend
    appd.Backend_declare(appd.BACKEND_DB, backendName)
    props := appd.ID_properties_map {
        "HOST": "localhost",
        "PORT": "5432",
        "DATABASE": "postcode",
        "VENDOR": "Postgres",
        "VERSION": "9",
    }
    rc = appd.Backend_set_identifying_properties(backendName, props)
    if(rc != 0) {
        log.Fatal("Backend ", rc)
    }
    rc = appd.Backend_add(backendName)
    if(rc != 0) {
        log.Fatal("Backend ", rc)
    }
    fmt.Println("Backend added")

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", indexHandler).Methods("GET")
    router.HandleFunc("/postcode/{code}", postcodeHandler).Methods("GET")
    router.HandleFunc(appd.WrapHandleFunc("Search", "/search/{code}", searchHandler)).Methods("GET")

    fmt.Println("Ready on port 3000...")

    log.Fatal(http.ListenAndServe(":3000", router))
}

func postcodeHandler(w http.ResponseWriter, r *http.Request) {

    type Result struct {
        Code string
        Xcoord int
        Ycoord int
    }

    var query string = "select xcoord, ycoord from postpntp where postcode = $1"
    var result Result

    fmt.Println("Postcode Handler called")
    fmt.Println(r.URL)
    // dump the headers
    /* for key, value := range r.Header {
        fmt.Println("Header", key, value[len(value) - 1])
    } */
    appd_correlation := r.Header.Get(appd.CORRELATION_HEADER_NAME)
    fmt.Println("Correlation", appd_correlation)

    vars := mux.Vars(r)
    code := vars["code"]

    bt := appd.BT_begin("postcode", appd_correlation)
    defer appd.BT_end(bt)
    fmt.Println("BT ", bt)
    if(appd.BT_is_snapshotting(bt) != 0) {
        fmt.Println("adding data")
        appd.BT_set_url(bt, r.URL.String())
        appd.BT_add_user_data(bt, "postcode", code)
    }

    exit := appd.Exitcall_begin(bt, backendName)
    fmt.Println("Exit handle", exit)
    db, err := sql.Open("postgres", "postgres://stevew:demosim@localhost/postcode")
    if(err != nil) {
        log.Fatal(err)
    }
    fmt.Println("Database connected")
    defer db.Close()

    w.Header().Set("Content-Type", "application/json")

    appd.Exitcall_set_details(exit, query)
    fmt.Println("Looking for ", code)
    var (
        x int
        y int
    )
    err = db.QueryRow(query, code).Scan(&x, &y)
    switch {
        case err == sql.ErrNoRows:
            fmt.Println("No match")
            result = Result {
                Code: "No match",
                Xcoord: 0,
                Ycoord: 0,
            }
            // report error
            appd.BT_add_error(bt, appd.ERROR_LEVEL_ERROR, "no match", 1)
        case err != nil:
            log.Fatal("Query ", err)
        default:
            fmt.Println("Got row")
            result = Result {
                Code: code,
                Xcoord: x,
                Ycoord: y,
            }
    }

    appd.Exitcall_end(exit)

    data, err := json.Marshal(result)
    if(err != nil) {
        fmt.Println("JSON Error: ", err)
    }

    fmt.Fprintf(w, string(data))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Postcode REST API")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("searchHandler called")

    // get the BT id from the header
    var bt uint64 = 0
    bta := r.Header.Get(appd.APPD_BT)
    if(bta != "") {
        i, err := strconv.ParseUint(bta, 10, 64)
        if(err != nil) {
            bt = 0
        } else {
            bt = i
        }
    }
    fmt.Println("BT header", bt)

    vars := mux.Vars(r)
    code := vars["code"]

    var query string = "select count(*) from postpntp where postcode like $1"
    exit := appd.Exitcall_begin(bt, backendName)
    appd.Exitcall_set_details(exit, query)
    db, err := sql.Open("postgres", "postgres://stevew:demosim@localhost/postcode")
    if(err != nil) {
        log.Fatal(err)
    }
    defer db.Close()

    var count int
    err = db.QueryRow(query, "%" + code + "%").Scan(&count)
    switch {
        case err != nil:
            log.Fatal("Query", err)
        default:
            fmt.Println("Got count", count)
    }
    appd.Exitcall_end(exit)

    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "Search for %q found %v", code, count)
}

