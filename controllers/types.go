package controllers

type GroupId uint8
type DeviceId string

type ControllerConfig struct {
	Hue *HueConfig `yaml:"hue"`

	Groupings []*DeviceGroup `yaml:"groupings"`
}

type DeviceGroup struct {
	Devices []DeviceId
	Name    string
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
