package main

import (
	"fmt"
	"log"

	lam "github.com/fn-code/lamcctv/internal/cctv"
	"github.com/fn-code/lamcctv/pkg/bot"
	uuid "github.com/satori/go.uuid"
)

func main() {
	var urls = []string{
		"http://192.168.38.134/",
		"http://192.168.38.135/",
		"http://192.168.38.136/",
		"http://192.168.38.137/",
		"http://192.168.38.138/",
		"http://192.168.38.139/",
		"http://192.168.38.140/",
		"http://192.168.38.141/",
		"http://192.168.38.142/",
		"http://192.168.38.143/",
		"http://192.168.38.144/",
		"http://192.168.38.145/",
		"http://192.168.38.146/",
		"http://192.168.38.147/",
		"http://192.168.38.148/",
		"http://192.168.38.149/",
		"http://192.168.38.150/",
		"http://192.168.38.151/",
		"http://192.168.38.152/",
		"http://192.168.38.153/",
		"http://192.168.38.154/",
		"http://192.168.38.155/",
		"http://192.168.38.156/",
		"http://192.168.38.157/",
		"http://192.168.38.158/",
		"http://192.168.38.159/",
		"http://192.168.38.160/",
		"http://192.168.38.161/",
	}

	ctv := make([]*lam.CCTV, 0)
	for _, v := range urls {
		cctv := &lam.CCTV{
			ID:   uuid.NewV4().String(),
			Name: fmt.Sprintf("IP %s", v),
			URL:  v,
		}
		ctv = append(ctv, cctv)
	}

	// Group ID -330310831
	bs, err := bot.New("TELEGRAM_APIKEY", -330310831)
	if err != nil {
		log.Println(err)
	}

	set := lam.New(ctv, bs)
	set.ProsesCCTV()

}
