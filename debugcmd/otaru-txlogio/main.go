package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nyaxt/otaru/btncrypt"
	"github.com/nyaxt/otaru/facade"
	"github.com/nyaxt/otaru/gcloud/auth"
	"github.com/nyaxt/otaru/gcloud/datastore"
)

var (
	flagConfigDir = flag.String("configDir", facade.DefaultConfigDir(), "Config dirpath")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s [options] {purge}\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flag.Usage = Usage
	flag.Parse()

	cfg, err := facade.NewConfig(*flagConfigDir)
	if err != nil {
		log.Printf("%v", err)
		Usage()
		os.Exit(2)
	}
	if flag.NArg() != 1 {
		Usage()
		os.Exit(2)
	}
	switch flag.Arg(0) {
	case "purge":
		break
	default:
		log.Printf("Unknown cmd: %v", flag.Arg(0))
		os.Exit(1)
	}

	clisrc, err := auth.GetGCloudClientSource(cfg.CredentialsFilePath, cfg.TokenCacheFilePath, false)
	if err != nil {
		log.Fatalf("Failed to init GCloudClientSource: %v", err)
	}

	key := btncrypt.KeyFromPassword(cfg.Password)
	c, err := btncrypt.NewCipher(key)
	if err != nil {
		log.Fatalf("Failed to init btncrypt.Cipher: %v", err)
	}
	dscfg := datastore.NewConfig(cfg.ProjectName, cfg.BucketName, c, clisrc)

	txlogio := datastore.NewDBTransactionLogIO(dscfg)

	switch flag.Arg(0) {
	case "purge":
		if err := txlogio.DeleteAllTransactions(); err != nil {
			log.Printf("DeleteAllTransactions() failed: %v", err)
		}
		break
	default:
		log.Printf("Unknown cmd: %v", flag.Arg(0))
		os.Exit(1)
	}
}