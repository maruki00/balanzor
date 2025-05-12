package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	pu, err := url.Parse("http://localhost:8081/")
	if err != nil {
		panic(err)
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(pu)

	//	server := http.NewerveMux()

	handler := func(reverseproxy *httputil.ReverseProxy) func(rw http.ResponseWriter, r *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Done"))
			reverseProxy.ServeHTTP(rw, r)
		}
	}

	http.HandleFunc("/lb", handler(reverseProxy))

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		panic(err)
	}

}
