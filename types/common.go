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
