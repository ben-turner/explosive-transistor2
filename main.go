package main

import (
	"errors"
	"fmt"
	"github.com/ben-turner/explosive-transistor2/controllers"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	ConfigEnv = "ET_CONFIG"
)

var (
	NoConfigError = errors.New("No configuration file specified")

	CurrentSong = "No Song"
)

func loadConfig() (*Config, error) {
	filename := os.Getenv(ConfigEnv)
	if filename == "" {
		return nil, NoConfigError
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err := yaml.Unmarshal(file, c); err != nil {
		return nil, err
	}

	return c, nil
}

func Run() int {
	c, err := loadConfig()
	if err != nil {
		log.Println("Failed to load Config:", err.Error())
		return 1
	}
	log.Println("Loaded config:")
	pretty.Printf("%#v\n", c)

	devs := make(map[string]controllers.Controller)
	for n, d := range c.Controllers {
		devs[n] = controllers.NewHue(d.Hue, d.Groups)
	}

	AddWatches(devs)

	http.HandleFunc("/cur_song/", getSong)
	http.HandleFunc("/set_song/", setSong)
	http.HandleFunc("/log/", writeLog)

	api := GetApiHandler(devs)
	http.Handle("/api/", api)
	fs := http.FileServer(http.Dir(c.WebViewDir))
	http.Handle("/", fs)
	err = http.ListenAndServe(c.ServerHost, nil)
	log.Println(err.Error())
	return 1
}

func getSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, CurrentSong)
}

func setSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	CurrentSong = string(b)
}

func writeLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	b, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	log.Println(string(b))
}

func main() {
	os.Exit(Run())
}
