package filemanager

import (
	"bytes"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/arvan/clientsimulator/config"
)

type Teacher struct {
	UUID     string `json:"uuid,omitempty"`
	UserName string `json:"username,omitempty"`
	Quote    string `json:"qoute"`
}

func SendQouteToServer(client Clients, cfg *config.Config, index int) {
	defer Wg.Done()
	var body []byte
	var readErr error
	var resp *http.Response

	for i := 0; i < cfg.Clients.NumberOfRequest; i++ {
		r := rand.Intn(cfg.Clients.Dilay)
		time.Sleep(time.Duration(r) * time.Millisecond)

		var err error
		var jsonStr []byte
		if cfg.Clients.PostDuplicateQoute {
			jsonStr = []byte("{ \"qoute\":\"" + client.Qoute + "\"}")
		} else {
			jsonStr = []byte("{ \"qoute\":\"" + client.Qoute + strconv.Itoa(i) + "\"}")
		}

		req, _ := http.NewRequest("POST", cfg.Server.Addr+"/api/v1/qoute/add/"+client.Uuid+"/"+client.Username, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err = client.Do(req)

		if err == nil {
			Zlogger.Info().Msgf("Request Response :   %+v", resp)
			body, readErr = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			client.CloseIdleConnections()
			if readErr == nil {
				Zlogger.Info().Msgf("Body Response :   %+v", body)
			}
		} else {
			Zlogger.Error().Msgf("Server is not responding")
		}

	}
}
