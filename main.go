package main

import (
  "log"
  "net/http"

  "filebot/internal/config"
  "filebot/internal/httpserver"
)

func main() {
  cfg, err := config.Load("")
  if err != nil {
    log.Fatalf("config load error: %v", err)
  }

  srv, err := httpserver.New(cfg)
  if err != nil {
    log.Fatalf("server init error: %v", err)
  }

  mux := http.NewServeMux()
  srv.Routes(mux)

  log.Printf("starting server on %s", cfg.Web.Address)
  if err := http.ListenAndServe(cfg.Web.Address, mux); err != nil {
    log.Fatalf("server error: %v", err)
  }
}