package main

import (
	"GenomeBustersBackend/webserver"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("Starting Busted")
	v, err := initializeConfiguration()
	if err != nil {
		fmt.Printf("unable to parse config file %s", err)
		return
	}
	port := ":" + strconv.Itoa(v.GetInt("port"))
	apiport := ":" + strconv.Itoa(v.GetInt("apiPort"))

	fmt.Println("Rebuilding Frontend")
	configureFrontend(v)
	npmBuild := exec.Command("npm", "run-script", "build")
	npmBuild.Dir = v.GetString("serverRoot")
	npmBuild.Stderr = os.Stderr
	if err = npmBuild.Run(); err != nil {
		panic(err)
	}
	fmt.Println("Frontend rebuilt")

	fileServer := http.FileServer(http.Dir(v.GetString("serverRoot") + "/build"))
	// fh := http.Handle("/", fileServer)
	// http.HandleFunc("/api/gene_search/", webserver.GeneSearch)

	keyboardInterrupt := make(chan os.Signal, 1)
	signal.Notify(keyboardInterrupt, os.Interrupt)

	server := &http.Server{Addr: port, Handler: fileServer}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error: %s", err)
			keyboardInterrupt <- nil
		}
	}()

	apiServer := &http.Server{Addr: apiport, Handler: http.HandlerFunc(webserver.GeneSearch)}
	go func() {
		if err := apiServer.ListenAndServe(); err != nil {
			fmt.Printf("Error: %s", err)
			keyboardInterrupt <- nil
		}
	}()

	<-keyboardInterrupt
	fmt.Printf("\nShutting Down Server...\n")
	if err := server.Shutdown(nil); err != nil {
		fmt.Printf("Error: %s", err)
	}
	fmt.Printf("Goodbye!\n")
}

func configureFrontend(v *viper.Viper) {
	filePath := v.GetString("serverRoot") + "/src/config.js"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		words := strings.Fields(line)

		switch words[0] {
		case "address:":
			lines[i] = "address: \"" + v.GetString("apiAddress") + ":" + strconv.Itoa(v.GetInt("apiPort")) + "/api/gene_search\","
			break
		}
	}
	if err := ioutil.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		panic(err)
	}
}

func initializeConfiguration() (*viper.Viper, error) {
	v := viper.New()
	v.SetDefault("serverRoot", "./GenomeBusters/polymorphs-frontend-master")
	v.SetDefault("port", 80)
	v.SetDefault("apiPort", 8080)
	v.SetDefault("apiAddress", "127.0.0.1")
	err := readInConfig(v)
	return v, err
}

func readInConfig(v *viper.Viper) error {
	v.SetConfigType("toml")
	v.SetConfigName(".busted")
	if runtime.GOOS != "windows" {
		v.AddConfigPath("$XDG_CONFIG_HOME/.config/busted/")
		v.AddConfigPath("/etc/busted/")
	}
	v.AddConfigPath(".")
	return v.ReadInConfig()
}
