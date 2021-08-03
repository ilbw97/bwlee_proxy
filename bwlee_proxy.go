package main

import (
	"bwlee_proxy/reverse"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (efgd *abcd) S_info_Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var info *reverse.Information
		get_info, err := json.Marshal(info)
		if err != nil {
			log.Println(err)
		}
		io.WriteString(w, string(get_info))
		w.Write([]byte("\n"))
		// jsonfile, err := os.Create("/root/go/src/bwlee_proxy/information.json")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer jsonfile.Close()

		// jsonfile.Write(get_info)
		// fmt.Println("JSON DATA WRRITEN TO \n", jsonfile.Name())
		// reverse.Rvproxy_handle()

	case "POST":
		var info *reverse.Information

		err := json.NewDecoder(r.Body).Decode(&info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("sending %v \n", info)
		log.Printf("Host : %v, Origin : %v, Port : %v \n", info.Host, info.Origin, info.Port)

		// reverse.Rvproxy_handle(info)
		efgd.reverseInfo.Rvproxy_handle(info)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "invalid type of method")
	}
}

type abcd struct {
	reverseInfo *reverse.Rpserver_array
}

func main() {

	var efgd abcd

	efgd.reverseInfo = reverse.Reverse_init()

	mux := http.NewServeMux()
	mux.HandleFunc("/", efgd.S_info_Handler)
	log.Fatal(http.ListenAndServe(":301", mux))
}
