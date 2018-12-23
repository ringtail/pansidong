package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/ringtail/pansidong/cache"
	"github.com/ringtail/pansidong/types"
	"log"
	"fmt"
)

type Server struct {
	Addr         string
	CacheManager *cache.CacheManager
}

func (s *Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/ips", s.NextIpsHandler).Methods(http.MethodGet)
	r.HandleFunc("/ip/{ip}/expire", s.ExpireHandler).Methods(http.MethodPost)
	srv := &http.Server{
		Handler: r,
		Addr:    s.Addr,
	}
	stopChan := make(chan struct{})
	go s.CacheManager.Loop(stopChan)
	if err := (srv.ListenAndServe()); err != nil {
		stopChan <- struct{}{}
		log.Fatal(err)
	}
}

func NewServer(config *types.Config) *Server {
	if err := config.Valid(); err != nil {
		log.Fatal("Please input valid config file: %s", err.Error())
	}
	cc := config.Config()
	s := &Server{
		Addr: fmt.Sprintf("%s:%d", config.GlobalConfig.Host, config.GlobalConfig.Port),
	}
	c, err := cache.NewCacheManager(cc)
	if err != nil {
		log.Fatal("Failed to create cache manager because of %s", err.Error())
	}
	s.CacheManager = c
	return s
}
