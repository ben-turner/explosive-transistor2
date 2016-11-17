package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type requestMethod string

const (
	Get    requestMethod = "GET"
	Post   requestMethod = "POST"
	Put    requestMethod = "PUT"
	Delete requestMethod = "DELETE"
)

var (
	BadHueResponseError error = errors.New("Bad response body from Hue")
	GetBodyNotNilError  error = errors.New("Get request body must be blank")
)

type HueConfig struct {
	BridgeIp string `yaml:"bridgeIp"`
	Key      string `yaml:"key"`
}

type hue struct {
	*HueConfig
	groups []*DeviceGroup
}

func NewHue(c *HueConfig, g []*DeviceGroup) Controller {
	return &hue{
		c,
		g,
	}
}

func (d *hue) sendRequest(m requestMethod, path string, body []byte) ([]byte, error) {
	if m == Get && body != nil {
		return nil, GetBodyNotNilError
	}

	url := fmt.Sprintf("http://%v/api/%v/%v", d.BridgeIp, d.Key, path)

	client := &http.Client{}
	request, err := http.NewRequest(string(m), url, bytes.NewReader(body))
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func (d *hue) set(dev DeviceId, s *State) error {
	log.Println("Setting")
	path := fmt.Sprintf("lights/%v/state", string(dev))

	bri, ok := s.Values["bri"]
	var body string
	if ok {
		body = fmt.Sprintf(`{"bri": %d, "on": %v}`, bri, s.On)
	} else {
		body = fmt.Sprintf(`{"on": %v}`, s.On)
	}

	res, err := d.sendRequest(Put, path, []byte(body))
	if err != nil {
		return err
	}

	log.Println(string(res))

	var resParsed []struct {
		Success *struct{}
		Error   *struct {
			Description string
		}
	}
	err = json.Unmarshal(res, &resParsed)
	if err != nil {
		return err
	}
	if resParsed[0].Error != nil {
		return errors.New(resParsed[0].Error.Description)
	}
	if resParsed[0].Success == nil {
		return BadHueResponseError
	}
	return nil
}

func (d *hue) Set(g GroupId, s *State) error {
	for _, dev := range d.groups[g].Devices {
		log.Println("Setting", dev, s.On)
		if err := d.set(dev, s); err != nil {
			return err
		}
	}
	return nil
}

func (d *hue) get(light string) (*State, error) {
	path := "lights/" + light + "/"
	res, err := d.sendRequest(Get, path, nil)
	if err != nil {
		return nil, err
	}

	log.Println(string(res))

	var parsedRes struct {
		State *struct {
			On  bool  `json:"on"`
			Bri uint8 `json:"bri"`
		} `json:"state"`
	}

	err = json.Unmarshal(res, &parsedRes)
	if err != nil {
		return nil, err
	}

	return &State{
		On: parsedRes.State.On,
		Values: map[string]uint8{
			"Bri": parsedRes.State.Bri,
		},
	}, nil
}

func (d *hue) Get(group GroupId) (*State, error) {
	groupVal := d.groups[group]
	return d.get(string(groupVal.Devices[0]))
}

func (d *hue) List() (map[GroupId]string, error) {
	return nil, errors.New("Not implemented")
	// path := "lights"
	// res, err := d.sendRequest(Get, path, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// var parsedRes map[string]*struct {
	// 	Name string
	// }
	// err = json.Unmarshal(res, &parsedRes)
	// if err != nil {
	// 	return nil, err
	// }

	// out := make(map[DeviceId]string)
	// for id, nameStruct := range parsedRes {
	// 	out[DeviceId(id)] = nameStruct.Name
	// }
	// return out, nil
}
