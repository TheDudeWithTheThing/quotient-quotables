package main

import (
    "encoding/json"
    "fmt"
    "github.com/garyburd/redigo/redis"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strings"
)

var server, password string = ParseRedistogoUrl()

var pool = &redis.Pool {
  MaxIdle: 3,
  IdleTimeout: 240,
  MaxActive: 10,
  Dial: func () (redis.Conn, error) {
    c, err := redis.Dial("tcp", server)
    if err != nil {
      return nil, err
    }
    if _, err := c.Do("AUTH", password); err != nil && len(password) > 0 {
      c.Close()
      return nil, err
    }
    return c, err
  },
}

func main() {
    var router = mux.NewRouter()

    router.HandleFunc("/quote/{author}", QuoteHandler)
    router.HandleFunc("/load", LoadHandler)

    http.Handle("/", router)

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      fmt.Println("panicing...")
      panic(err)
    }
}


func LoadHandler(res http.ResponseWriter, req *http.Request) {
  var client = pool.Get()
  defer client.Close()
  quoteFiles, _ := filepath.Glob("*.quote")

  for _, file := range quoteFiles {
    fileContents, _ := ioutil.ReadFile(file)
    splitFileContents := strings.Split(string(fileContents), "\n")
    authorName := splitFileContents[0]
    splitFileContents = splitFileContents[1:]

    uri := strings.Split(file, ".")
    redisKey := "quote:" + uri[0]

    client.Do("SET", "realname:" + uri[0], []byte(authorName))

    for _, quote := range splitFileContents {
      if len(quote) > 0  {
        client.Do("SADD", redisKey, []byte(quote))
      }
    }
  }

  fmt.Fprintln(res, "loading complete")
}

func QuoteHandler(res http.ResponseWriter, req *http.Request) {
  var client = pool.Get()
  defer client.Close()
  vars := mux.Vars(req)
  author := vars["author"]
  authorName, err := redis.String(client.Do("GET", "realname:" + author))
  if err != nil {
    http.NotFound(res, req)
  } else {
    quote, _ := redis.String(client.Do("SRANDMEMBER", "quote:" + author))

    m := make(map[string]string)
    m["author"] = string(authorName)
    m["quote"] = string(quote)
    value, _ := json.Marshal(m)

    res.Header().Set("Content-type", "application/json; charset=utf-8")
    fmt.Fprintln(res, string(value))
  }
}

func ParseRedistogoUrl() (string, string) {
  redisUrl := os.Getenv("REDISTOGO_URL")
  fmt.Println(redisUrl)

  redisInfo, _ := url.Parse(redisUrl)
  server := redisInfo.Host
  fmt.Println("SERVER " + server)

  password := ""
  if redisInfo.User != nil {
    password, _ = redisInfo.User.Password()
  }
  return server, password
}
