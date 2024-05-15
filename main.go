package main

import (
	"io"
	"log"
	"net/http"
)

var settings Settings

func main() {
	log.Print("Starting OAuth2 proxy service")
	settings = loadSettings("proxy.ini")

	http.HandleFunc("/", func(w http.ResponseWriter, incomingReq *http.Request) {
		authorizationFail := false
		proxy := func() {
			log.Print("proxying " + incomingReq.RequestURI)

			token, err := getToken()
			if err != nil {
				log.Print(err)
				return
			}

			proxyReq, err := http.NewRequest(incomingReq.Method, settings.Proxy.Server+incomingReq.RequestURI, incomingReq.Body)
			if err != nil {
				log.Fatal(err)
			}
			proxyReq.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)
			proxyResp, err := http.DefaultClient.Do(proxyReq)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Print(err)
				}
			}(proxyResp.Body)

			if proxyResp.StatusCode == http.StatusUnauthorized && !authorizationFail {
				authorizationFail = true
				return
			}

			respBytes, err := io.ReadAll(proxyResp.Body)
			if err != nil {
				log.Print(err)
				return
			}

			for k, v := range proxyResp.Header {
				for _, vv := range v {
					w.Header().Set(k, vv)
				}
			}

			w.WriteHeader(proxyResp.StatusCode)
			_, err = w.Write(respBytes)
			if err != nil {
				log.Print(err)
				return
			}
		}
		proxy()
		if authorizationFail {
			log.Print("retrying proxy")
			invalidateToken()
			proxy()
		}
	})

	log.Fatal(http.ListenAndServe(settings.Webservice.Bind, nil))
}
