package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Debug      bool
	Token      string
	WebAppAddr string
	DBAddr     string
	QRText     string
	BotID      int64
}

func (c Config) Print() {
	msg := `{
	"debug": "%t",
	"token": "***",
	"web app url": "%s",
	"db addr": "%s",
	"qr text": "%s",
	"bot id": "%d",
}`
	log.Printf(msg, c.Debug, c.WebAppAddr, c.DBAddr, c.QRText, c.BotID)
}

func New() *Config {
	// if need additional logging
	debug := false
	debugStr := strings.ToLower(os.Getenv("REQ_LOG"))
	if debugStr == "true" {
		debug = true
	}

	botIDStr := must("BOT_ID")

	id, err := strconv.ParseInt(botIDStr, 10, 64)
	if err != nil || id == 0 {
		log.Fatalf("config: bot id: parse: %s", botIDStr)
	}

	return &Config{
		Debug:      debug,
		Token:      must("BOT_TOKEN"),
		WebAppAddr: must("WEB_URL"),
		DBAddr:     must("DATABASE_URL"),
		QRText:     must("QR_TEXT"),
		BotID:      id,
	}
}
