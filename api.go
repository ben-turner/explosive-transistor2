package main

import (
	"encoding/json"
	"fmt"
	"github.com/ben-turner/explosive-transistor2/controllers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func wrapGetFunc(dev controllers.Controller) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		vars := mux.Vars(r)
		d, ok := vars["deviceId"]
		if !ok {
			fmt.Fprint(w, "No device ID provided")
			return
		}

		groupIndex, err := strconv.ParseUint(d, 10, 8)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		state, err := dev.Get(controllers.GroupId(groupIndex))
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		res, err := json.Marshal(state)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, string(res))
	}
}

func wrapSetFunc(dev controllers.Controller) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		vars := mux.Vars(r)
		d, ok := vars["deviceId"]
		if !ok {
			fmt.Fprint(w, "No device ID provided")
			return
		}

		stateBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		state := &controllers.State{}
		err = json.Unmarshal(stateBytes, state)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}

		groupIndex, err := strconv.ParseUint(d, 10, 8)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		err = dev.Set(controllers.GroupId(groupIndex), state)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, "ok")
	}
}

func GetApiHandler(devs map[string]controllers.Controller) http.Handler {
	r := mux.NewRouter()
	for n, d := range devs {
		log.Printf("Creating route: /api/%v/{deviceId}/", n)
		r.HandleFunc("/api/"+n+"/{deviceId}", wrapGetFunc(d)).Methods("GET")
		r.HandleFunc("/api/"+n+"/{deviceId}", wrapSetFunc(d)).Methods("POST")
	}
	return r
}
