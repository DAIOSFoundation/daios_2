package types

type Address [15]byte
type Hash string

func NewAddress(s string) (address Address) {
	a := new(Address)

	copy(a[:], s)

	return *a
}

func ConvertAddress(s []byte) (address Address) {
	a := new(Address)

	copy(a[:], s)

	return *a
}

const k = `/key/swarm/psk/1.0.0/
/base16/
289bf3546b030a5d17e32d695c755aad8471034738b879645ff238cabeb15266`
