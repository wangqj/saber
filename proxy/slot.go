package proxy

import "hash/crc32"

type Slot struct {
	id     int
	node   *Node
	status int
}

func FindSlot(key string) (s *Slot) {

	return nil
}
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}