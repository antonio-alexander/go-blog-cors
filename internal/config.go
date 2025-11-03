package internal

import (
	"flag"
	"strconv"
	"strings"
)

const (
	EnvNameAddress          string = "ADDRESS"
	EnvNamePort             string = "PORT"
	EnvNameProxyAddress     string = "PROXY_ADDRESS"
	EnvNameProxyPort        string = "PROXY_PORT"
	EnvNameRedirectAddress  string = "REDIRECT_ADDRESS"
	EnvNameRedirectPort     string = "REDIRECT_PORT"
	EnvNameAllowCredentials string = "ALLOW_CREDENTIALS"
	EnvNameAllowedOrigins   string = "ALLOWED_ORIGINS"
	EnvNameAllowedMethods   string = "ALLOWED_METHODS"
	EnvNameAllowedHeaders   string = "ALLOWED_HEADERS"
	EnvNameCorsDebug        string = "CORS_DEBUG"
	EnvNameUsername         string = "USERNAME"
	EnvNamePassword         string = "PASSWORD"
	EnvNameCors             string = "CORS"
	EnvNameCorsOrigin       string = "CORS_ORIGIN"
)

var (
	configCorsOrigin       string = "http://localhost:10000"
	configAddress          string = "localhost"
	configPort             string = "8080"
	configUsername         string = "username"
	configPassword         string = "password"
	configAllowCredentials bool   = false
	configCorsDebug        bool   = false
	configCors             string = ""
	configProxyAddress     string = ""
	configProxyPort        string = ""
	configRedirectAddress  string = ""
	configRedirectPort     string = ""
	configAllowedMethods          = []string{}
	configAllowedHeaders          = []string{}
	configAllowedOrigins          = []string{}
)

func configFromEnvs(envs map[string]string) {
	if _, ok := envs[EnvNameAddress]; ok {
		configAddress = envs[EnvNameAddress]
	}
	if _, ok := envs[EnvNamePort]; ok {
		configPort = envs[EnvNamePort]
	}
	if _, ok := envs[EnvNameProxyAddress]; ok {
		configProxyAddress = envs[EnvNameProxyAddress]
	}
	if _, ok := envs[EnvNameProxyPort]; ok {
		configProxyPort = envs[EnvNameProxyPort]
	}
	if _, ok := envs[EnvNameRedirectAddress]; ok {
		configRedirectAddress = envs[EnvNameRedirectAddress]
	}
	if _, ok := envs[EnvNameRedirectPort]; ok {
		configRedirectPort = envs[EnvNameRedirectPort]
	}
	if s, ok := envs[EnvNameAllowedOrigins]; ok && s != "" {
		configAllowedOrigins = strings.Split(s, ",")
	}
	if s, ok := envs[EnvNameAllowedMethods]; ok && s != "" {
		configAllowedMethods = strings.Split(s, ",")
	}
	if s, ok := envs[EnvNameAllowedHeaders]; ok && s != "" {
		configAllowedHeaders = strings.Split(s, ",")
	}
	if s, ok := envs[EnvNameAllowCredentials]; ok && s != "" {
		configAllowCredentials, _ = strconv.ParseBool(s)
	}
	if _, ok := envs[EnvNameUsername]; ok {
		configUsername = envs[EnvNameUsername]
	}
	if _, ok := envs[EnvNamePassword]; ok {
		configPassword = envs[EnvNamePassword]
	}
	if _, ok := envs[EnvNameCors]; ok {
		configCors = envs[EnvNameCors]
	}
	if _, ok := envs[EnvNameCorsDebug]; ok {
		configCorsDebug, _ = strconv.ParseBool(envs[EnvNameCorsDebug])
	}
	if _, ok := envs[EnvNameCorsOrigin]; ok {
		configCorsOrigin = envs[EnvNameCorsOrigin]
	}
}

func configFromArgs(args []string) error {
	var address, port string

	serverFlags := flag.NewFlagSet("server", flag.ContinueOnError)
	serverFlags.StringVar(&address, "address", "", "the listening address")
	serverFlags.StringVar(&port, "port", "", "the listening port")
	if err := serverFlags.Parse(args); err != nil {
		return err
	}
	if address != "" {
		configAddress = address
	}
	if port != "" {
		configPort = port
	}
	return nil
}
