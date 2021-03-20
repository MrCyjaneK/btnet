package btclient

import (
	"log"

	"github.com/anacrolix/torrent"
)

func handleNewPeer(peer *torrent.Peer) {
	log.Println("Got new peer:")
	log.Println(" - Net:", peer.Network)
	log.Println(" - Name:", peer.PeerClientName)
	log.Println(" - String:", peer.String())
}
