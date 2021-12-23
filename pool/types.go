package pool

type relayCounter struct {
	prevNum uint64
	num     uint64
	ts      uint64
}

type node struct {
	userId    string
	chordId   string
	pubKey    string
	version   string
	syncState map[uint64]string

	relayNum map[uint64]uint64
	uptime   map[uint64]uint64
}
