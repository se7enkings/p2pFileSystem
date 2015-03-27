package ndp

func GetPeerList() map[string]Peer {
	PlMutex.Lock()
	defer PlMutex.Unlock()
	return peerList
}
func GetPeerAddr(name string) string {
	PlMutex.Lock()
	defer PlMutex.Unlock()
	return peerList[name].Addr
}
func GetPeerFromJson(message []byte) (Peer, error) {
	peer, err := json2peer(message)
	return peer, err
}
