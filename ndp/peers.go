package ndp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/CRVV/p2pFileSystem/settings"
	"math/rand"
	"sync"
	"time"
)

var myself Peer

func GetPeerList() *peerTable {
	return &peerList
}
func GetPeerAddr(name string) (string, error) {
	peerList.RLock()
	defer peerList.RUnlock()
	_, ok := peerList.M[name]
	if ok {
		return peerList.M[name].Addr, nil
	}
	return "", errors.New("requested address of an unknown peer")
}
func GetPeerFromJson(message []byte) (Peer, error) {
	peer, err := json2peer(message)
	return peer, err
}

type peerTable struct {
	M map[string]Peer //key: Username
	sync.RWMutex
}
type Peer struct {
	Username string
	Addr     string `json:"-"`
	ID       string
	Group    string
}

const ReloadPeerListNotice int = 1
const NewPeerNotice int = 2
const PeerMissingNotice int = 3

type PeerListNotice struct {
	NoticeType int
	PeerName   string
}

func peer2Json(message Peer) ([]byte, error) {
	b, err := json.Marshal(message)
	return b, err
}
func json2peer(message []byte) (Peer, error) {
	cm := Peer{}
	err := json.Unmarshal(message, &cm)
	return cm, err
}
func genID() string {
	rand.Seed(time.Now().UnixNano())
	idd := make([]byte, 16)
	for i, _ := range idd {
		idd[i] = byte(rand.Intn(256))
	}
	return base64.StdEncoding.EncodeToString(idd)
}
func Init() {
	myself = Peer{
		Username: settings.GetSettings().GetUsername(),
		Group:    settings.GetSettings().GetGroupName(),
		ID:       genID(),
	}
}
