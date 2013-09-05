package main

import (
    "fmt"
    "net/http"
    "os"
    "encoding/json"
    "github.com/gorilla/mux"
    "io/ioutil"
    "strings"
    "math/rand"
)

func main() {
    var router = mux.NewRouter()

    router.HandleFunc("/quote/{author}", QuoteHandler)

    http.Handle("/", router)

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
      panic(err)
    }
}

func QuoteHandler(res http.ResponseWriter, req *http.Request) {
  res.Header().Set("Content-type", "application/json; charset=utf-8")
  // see if file exists
  // read first line for full author name
  // make json
  vars := mux.Vars(req)
  author := vars["author"]
  authorFile := "./" + author + ".quote"
  if _, err := os.Stat(authorFile); os.IsNotExist(err) {
    // return 404
  } else {
    fileContents, _ := ioutil.ReadFile(authorFile)
    splitFileContents := strings.Split(string(fileContents), "\n")
    authorName := splitFileContents[0]
    splitFileContents = splitFileContents[1:]
    quoteIndex := rand.Int31n(int32(len(splitFileContents) - 1))
    quote := splitFileContents[quoteIndex]

    m := make(map[string]string)
    m["author"] = authorName
    m["quote"] = quote
    value, _ := json.Marshal(m)
    fmt.Fprintln(res, string(value))
  }

}
