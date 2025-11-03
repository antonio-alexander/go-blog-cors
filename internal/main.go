package internal

import (
	"fmt"
	"os"
	"sync"
)

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
		fmt.Printf("configured with proxy/redirect: http://%s:%s\n", configProxyAddress, configProxyPort)
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
