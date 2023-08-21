package filemanager

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/arvan/clientsimulator/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	configFile string
	Wg         sync.WaitGroup
	cfg        *config.Config
	clientList []Clients
	Zlogger    = log.With().Str("service", "Arvan-Qoute").Logger()
)

func LoadFile(cfgPtr *config.Config) {
	cfg = cfgPtr

	file, err := os.OpenFile(cfg.Clients.FileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		Zlogger.Error().Msg(err.Error())
	}

	defer file.Close()

	LineCount, errLC := lineCounter(file)
	if errLC != nil {
		Zlogger.Error().Msg(errLC.Error())
	}

	if LineCount <= cfg.Clients.NumberOfClients {
		for i := 0; i < cfg.Clients.NumberOfClients-LineCount; i++ {
			u, _ := uuid.NewRandom()
			client := &User{
				UUID:     u.String(),
				UserName: base32.StdEncoding.EncodeToString(u.NodeID())}
			clientBytr, _ := json.Marshal(client)
			file.WriteString(string(clientBytr) + "\n")
		}
	}
}

func RunBroadcast() {
	index := 1

	file, err := os.OpenFile(cfg.Clients.FileName, os.O_RDONLY, 0777)
	if err != nil {
		Zlogger.Error().Msg(err.Error())
	}
	defer file.Close()
	data := make([]byte, 78)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}
		newdata := Clients{"", "", ""}
		err = json.Unmarshal([]byte(string(data[:n])), &newdata)
		if err != nil {
			Zlogger.Error().Msg(err.Error())
			return
		}
		Wg.Add(1)
		newdata.Qoute = cfg.Clients.Qoute
		go SendQouteToServer(newdata, cfg, index)
		index++
		clientList = append(clientList, newdata)
	}
	Wg.Wait()
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
