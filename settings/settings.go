package settings

import (
	"encoding/json"
	"github.com/CRVV/p2pFileSystem/logger"
	"io/ioutil"
	"sync"
)

var settings Settings
var mutex sync.Mutex = sync.Mutex{}

type Settings struct {
	Username  string
	GroupName string
	SharePath string
	Ignore    map[string]bool
}

func GetSettings() *Settings {
	mutex.Lock()
	if settings.Username == "" {
		configFile, err := ioutil.ReadFile("config.json")
		if err != nil {
			settings = Settings{"crvv", "Group", ".", make(map[string]bool)}
			settings.saveSettings()
		} else {
			settings = Settings{}
			err = json.Unmarshal(configFile, &settings)
			logger.Error(err)
		}
	}
	mutex.Unlock()
	return &settings
}
func (s *Settings) saveSettings() {
	configFile, err := json.Marshal(s)
	logger.Warning(err)
	ioutil.WriteFile("config.json", configFile, 0644)
}

func (s *Settings) GetSharePath() string {
	mutex.Lock()
	defer mutex.Unlock()
	return s.SharePath
}
func (s *Settings) SetSharePath(path string) {
	mutex.Lock()
	s.SharePath = path
	s.saveSettings()
	mutex.Unlock()
}
func (s *Settings) IsIgnored(name string) bool {
	mutex.Lock()
	_, ok := s.Ignore[name]
	mutex.Unlock()
	if name[0] == '.' {
		ok = true
	}
	return ok
}
func (s *Settings) GetIgnoreList() []string {
	list := make([]string, 8)
	mutex.Lock()
	for name, _ := range s.Ignore {
		list = append(list, name)
	}
	mutex.Unlock()
	return list
}
func (s *Settings) AddIgnore(name string) {
	mutex.Lock()
	s.Ignore[name] = true
	s.saveSettings()
	mutex.Unlock()
}
func (s *Settings) DeleteIgnore(name string) {
	mutex.Lock()
	delete(s.Ignore, name)
	s.saveSettings()
	mutex.Unlock()
}

func (s *Settings) GetUsername() string {
	mutex.Lock()
	defer mutex.Unlock()
	return s.Username
}
func (s *Settings) SetUsername(name string) {
	mutex.Lock()
	s.Username = name
	s.saveSettings()
	mutex.Unlock()
}
func (s *Settings) GetGroupName() string {
	mutex.Lock()
	defer mutex.Unlock()
	return s.GroupName
}
func (s *Settings) SetGroupName(name string) {
	mutex.Lock()
	s.GroupName = name
	s.saveSettings()
	mutex.Unlock()
}
