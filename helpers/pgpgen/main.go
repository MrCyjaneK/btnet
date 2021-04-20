package main

import (
	"log"
	"os"

	"git.mrcyjanek.net/mrcyjanek/btnet/utils"
	g "github.com/AllenDang/giu"
	"github.com/sqweek/dialog"
)

var note = `Use this simple program to generate PGP key.

Please note, that for performace reasons you should use one PGP key per site, because PGP key is used to find newer version of website, so using one key would cause some mess.
But, we can't tell you to not do that, it's an open network :)`
var status = "waiting"

var outdir string
var name string
var exec bool

func genpgp() {
	go func() {
		status = "generating"
		utils.GenerateKeys(outdir, name)
		status = "Done!"
	}()
}

func pickdir() {
	go func() {
		if exec == true {
			return
		}
		exec = true
		out, err := dialog.Directory().Title("Output Directory").Browse()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(out)
		exec = true
	}()
}

func loop() {
	g.SingleWindow("PGP generator").Layout(
		g.Label(note).Wrapped(true),
		g.InputText("Name", &name),
		g.Line(
			g.InputText("Save Path", &outdir),
			g.Button("[O]").OnClick(pickdir),
		),
		g.Button("Generate!").OnClick(genpgp),
		g.Label(status),
	)
}

func main() {
	var err error
	outdir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wnd := g.NewMasterWindow("PGP generator", 400, 200, 0, nil)
	wnd.Run(loop)
}
