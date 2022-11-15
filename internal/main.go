package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func endpointHelloWorld(writer http.ResponseWriter, request *http.Request) {
	//TODO: check headers to confirm that GET method and ORIGIN is allowed
	bytes, _ := json.Marshal(&Message{
		Message: "Hello, World!",
	})
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := writer.Write(bytes); err != nil {
		fmt.Printf("error (endpointHelloWorld): %s", err)
	}
}

func endpointAuthorize(writer http.ResponseWriter, request *http.Request) {
	//TODO: check headers to confirm that POST method, ORIGIN is allowed,
	// Allow Credentials and the Authorization header is allowed
	username, password, ok := request.BasicAuth()
	switch {
	default:
		bytes, _ := json.Marshal(&Error{
			StatusCode: http.StatusUnauthorized,
			Error:      "invalid credentials",
		})
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.WriteHeader(http.StatusUnauthorized)
		if _, err := writer.Write(bytes); err != nil {
			fmt.Printf("error (endpointAuthorize): %s", err)
		}
	case !ok:
		bytes, _ := json.Marshal(&Error{
			StatusCode: http.StatusBadRequest,
			Error:      "unable to get basic auth",
		})
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.WriteHeader(http.StatusBadRequest)
		if _, err := writer.Write(bytes); err != nil {
			fmt.Printf("error (endpointAuthorize): %s", err)
		}
	case username == configUsername && password == configPassword:
		bytes, _ := json.Marshal(&Message{
			Message: "login successful",
		})
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.WriteHeader(http.StatusOK)
		if _, err := writer.Write(bytes); err != nil {
			fmt.Printf("error (endpointAuthorize): %s", err)
		}
	}
}

func endpointProxy(writer http.ResponseWriter, request *http.Request) {
	var uri string

	switch endpoint := request.URL.Query().Get("endpoint"); endpoint {
	default:
		uri = fmt.Sprintf("http://%s:%s/%s", configProxyAddress, configProxyPort, endpoint)
	case "":
		uri = fmt.Sprintf("http://%s:%s", configProxyAddress, configProxyPort)
	}
	fmt.Printf("attempting to proxy to: %s\n", uri)
	target, err := url.Parse(uri)
	if err != nil {
		bytes, _ := json.Marshal(&Error{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		writer.WriteHeader(http.StatusInternalServerError)
		if _, err := writer.Write(bytes); err != nil {
			fmt.Printf("error (endpointProxy): %s", err)
		}
		return
	}
	request.URL = target
	httputil.NewSingleHostReverseProxy(target).ServeHTTP(writer, request)
}

func launchServer(wg *sync.WaitGroup, chErr chan error) *http.Server {
	server := &http.Server{Addr: configAddress + ":" + configPort}
	started := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()

		router := mux.NewRouter()
		router.HandleFunc("/", endpointHelloWorld).Methods(http.MethodGet)
		router.HandleFunc("/authorize", endpointAuthorize).Methods(http.MethodPost)
		if configProxyAddress != "" && configProxyPort != "" {
			router.HandleFunc("/proxy", endpointProxy).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
		}
		switch strings.ToLower(configCors) {
		default:
			fmt.Println("cors disabled")
			server.Handler = router
		case "cors":
			fmt.Println("cors enabled (cors)")
			server.Handler = cors.New(cors.Options{
				AllowedOrigins:   configAllowedOrigins,
				AllowCredentials: configAllowCredentials,
				AllowedMethods:   configAllowedMethods,
				AllowedHeaders:   configAllowedHeaders,
				Debug:            configCorsDebug,
			}).Handler(router)
		case "gorilla":
			fmt.Println("cors enabled (gorilla)")
			var corsOptions []handlers.CORSOption

			if configAllowCredentials {
				corsOptions = append(corsOptions, handlers.AllowCredentials())
			}
			corsOptions = append(corsOptions, handlers.AllowedHeaders(configAllowedHeaders))
			corsOptions = append(corsOptions, handlers.AllowedMethods(configAllowedMethods))
			corsOptions = append(corsOptions, handlers.AllowedOrigins(configAllowedOrigins))
			server.Handler = handlers.CORS(corsOptions...)(router)
		}
		if configCors != "" {
			fmt.Printf("cors allowed origins: %s\n", strings.Join(configAllowedOrigins, ","))
			fmt.Printf("cors allow credentials: %t\n", configAllowCredentials)
			fmt.Printf("cors allowed methods: %s\n", strings.Join(configAllowedMethods, ","))
			fmt.Printf("cors allowed headers: %s\n", strings.Join(configAllowedHeaders, ","))
		}
		close(started)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			chErr <- err
		}
		fmt.Println("closed server")
	}()
	<-started
	return server
}

func Main(pwd string, args []string, envs map[string]string, chSignalInt chan os.Signal) error {
	var err error

	chErr := make(chan error, 1)
	wg := new(sync.WaitGroup)
	configFromEnvs(envs)
	if err := configFromArgs(args); err != nil {
		return err
	}
	server := launchServer(wg, chErr)
	fmt.Printf("started server listening on %s:%s\n", configAddress, configPort)
	//KIM: this is incredibly insecure, if you copy+paste this beware...
	fmt.Printf("configured with basic auth: %s/%s\n", configUsername, configPassword)
	if configProxyAddress != "" && configProxyPort != "" {
		fmt.Printf("configured with proxy: http://%s:%s\n", configProxyAddress, configProxyPort)
	}
	select {
	case err = <-chErr:
	case osSignal := <-chSignalInt:
		fmt.Printf("signal received: %s\n", osSignal)
	}
	server.Close()
	wg.Wait()
	close(chErr)
	return err
}
