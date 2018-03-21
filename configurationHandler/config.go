package configurationHandler

import (
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

var configuration = initializeConfiguration()

// GetConfig retrives configuration for the server
func GetConfig() *viper.Viper {
	return configuration
}

func init() {
	//initializeVisualOutput()
	if !configuration.GetBool("skipRebuild") {
		log.Println("Rebuilding Frontend")
		configureFrontend(configuration)
		npmBuild := exec.Command("npm", "run-script", "build")
		npmBuild.Dir = configuration.GetString("serverRoot")
		npmBuild.Stderr = os.Stderr
		if err := npmBuild.Run(); err != nil {
			panic(err)
		}
		log.Println("Frontend rebuilt")
	}
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
		log.Fatalf("Unable to parse configurations file or arguments, try \"busted --help\"\n%v", err)
	}

	if v.GetBool("help") {
		log.Fatalln(flag.ErrHelp)
	}

	return v
}

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
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	if err != nil {
		log.Fatalf("Unable to open log file at %s:\n\t%v", LogFilePath, err)
	}
	if v.GetBool("LogToConsole") {
		log.SetOutput(io.MultiWriter(LogFile, os.Stderr))
	} else {
		log.SetOutput(LogFile)
	}

	return v.ReadInConfig()
}
