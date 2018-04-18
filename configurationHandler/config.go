/**
package configurationHandler

This package handles loading in configuration data from .butsted.toml, and command line flags.
It is also responsible for setting the logger for the project in the global package.

using the pflag and viper libraries for command line options and configuration files
*/
package configurationHandler

import (
	"GenomeBustersBackend/global"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// configuration is the object representing the current configuration of the program
// Its initialization is run at package import time, prior to the init method.
var configuration = initializeConfiguration()

// GetConfig retrives configuration for the server
func GetConfig() *viper.Viper {
	return configuration
}

// init rebuilds the front end if the configuration states to
func init() {
	//initializeVisualOutput()
	if !configuration.GetBool("skipRebuild") {
		global.Log.Println("Rebuilding Frontend")
		configureFrontend(configuration)
		npmBuild := exec.Command("npm", "run-script", "build")
		npmBuild.Dir = configuration.GetString("serverRoot")
		npmBuild.Stderr = os.Stderr
		if err := npmBuild.Run(); err != nil {
			panic(err)
		}
		global.Log.Println("Frontend rebuilt")
	}
}

// configureFrontend wull rebuild the /src/config.js file based on the busted.toml file.
// This way the front end can be forced to respect certian variables for the backend such
// as port mapping.
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

// initializeConfiguration sets up the potential configuration values that can be used, as well as default values for those variables.
// Generally speaking, if a new configuration option is needed, place it in here prior to "flag.parse()".
// Note that command line options override config file options
func initializeConfiguration() *viper.Viper {
	flag.Int("apiPort", 8080, "Overids the configuration files port for")
	flag.Bool("skipRebuild", false, "Skips rebuilding the frontend")
	flag.String("serverRoot", "./GenomeBusters/polymorphs-frontend-master", "Sets the location of the front end source code root")
	flag.Int("port", 80, "Sets the listening port for the server")
	flag.String("apiAddress", "127.0.0.1", "Address for the api server the front end should look for")
	flag.String("LogFile", "busted.log", "File to save log to. Defaults to './busted.log'")
	flag.Bool("color-256", true, "Determines whether or not to use 256 colors. Windows consoles do not support this, and as such this will have no effect there. Defaults to true")
	flag.Bool("LogToConsole", false, "Instead of getting the interactive prompt, print log to stdout. Exit with ^c")
	flag.Bool("help", false, "Show the help text")
	flag.Parse()

	v := viper.New()
	v.BindPFlags(flag.CommandLine)
	err := readInConfig(v)
	if err != nil {
		global.Log.Fatalf("Unable to parse configurations file or arguments, try \"busted --help\"\n%v", err)
	}

	if v.GetBool("help") {
		global.Log.Fatalln(flag.ErrHelp)
	}

	return v
}

// readInConfig reads in the .busted.toml and applies configuration options.
func readInConfig(v *viper.Viper) error {
	v.SetConfigType("toml")
	v.SetConfigName(".busted")
	if runtime.GOOS != "windows" {
		v.AddConfigPath("$XDG_CONFIG_HOME/.config/busted/")
		v.AddConfigPath("/etc/busted/")
	}
	v.AddConfigPath(".")

	LogFilePath := v.GetString("LogFile")
	LogFile, err := os.Create(LogFilePath)
	if err != nil {
		global.Log = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		log.Fatalf("Unable to open log file at %s:\n\t%v", LogFilePath, err)
	} else {
		if v.GetBool("LogToConsole") {
			global.Log = log.New(io.MultiWriter(LogFile, os.Stderr), "", log.Ldate|log.Ltime|log.Lmicroseconds)
		} else {
			global.Log = log.New(LogFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
		}
	}

	return v.ReadInConfig()
}
