package uid

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"os"
	"sync"
	"time"
	"strconv"
)

const epochStart = 122192928000000000
const saltStr = "thecatjumpingdownthebedshedsomefur"

var (
	initOnce   	  sync.Once
	epochFunc     = func() uint64 { return epochStart + uint64(time.Now().UnixNano()/100) }
	hardwareAddr  [6]byte
	posixUID      = os.Getuid()
	posixGID      = os.Getgid()
	shaGen		  = sha256.New()
	salt 		  = []byte(saltStr)
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

	pId := []byte(strconv.Itoa(posixUID))
	gId := []byte(strconv.Itoa(posixGID))
	time := []byte(strconv.FormatUint(epochFunc(), 10))

	hashing := append( append( append( append( append( salt, pId... ), gId... ), hardwareAddr[:]... ), time... ), []byte(name)... )
	u := shaGen.Sum(hashing)

	return hex.EncodeToString(u)
}
