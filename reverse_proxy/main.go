package main

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/AdityaP1502/Instant-Messanging/reverse_proxy/config"
	"github.com/gorilla/mux"
)

func ForwardClientCertMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			cert := r.TLS.PeerCertificates[0]
			block := &pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			}
			certBytes := pem.EncodeToMemory(block)
			encodedCert := base64.StdEncoding.EncodeToString(certBytes)

			r.Header.Add("x-client-cert", encodedCert)
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	path := "config/app.config.json"
	config, err := config.ReadJSONConfiguration(path)

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy\n"))
		w.WriteHeader(200)
	}).Methods("GET")

	// Always forward client certificate to endpoint
	r.Use(ForwardClientCertMiddleware)

	ver := r.PathPrefix(fmt.Sprintf("/%s", config.Version)).Subrouter()

	// authEndpoint, err := url.Parse(
	// 	fmt.Sprintf(
	// 		"%s://%s:%d",
	// 		config.Services.Account.Scheme,
	// 		config.Services.Auth.Host,
	// 		config.Services.Auth.Port,
	// 	),
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(authEndpoint)

	// authProxy := httputil.NewSingleHostReverseProxy(authEndpoint)

	// auth := ver.PathPrefix("/auth").Subrouter()

	// auth.Use(ForwardClientCertMiddleware)
	// auth.HandleFunc("/{rest:[a-zA-Z0-9=\\-\\/]+}", func(w http.ResponseWriter, r *http.Request) {
	// 	r.URL.Host = authEndpoint.Host
	// 	r.URL.Scheme = "https"
	// 	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	// 	r.Host = authEndpoint.Host

	// 	authProxy.ServeHTTP(w, r)
	// })

	// accountEndpoint, err := url.Parse(
	// 	fmt.Sprintf(
	// 		"%s://%s:%d",
	// 		config.Services.Account.Scheme,
	// 		config.Services.Account.Host,
	// 		config.Services.Account.Port,
	// 	),
	// )

	// fmt.Println(accountEndpoint)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// accountProxy := httputil.NewSingleHostReverseProxy(accountEndpoint)

	// account := ver.PathPrefix("/account").Subrouter()
	// account.HandleFunc("/{rest:[a-zA-Z0-9=\\-\\/]+}", func(w http.ResponseWriter, r *http.Request) {
	// 	r.URL.Host = accountEndpoint.Host
	// 	r.URL.Scheme = "https"
	// 	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	// 	r.Host = accountEndpoint.Host

	// 	accountProxy.ServeHTTP(w, r)
	// })

	for _, service := range config.Services {
		url, err := url.Parse(fmt.Sprintf("%s://%s:%d", service.Scheme, service.Host, service.Port))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Assigning %s service (%s://%s/%s/%s) to the proxy\n",
			service.Endpoint,
			url.Scheme,
			config.Version,
			url.Host,
			service.Endpoint,
		)

		proxy := httputil.NewSingleHostReverseProxy(url)
		subrouter := ver.PathPrefix(fmt.Sprintf("/%s", service.Endpoint)).Subrouter()

		subrouter.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
			r.URL.Host = url.Host
			r.URL.Scheme = url.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = url.Host

			proxy.ServeHTTP(w, r)
		})

		subrouter.HandleFunc("/{rest:.*}", func(w http.ResponseWriter, r *http.Request) {
			r.URL.Host = url.Host
			r.URL.Scheme = url.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = url.Host

			proxy.ServeHTTP(w, r)
		})
	}

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
			r,
		)
	}()

	fmt.Printf("Reverse proxy running on %s:%d\n", config.Server.Host, config.Server.Port)
	wg.Wait()
}
