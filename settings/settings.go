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
	Ignore    map[string]int
}

func GetSettings() *Settings {
	mutex.Lock()
	if settings.Username == "" {
		configFile, err := ioutil.ReadFile("config.json")
		if err != nil {
			settings = Settings{"crvv", "Group", "test/testLocalFolder", make(map[string]int)}
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
	mutex.Lock()
	configFile, err := json.Marshal(s)
	mutex.Unlock()
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
	s.Ignore[name] = 1
	mutex.Unlock()
	s.saveSettings()
}
func (s *Settings) DeleteIgnore(name string) {
	mutex.Lock()
	delete(s.Ignore, name)
	mutex.Unlock()
	s.saveSettings()
}

func (s *Settings) GetUsername() string {
	mutex.Lock()
	defer mutex.Unlock()
	return s.Username
}
func (s *Settings) SetUsername(name string) {
	mutex.Lock()
	s.Username = name
	mutex.Unlock()
	s.saveSettings()
}
func (s *Settings) GetGroupName() string {
	mutex.Lock()
	defer mutex.Unlock()
	return s.GroupName
}
func (s *Settings) SetGroupName(name string) {
	mutex.Lock()
	s.GroupName = name
	mutex.Unlock()
	s.saveSettings()
}
