package utils

import (
	"log"
	"strings"

	"github.com/anacrolix/torrent"
)

// IsValidMagnetBTnet - Tells is asd.btnet is a valid magnet url
func IsValidMagnetBTnet(url string) (bool, string) {
	parts := strings.Split(url, "/")
	host := parts[0]
	parts = strings.Split(host, ".")
	if len(parts) < 2 {
		return false, "len(parts) < 2"
	}
	if parts[1] != "btnet" {
		return false, "domain not end with .btnet"
	}

	_, err := torrent.TorrentSpecFromMagnetUri(GetMagnetURI(parts[0]))
	if err != nil {
		log.Println(err)
		return false, err.Error()
	}
	return true, "OK"
}

// GetMagnetFromBTnet - Get magnet directly from url.
func GetMagnetFromBTnet(url string) string {
	parts := strings.Split(url, "/")
	host := parts[0]
	parts = strings.Split(host, ".")
	if len(parts) < 2 {
		return "len(parts)_<_2"
	}
	if parts[1] != "btnet" {
		return "domain_not_end_with_.btnet"
	}

	return GetMagnetURI(parts[0])
}

// GetMagnetURI - Generate magnet uri.
func GetMagnetURI(uri string) string {
	magnet := "magnet:?xt=urn:btih:" + uri
	// TEMP: Add tracker
	magnet = magnet + "&tr=udp%3a%2f%2ftracker.opentrackr.org%3a1337%2fannounce"
	return magnet
}
