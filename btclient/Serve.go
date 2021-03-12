package btclient

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"git.mrcyjanek.net/mrcyjanek/btnet/utils"
)

// Serve - simple way to serve whatever is inside the magnet.
func Serve(wr http.ResponseWriter, host string, path string) {
	magnet := utils.GetMagnetFromBTnet(host)
	T, err := cl.AddMagnet(magnet)
	//wr.WriteHeader(http.StatusProcessing)
	if err != nil {
		wr.WriteHeader(http.StatusInternalServerError)
		wr.Write([]byte(err.Error()))
		return
	}
	log.Println("Waiting to get torrent info.")
	<-T.GotInfo()
	log.Println("OK")
	//wr.WriteHeader(http.StatusOK)
	files := T.Files()
	for _, file := range files {
		torrentpath := strings.Replace(file.Path(), T.Name(), "", -1)
		if path == torrentpath || path+"index.html" == torrentpath {
			// Download file and show it
			file.Download()
			srcReader := file.NewReader()
			defer srcReader.Close()
			//2021/03/12 11:43:43 Done! Written total: 16384 423036 423036
			log.Println("Downloading", file.Path())
			bs := make([]byte, file.FileInfo().Length)
			srcReader.SetReadahead(file.FileInfo().Length)
			written, err := srcReader.Read(bs)
			if int64(written) != file.FileInfo().Length {
				log.Println("[INFO] File not yet ready, calling loop in 5 seconds", int64(written), "!=", file.FileInfo().Length)
				log.Println("[INFO] \\___", file.FileInfo().Length-int64(written), "left")
				time.Sleep(time.Second * 5)
				Serve(wr, host, path)
				return
			}
			if err != nil && err.Error() != "EOF" {
				log.Println(err)
				wr.Write([]byte(err.Error()))
				return
			}
			ct := http.DetectContentType(bs)
			if strings.Split(ct, "; ")[0] == "text/plain" {
				if strings.HasSuffix(file.Path(), ".css") {
					ct = "text/css"
				} else if strings.HasSuffix(file.Path(), ".js") {
					ct = "text/javascript"
				}
			}
			wr.Header().Set("Content-Type", ct)
			copied, err := io.Copy(wr, bytes.NewReader(bs))
			if err != nil {
				wr.Write([]byte(err.Error()))
			}
			log.Println("Done! Written total:", written, copied, file.FileInfo().Length)
			//log.Println(file.DisplayPath())
			return
		}
	}
	wr.Write([]byte(utils.MakeNiceList(T)))
}
