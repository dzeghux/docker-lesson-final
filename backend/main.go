package main

import (
 "database/sql"
 "encoding/json"
 "log"
 "net/http"
 "os"
 "strings"

 _ "github.com/lib/pq"
)

type Item struct {
 ID   int    `json:"id"`
 Name string `json:"name"`
}

func main() {
 dsn := "host=" + env("DB_HOST", "db") +
  " port=" + env("DB_PORT", "5432") +
  " user=" + env("DB_USER", "app") +
  " password=" + env("DB_PASSWORD", "secret") +
  " dbname=" + env("DB_NAME", "appdb") +
  " sslmode=disable"

 db, err := sql.Open("postgres", dsn)
 if err != nil {
  log.Fatal(err)
 }
 defer db.Close()

 staticDir := env("STATIC_DIR", "/app/static")

 http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
  rows, err := db.Query("SELECT id, name FROM items ORDER BY id")
  if err != nil {
   http.Error(w, err.Error(), 500)
   return
  }
  defer rows.Close()

  items := []Item{}
  for rows.Next() {
   var it Item
   if err := rows.Scan(&it.ID, &it.Name); err != nil {
    http.Error(w, err.Error(), 500)
    return
   }
   items = append(items, it)
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(items)
 })

 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  if path == "/" {
   path = "/index.html"
  }
  
  filePath := staticDir + path
  data, err := os.ReadFile(filePath)
  if err != nil {
   http.NotFound(w, r)
   return
  }
  
  if strings.HasSuffix(path, ".js") {
   w.Header().Set("Content-Type", "application/javascript")
  } else if strings.HasSuffix(path, ".css") {
   w.Header().Set("Content-Type", "text/css")
  } else if strings.HasSuffix(path, ".html") {
   w.Header().Set("Content-Type", "text/html")
  }
  w.Write(data)
 })

 addr := ":" + env("HTTP_PORT", "8080")
 log.Printf("listening on %s", addr)
 if err := http.ListenAndServe(addr, nil); err != nil {
  log.Fatal(err)
 }
}

func env(key, def string) string {
 if v := os.Getenv(key); v != "" {
  return v
 }
 return def
}