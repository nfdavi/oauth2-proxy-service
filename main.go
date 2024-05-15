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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		authorizationFail := false
		proxy := func() {
			log.Print("proxying " + r.RequestURI)

			token, err := getToken()
			if err != nil {
				log.Print(err)
				return
			}

			req, err := http.NewRequest(r.Method, settings.Proxy.Server+r.RequestURI, r.Body)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Add("Authorization", token.TokenType+" "+token.AccessToken)
			resp, err := http.DefaultClient.Do(req)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Print(err)
				}
			}(resp.Body)

			if resp.StatusCode == http.StatusUnauthorized && !authorizationFail {
				authorizationFail = true
				return
			}

			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Print(err)
				return
			}

			for k, v := range resp.Header {
				for _, vv := range v {
					w.Header().Set(k, vv)
				}
			}

			w.WriteHeader(resp.StatusCode)
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
