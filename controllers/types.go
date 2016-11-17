package controllers

type GroupId uint8
type DeviceId string

type ControllerConfig struct {
	Hue *HueConfig `yaml:"hue"`

	Groups []*DeviceGroup `yaml:"groups"`
}

type DeviceGroup struct {
	Devices []DeviceId `yaml:"devices"`
	Name    string     `yaml:"name"`
}

type Controller interface {
	Set(deviceId GroupId, s *State) error
	Get(deviceId GroupId) (*State, error)
	List() (map[GroupId]string, error)
}

type State struct {
	On     bool
	Values map[string]uint8
}
