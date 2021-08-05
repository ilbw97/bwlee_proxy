package reverse

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type RPSever struct {
	Proxy *httputil.ReverseProxy
	sinfo ServerInfo
	// ServerInfo
	// Host string
}

type RpServerArray struct {
	// Proxy []RPSever
	server []RPSever

	Port [65535]bool //range of port number
}

// type RpServerArray []struct {
// 	Proxy *httputil.ReverseProxy
// 	Host  string
// 	Port  [65535]bool
// }
type ServerInfo struct {
	Host   string `json:"host"`
	Origin string `json:"origin"`
	Port   int    `json:"port"`
}

func (rpserver *RpServerArray) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("###rpserver.ServeHTTP start###")
	// if len(rpserver.Proxy) == 0 {
	// 	log.Printf("%d", len(rpserver.Proxy))
	// }
	// if rpserver.Proxy == r.Host {
	// for _, rp := range rpserver.server {
	// 	if rp.sinfo.Host == r.Host {
	// 		rp.Proxy.ServeHTTP(w, r)
	// 		log.Println("###rps.ServeHTTP success###")
	// 	}
	// }
	// }
	for _, rvproxy := range rpserver.server {
		log.Printf("host %s rvproxy : %v", r.Host, rvproxy)
		if rvproxy.sinfo.Host == r.Host {
			rvproxy.Proxy.ServeHTTP(w, r)
			break
		}
	}
}

// func (rps *RPSever) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	log.Println("###rps.ServeHTTP start###")
// 	log.Printf("rps.Host : %v\n", rps.Host)
// 	log.Printf("requst.Host : %v\n", r.Host)
// 	if rps.Host == r.Host {
// 		rps.Proxy.ServeHTTP(w, r)
// 		log.Println("###rps.ServeHTTP success###")
// 	}
// 	// rps.Proxy.ServeHTTP(w, r)
// }

// Directing ReverseProxy and append to array`
func (rpserver *RpServerArray) RvProxyHandle(serverinfo *ServerInfo) {
	log.Println("###RvProxyHandle start###")
	//target 설정
	//target_url은 api로 부터 입력받은 host:port. ex) http://www.facebook.com:10820
	target_url := serverinfo.Host + ":" + strconv.Itoa(serverinfo.Port)
	target, err := url.Parse(target_url)
	if err != nil {
		panic(err)
	}
	log.Printf("taget_url : %v\n", target_url)
	log.Printf("target : %v\n", target)
	log.Printf("target.Scheme : %v\n", target.Scheme)
	log.Printf("target.Host : %v\n", target.Host)
	log.Printf("target.Port : %v\n", target.Port())

	rvp := httputil.NewSingleHostReverseProxy(target)
	rvp.Director = func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", serverinfo.Origin)
		r.URL.Scheme = "http"
		r.URL.Host = serverinfo.Origin

		log.Println("########Reverse Proxy########")
		log.Printf("%v -> %v \n", serverinfo.Host, r.URL.Host)
	}

	mux := http.NewServeMux()
	log.Println("mux init")
	// rvproxy := rpserver.Proxy

	var rps RPSever
	rps.Proxy = rvp
	rps.sinfo.Host = serverinfo.Host
	rpserver.server = append(rpserver.server, rps)

	// var rp RpServerArray

	mux.HandleFunc("/", rpserver.ServeHTTP)
	// mux.HandleFunc("/", rps.Proxy.ServeHTTP)

	// for _, rvproxy := range rpserver.Proxy {
	// 	log.Println("rpserver.proxy loop start")
	// 	rvproxy.Proxy = rvp
	// 	rvproxy.Host = serverinfo.Host
	// 	mux.HandleFunc("/", rvproxy.Proxy.ServeHTTP)
	// 	log.Printf("rvproxy.proxy : %v\n", rvproxy.Proxy)
	// 	log.Printf("rvproxy.Host : %v\n", rvproxy.Host)
	// }

	log.Printf("rpserver.Proxy : %v", *rvp)
	// for _, rvproxy := range rpserver.Proxy {

	if !rpserver.Port[serverinfo.Port] {
		rpserver.Port[serverinfo.Port] = true
		log.Printf("check port : %v\n", rpserver.Port[serverinfo.Port])
		go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverinfo.Port), mux))
		log.Printf("%d port server started\n", serverinfo.Port)
	} else {
		log.Printf("%d port already use\n", serverinfo.Port)
	}
	// }

}

//init
func RvInit() *RpServerArray {
	log.Println("###RvInit start###")

	var rpserver RpServerArray
	mux := http.NewServeMux()

	mux.HandleFunc("/", rpserver.ServeHTTP)

	return &rpserver
}
