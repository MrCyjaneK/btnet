package main

import (
	"git.mrcyjanek.net/mrcyjanek/btnet/btclient"
	"git.mrcyjanek.net/mrcyjanek/btnet/proxy"
)

func main() {
	btclient.Init()
	proxy.Init()
}
