package btclient

import (
	"log"
	"os"

	"github.com/anacrolix/torrent"
)

var cl *torrent.Client

// Init - Load all the things that are needed to work properly
func Init() {
	var err error
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	config := torrent.NewDefaultClientConfig()
	config.DataDir = dir + "/.BTnet"
	config.Seed = true
	config.HTTPUserAgent = "BTnet v0.0.0 (websites over torrent)"
	config.ExtendedHandshakeClientVersion = "btnet dev 20210311"
	config.Bep20 = "-BN0000-"
	config.Callbacks.NewPeer = append(config.Callbacks.NewPeer, handleNewPeer)

	cl, err = torrent.NewClient(config)
	if err != nil {
		log.Panic(err)
	}
}
