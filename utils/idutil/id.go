package idutil

import (
	"hash/crc32"
	"math/rand"

	"strconv"
	"time"

	"github.com/go-leo/sonyflake"
	"github.com/google/uuid"
)

const (
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"
)

var sf *sonyflake.Sonyflake

func init() {
	rand.Seed(time.Now().Unix())
	st := sonyflake.Settings{
		StartTime: time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC),
		Sequence:  func() uint16 { return uint16(rand.Intn(1024)) },
	}
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func GetSnowflakeID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	return id
}

func GetUUID() string {
	return uuid.NewString()
}

func GetRefreshToken() string {
	return randString(Alphabet62, 50)
}

func randString(letters string, n int) string {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	if _, err := rand.Read(randomness); err != nil {
		panic(err)
	}

	l := len(letters)
	// fill output
	for pos := range output {
		// get random item
		random := randomness[pos]

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letters[randomPos]
	}

	return string(output)
}

// String hashes a string to a unique hashcode.
// https://github.com/hashicorp/terraform/blob/master/helper/hashcode/hashcode.go
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
func HashCode(s string) int {
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

func HashCodeUInt64(id uint64) int {
	s := strconv.FormatUint(id, 10)
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
