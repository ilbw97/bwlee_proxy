package reverse

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type Rpserver_array struct {
	Proxy []*httputil.ReverseProxy
	Port  [65535]bool
}

type Information struct {
	Host   string `json:"host"`
	Origin string `json:"origin"`
	Port   int    `json:"port"`
}

func (rpserver *Rpserver_array) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(rpserver.Proxy) == 0 {
		log.Fatal(len(rpserver.Proxy))
	}

	for _, rproxy := range rpserver.Proxy {
		log.Println(rproxy)

		rproxy.ServeHTTP(w, r)
	}
}

func (rpserver *Rpserver_array) Rvproxy_handle(info *Information) *httputil.ReverseProxy {
	target_url := info.Host + ":" + strconv.Itoa(info.Port)
	target, err := url.Parse(target_url)
	if err != nil {
		panic(err)
	}

	log.Printf("taget_url : %v\n", target_url)
	log.Printf("taget : %v\n", target)
	log.Printf("target.Scheme : %v\n", target.Scheme)
	log.Printf("target.Host : %v\n", target.Host)
	log.Printf("target.Port : %v\n", target.Port())
	rvproxy := httputil.NewSingleHostReverseProxy(target)

	rvproxy.Director = func(r *http.Request) {
		r.Header.Add("X-Forwarded-Host", r.Host)
		r.Header.Add("X-Origin-Host", info.Origin)
		r.URL.Scheme = "http"
		r.URL.Host = info.Origin

		log.Printf("info.Origin : %v\n", info.Origin)
		log.Printf("%v -> %v \n", info.Host, r.URL.Host)
	}

	//reverse proxy array에 저장하기
	// var rparray Rpserver_array
	// rparray.Proxy = append(rparray.Proxy, rvproxy)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", rparray.ServeHTTP)
	rpserver.Proxy = append(rpserver.Proxy, rvproxy)

	if rpserver.Port[info.Port] == false {
		rpserver.Port[info.Port] = true
		go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", info.Port), nil))
	}

	// if target.Port() !=

	// Reverse_init()

	// var rparray Rpserver_array
	// rparray.proxy = append(rparray.proxy, rvproxy)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", rparray.ServeHTTP)

	// log.Printf("rpserver port : %v\n", info.Port)
	return rvproxy
}

func Reverse_init() *Rpserver_array {
	var array Rpserver_array
	mux := http.NewServeMux()
	mux.HandleFunc("/", array.ServeHTTP)
	return &array
}
