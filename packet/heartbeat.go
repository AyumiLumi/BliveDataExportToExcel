<<<<<<< HEAD
package packet

// NewHeartBeatPacket 构造心跳包
func NewHeartBeatPacket() []byte {
	pkt := NewPacket(1, HeartBeat, nil)
	return pkt.Build()
}
=======
package packet

// NewHeartBeatPacket 构造心跳包
func NewHeartBeatPacket() []byte {
	pkt := NewPacket(1, HeartBeat, nil)
	return pkt.Build()
}
>>>>>>> e20f45e8c9dc9dc6e202d459cd56a928afdd5f95
