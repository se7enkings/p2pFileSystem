package ndp

func GetPeerList() map[string]Peer {
	plMutex.Lock()
	defer plMutex.Unlock()
	return peerList
}
func GetPeerAddr(name string) string {
	plMutex.Lock()
	defer plMutex.Unlock()
	return peerList[name].Addr
}
func GetPeerFromJson(message []byte) (Peer, error) {
	peer, err := json2peer(message)
	return peer, err
}
