package controllers

// import (
// 	"github.com/tarm/serial"
// )

// const (
// 	arduinoRead  byte = 0x00
// 	arduinoWrite byte = 0x01
// )

// var arduinoPrefix = []byte{0x02, 0x02, 0x02}

// type ArduinoConfig struct {
// 	SerialPort string
// 	BaudRate   int
// }

// type arduino struct {
// 	*ArduinoConfig

// 	groups []*DeviceGroup
// 	port   *serial.Port
// }

// func NewArduino(c *ArduinoConfig, g []*DeviceGroup) (Controller, error) {
// 	p, err := serial.OpenPort(&serial.Config{
// 		Name: c.SerialPort,
// 		Baud: c.BaudRate,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &arduino{
// 		ArduinoConfig: c,
// 		groups:        g,
// 		port:          p,
// 	}, nil
// }

// func (a *arduino) Set(deviceId GroupId, s *State) error {
// 	msg := append(arduinoPrefix, arduinoWrite, byte(d.Ports[port]), byte(val))
// 	err := d.write(msg)
// }

// func (a *arduino) Get(deviceId GroupId) (*State, error) {

// }

// func (a *arduino) List() (map[GroupId]string, error) {

// }
