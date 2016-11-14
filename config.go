package main

import (
	"github.com/ben-turner/explosive-transistor2/controllers"
)

type TlsConfig struct {
	PublicKey  string `yaml:"publicKey"`
	PrivateKey string `yaml:"privateKey"`
}

type Config struct {
	Controllers map[string]*controllers.ControllerConfig `yaml:"controllers"`

	TlsConfig  *TlsConfig `yaml:"tlsConfig"`
	ServerHost string     `yaml:"serverHost"`
	WebViewDir string     `yaml:"webViewDir"`
}
