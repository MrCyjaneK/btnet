package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"git.mrcyjanek.net/mrcyjanek/btnet/utils"
	"github.com/sqweek/dialog"
)

var BTnet utils.BTnetJson

func main() {
	dialog.Message("%s", "Please choose directory with files of the site.").Info()
	site_data, err := dialog.Directory().Title("Input Directory").Browse()
	if err != nil {
		log.Fatal(err)
	}
	dialog.Message("%s", "Please choose key.priv.pgp.asc").Info()
	pgp_priv, err := dialog.File().Title("Input Directory").Load()
	if err != nil {
		log.Fatal(err)
	}
	dialog.Message("%s", "Please choose key.pgp.asc").Info()
	pgp_pub, err := dialog.File().Title("Input Directory").Load()
	if err != nil {
		log.Fatal(err)
	}
	BTnet.Build = time.Now().String()
	err = os.RemoveAll("_btnet")
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return err
			}
			if !info.IsDir() {
				fmt.Println(path, info.Size())
				data, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatal(err)
				}
				hasher := sha512.New()
				hasher.Write(data)
				shasum := fmt.Sprintf("%x", hasher.Sum(nil))
				log.Println(shasum)
				BTnet.Files = append(BTnet.Files, utils.Files{
					Path:   path,
					Sha512: shasum,
				})
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	os.Chdir(site_data)
	pgpasc, err := ioutil.ReadFile(pgp_pub)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	btnetjson, err := json.MarshalIndent(BTnet, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir("_btnet", 0750)
	ioutil.WriteFile("_btnet/btnet.json", btnetjson, 0750)
	err = ioutil.WriteFile("_btnet/pgp.asc", pgpasc, 0750)
	if err != nil {
		log.Fatal(err)
	}
	utils.SignFile(pgp_pub, pgp_priv, "_btnet/btnet.json", "_btnet/btnet.json.asc")
}
