package ndp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CRVV/p2pFileSystem/settings"
	"math/rand"
	"sync"
	"time"
)

var myself Peer

func GetPeerList() map[string]Peer {
	return peerList.GetMap()
}
func GetPeerAddr(name string) (string, error) {
	peerList.RLock()
	defer peerList.RUnlock()
	_, ok := peerList.M[name]
	if ok {
		return peerList.M[name].Addr, nil
	}
	return "", errors.New(fmt.Sprintf("requested address of an unknown peer %s", name))
}
func GetPeerFromJson(message []byte) (Peer, error) {
	peer, err := json2peer(message)
	return peer, err
}

type peerTable struct {
	M map[string]Peer //key: Username
	sync.RWMutex
}

func (t *peerTable) Exist(name string) bool {
	t.RLock()
	_, ok := t.M[name]
	t.RUnlock()
	return ok
}
func (t *peerTable) GetMap() map[string]Peer {
	t.RLock()
	m := make(map[string]Peer)
	for name, peer := range t.M {
		m[name] = peer
	}
	t.RUnlock()
	return m
}
func (t *peerTable) ReplaceByNewMap(m map[string]Peer) {
	t.Lock()
	t.M = m
	t.Unlock()
}
func (t *peerTable) Delete(name string) {
	t.Lock()
	defer t.Unlock()
	delete(t.M, name)
}
func (t *peerTable) Add(peer Peer) {
	t.Lock()
	defer t.Unlock()
	t.M[peer.Username] = peer
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
