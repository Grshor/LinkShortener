package pkg

import (
	"hash/fnv"
	"strconv"
)

func HashFnv(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.FormatUint(uint64(h.Sum32()), 10)
}
