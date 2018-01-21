package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/spf13/viper"

	"github.com/GenomeBustersBackend/webserver"
)

func main() {
	fmt.Println("Starting Busted")
	v, err := initializeConfiguration()
	if err != nil {
		fmt.Printf("unable to parse config file %s", err)
		return
	}
	port := ":" + strconv.Itoa(v.GetInt("port"))

	fileServer := http.FileServer(http.Dir(v.GetString("serverRoot")))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/gene_search/", webserver.GeneSearch)

	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)

	server := &http.Server{Addr: port}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error: %s", err)
		}
	}()

	<-keyboardInterrupt
	fmt.Printf("\nShutting Down Server...\n")
	if err := server.Shutdown(nil); err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Goodbye!\n")
}

func initializeConfiguration() (*viper.Viper, error) {
	v := viper.New()
	v.SetDefault("serverRoot", "./GenomeBusters/polymorphs-frontend-master/build")
	v.SetDefault("port", 80)
	err := readInConfig(v)
	return v, err
}

func readInConfig(v *viper.Viper) error {
	v.SetConfigType("toml")
	v.SetConfigName(".busted")
	v.AddConfigPath("$XDG_CONFIG_HOME/.config/busted/")
	v.AddConfigPath("/etc/busted/")
	v.AddConfigPath(".")
	return v.ReadInConfig()
}
