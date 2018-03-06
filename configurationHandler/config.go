package configurationHandler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var configuration *viper.Viper
var err error

// GetConfig retrives configuration for the server
func GetConfig() (*viper.Viper, error) {
	return configuration, err
}

func init() {
	//initializeVisualOutput()
	configuration, err = initializeConfiguration()
	if !configuration.GetBool("skipRebuild") {
		fmt.Println("Rebuilding Frontend")
		configureFrontend(configuration)
		npmBuild := exec.Command("npm", "run-script", "build")
		npmBuild.Dir = configuration.GetString("serverRoot")
		npmBuild.Stderr = os.Stderr
		if err = npmBuild.Run(); err != nil {
			panic(err)
		}
		fmt.Println("Frontend rebuilt")
	}
}

/*
func initializeVisualOutput() {
	err := termui.Init()
	if err == nil {
		go visualOut()
	} else {
		fmt.Printf("error occurred when initializing termui %v", err)
	}
}
*/

func configureFrontend(v *viper.Viper) {
	filePath := v.GetString("serverRoot") + "/src/config.js"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		words := strings.Fields(line)

		if len(words) >= 2 {
			switch words[1] {
			case "address":
				lines[i] = "    static address = \"" + v.GetString("apiAddress") + ":" + strconv.Itoa(v.GetInt("apiPort")) + "/api/gene_search\";"
				break
			}
		}
	}
	if err := ioutil.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		panic(err)
	}
}

func initializeConfiguration() (*viper.Viper, error) {
	flag.Int("apiPort", 8080, "Overids the configuration files port for")
	flag.Bool("skipRebuild", false, "Skips rebuilding the frontend")
	flag.String("serverRoot", "./GenomeBusters/polymorphs-frontend-master", "Sets the location of the front end source code root")
	flag.Int("port", 80, "Sets the listening port for the server")
	flag.String("apiAddress", "127.0.0.1", "Address for the api server the front end should look for")
	flag.Parse()
	v := viper.New()
	v.BindPFlags(flag.CommandLine)
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
