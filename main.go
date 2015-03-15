package main

import (
	"flag"
	"fmt"
	"github.com/ihsw/deploy-symfony2-client/Config"
)

func main() {
	configPath := flag.String("config", "", "Config path")
	flag.Parse()

	var (
		file Config.File
		err  error
	)
	if file, err = Config.NewFile(*configPath); err != nil {
		fmt.Println(fmt.Sprintf("Config.NewFile() failed: %s", err.Error()))
		return
	}

	var wrapper Config.Wrapper
	if wrapper, err = Config.NewWrapper(file); err != nil {
		fmt.Println(fmt.Sprintf("Config.NewWrapper() failed: %s", err.Error()))
		return
	}

	fmt.Println(fmt.Sprintf("wrapper: %v", wrapper))
}
