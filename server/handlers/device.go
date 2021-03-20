package handlers

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v4"
	"miikka.xyz/rs/config"
	"miikka.xyz/rs/store"
)

// These will be used in device, when reading responses
const okMessage = "ok"
const errorMessage = "error"
const invalidRequest = "hahaa"

// TagToUsername changes RFID-tag to username
func TagToUsername(db *store.DB, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Take stuff after route
		tag := r.URL.Path[len(conf.DeviceTagRoute):]
		log.Info(conf.DeviceTagRoute)
		log.Info("request with tag: ", tag)

		// Check header. This header is set in device
		token := r.Header.Get("Arduino-Token")
		if token != conf.DeviceToken && conf.IsProduction {
			log.Error("invalid token ", token)
			w.Write([]byte(invalidRequest))
			return
		}

		username, err := db.GetUsernameByTag(tag)
		if err != nil {
			// Username was not found
			if err == pgx.ErrNoRows {
				log.Info("unknown tag")
				// Save unknown tag so it can be registered
				err = db.SaveAudit(config.EventUnknownTag, tag, "")
				if err != nil {
					log.Error("audit trail failed ", err)
				}
			} else {
				log.Error(err)
			}
			w.Write([]byte(tag))
			return
		}

		log.Info("username ", *username, " found")
		w.Write([]byte(*username))
	}
}

// SaveGame saves game
func SaveGame(db *store.DB, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Arduino-Token")
		// Only Arduino is allowed
		if token != conf.DeviceToken && conf.IsProduction {
			log.Error("invalid token ", token)
			w.Write([]byte(invalidRequest))
			return
		}

		agent := r.Header.Get("User-Agent")
		// Only Arduino is allowed
		if agent != "ArduinoWiFi/1.1" && conf.IsProduction {
			log.Error("client wasn't arduino")
			w.Write([]byte(invalidRequest))
			return
		}

		payload := r.URL.Path[len(conf.DeviceSaveRoute):]
		log.Info("payload: ", payload)

		payload = strings.ReplaceAll(payload, " ", "")
		players := strings.SplitN(payload, ",", 2)

		if len(players) < 2 {
			log.Error("not enough players")
			w.Write([]byte(errorMessage))
			return
		}

		if players[0] == "" || players[1] == "" {
			log.Error("empty value")
			w.Write([]byte(errorMessage))
			return
		}

		if players[0] == players[1] {
			log.Error("winner and loser can't be same player")
			w.Write([]byte(errorMessage))
			return
		}

		err := db.CreateGameIOT(players[0], players[1], conf.DeviceSport)
		if err != nil {
			log.Error(err)
			w.Write([]byte(errorMessage))
			return
		}

		// Needed for client to show right message
		w.Write([]byte(okMessage))
	}

}
