package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arvan/clientsimulator/config"
	"github.com/arvan/clientsimulator/pkg/filemanager"
)

var (
	configFile string
	cfg        *config.Config
)

func main() {
	flag.StringVar(&configFile, "c", "config.yml", "config file")
	if !config.Load(configFile) {
		log.Fatal()
	}

	cfg = config.Get()
	fmt.Println(cfg.Clients.FileName)
	fmt.Println(cfg.Clients.NumberOfClients)
	fmt.Println(cfg.Clients.ResetUUID)

	if cfg.Clients.ResetUUID == true {
		e := os.Remove(cfg.Clients.FileName)
		if e != nil {
			log.Fatal(e)
		}
	}

	filemanager.LoadFile(cfg)
	filemanager.RunBroadcast()
	time.Sleep(10 * time.Second)
	fmt.Println("Terminating Program")
}
