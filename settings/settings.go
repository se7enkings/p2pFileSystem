package settings

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var settings Settings
var mutex sync.Mutex = sync.Mutex{}

type Settings struct {
	Username  string
	GroupName string
	SharePath string
	Ignore    map[string]int
}

func GetSettings() *Settings {
	mutex.Lock()
	if settings.Username == "" {
		configFile, err := ioutil.ReadFile("config.json")
		if err != nil {
			settings = Settings{"crvv", "Group", "test/testLocalFolder", make(map[string]int)}
			settings.SaveSettings()
		} else {
			settings = Settings{}
			err = json.Unmarshal(configFile, &settings)
			checkError(err)
		}
	}
	mutex.Unlock()
	return &settings
}
func (s *Settings) SaveSettings() {
	configFile, err := json.Marshal(s)
	// TODO: an error should not panic
	checkError(err)
	ioutil.WriteFile("config.json", configFile, 0644)
}

func (s *Settings) GetSharePath() string {
	return s.SharePath
}
func (s *Settings) SetSharePath(path string) {
	s.SharePath = path
	s.SaveSettings()
}
func (s *Settings) IsIgnored(name string) bool {
	_, ok := s.Ignore[name]
	return ok
}
func (s *Settings) GetIgnoreList() []string {
	list := make([]string, 8)
	for name, _ := range s.Ignore {
		list = append(list, name)
	}
	return list
}
func (s *Settings) AddIgnore(name string) {
	s.Ignore[name] = 1
	s.SaveSettings()
}
func (s *Settings) DeleteIgnore(name string) {
	delete(s.Ignore, name)
	s.SaveSettings()
}

func (s *Settings) GetUsername() string {
	return s.Username
}
func (s *Settings) SetUsername(name string) {
	s.Username = name
	s.SaveSettings()
}
func (s *Settings) GetGroupName() string {
	return s.GroupName
}
func (s *Settings) SetGroupName(name string) {
	s.GroupName = name
	s.SaveSettings()
	// TODO: changing group need to reconnect
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
