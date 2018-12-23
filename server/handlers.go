package server

import (
	"net/http"
	"github.com/ringtail/pansidong/types"
	"encoding/json"
	"strconv"
	"github.com/gorilla/mux"
)

func (s Server) NextIpsHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	limit := vals.Get("limit")
	l := types.DefaultPoolSize
	if limit != "" {
		sl, err := strconv.Atoi(limit)
		if err == nil {
			l = sl
		}
	}

	options := &types.ListOptions{
		Limit: l,
	}
	ips, err := s.CacheManager.FetchNextIps(options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	b, err := json.Marshal(ips)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s Server) ExpireHandler(w http.ResponseWriter, r *http.Request) {
	//ips := make([]string, 0)
	//b, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	w.WriteHeader(http.StatusPaymentRequired)
	//	w.Write([]byte("can not expire those ips,because of " + err.Error()))
	//	return
	//}
	//err = json.Unmarshal(b, ips)
	//if err != nil {
	//	w.WriteHeader(http.StatusPaymentRequired)
	//	w.Write([]byte("can not expire those ips,because of " + err.Error()))
	//	return
	//}

	vars := mux.Vars(r)
	if vars["ip"] == "" {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}
	err := s.CacheManager.ExpireIp(vars["ip"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
