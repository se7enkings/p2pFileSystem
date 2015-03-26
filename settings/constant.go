package settings

const MaxMessageSize = ^uint32(0)

//const BlockSize = 4194304 // 4 M
const MessageHeaderSize = 4
const NeighborDiscoveryMessageBufferSize = 1024

const BroadcastAddress = "255.255.255.255"
const CommunicationPort = ":1539"
const NeighborDiscoveryPort = ":1540"

const FileSystemListProtocol = "fslp"
const FileSystemRequestProtocol = "fsrp"
const NeighborDiscoveryProtocol = "ndpl"
const NeighborDiscoveryProtocolEcho = "ndpe"

const InvalidUsername = "iune"
