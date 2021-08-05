package main

import (
	"bwlee_proxy/reverse"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type serverinfo struct {
	rparray *reverse.RpServerArray
}

func (rvinfo *serverinfo) ServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		log.Println("###ServerInfoHandler 'POST' start###")
		var serverinfo *reverse.ServerInfo

		err := json.NewDecoder(r.Body).Decode(&serverinfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Host : %s \n", serverinfo.Host)
		log.Printf("Origin : %s \n", serverinfo.Origin)
		log.Printf("Port : %d \n", serverinfo.Port)
		rvinfo.rparray.RvProxyHandle(serverinfo)

	case "GET":
		io.WriteString(w, string(r.Method)+"\n")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func main() {
	log.Println("###bwlee_proxy main start###")
	var rvinfo serverinfo

	rvinfo.rparray = reverse.RvInit()
	log.Println("###reverse.RvInit() end###")
	mux := http.NewServeMux()
	mux.HandleFunc("/", rvinfo.ServerInfoHandler)
	log.Fatal(http.ListenAndServe(":301", mux))
	log.Printf("###bwlee_api from 301 port start###\n")
}
