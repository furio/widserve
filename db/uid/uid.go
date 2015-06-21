package uid

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"sync"
	"time"
)

const epochStart = 122192928000000000
const salt = []byte("thecatjumpingdownthebedshedsomefur")

var (
	initOnce   	  sync.Once
	epochFunc     = func() int64 { return epochStart + uint64(time.Now().UnixNano()/100) }
	hardwareAddr  [6]byte
	posixUID      = uint32(os.Getuid())
	posixGID      = uint32(os.Getgid())
	shaGen		  = sha256.New()
)

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

func initHardwareAddr() {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				copy(hardwareAddr[:], iface.HardwareAddr)
				return
			}
		}
	}

	safeRandom(hardwareAddr[:])

	hardwareAddr[0] |= 0x01
}

func NewUid(name string) string {
	initOnce.Do( func () {
		initHardwareAddr()
	})

	u := shaGen.Sum(append(salt, []byte(posixUID), []byte(posixGID), hardwareAddr, []byte(epochFunc()), []byte(name)))

	return hex.EncodeToString(u)
}
