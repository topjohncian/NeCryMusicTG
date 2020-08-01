package main

import (
	"encoding/json"
	"fmt"
	httputils "github.com/topjohncian/NeCryMusicTG/pkg/http"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"time"
)

var b *tb.Bot

type Sentence struct {
	UID  int    `json:"uid"`
	Type int    `json:"type"`
	From string `json:"from"`
	Text string `json:"text"`
}

func main() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:7890", nil, proxy.Direct)
	if err != nil {
		log.Fatal("can't connect to the proxy:", err)
	}

	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial

	b, err = tb.NewBot(tb.Settings{
		Token:  "1087503176:AAHXJrN0ZUw5GxeRjhmpC5jXd7DdrwnWxbc",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: httpClient,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/cry", func(m *tb.Message) {
		res, err := httputils.Get("http://api.heerdev.top/nemusic/random")
		if err != nil {
			log.Println("请求语录时发生错误：", err)
			b.Reply(m, "服务器去火星了，等会儿再试试吧")
		}

		sentence := &Sentence{}
		err = json.Unmarshal(res, sentence)
		if err != nil {
			log.Println("请求语录时发生错误：", err)
			b.Reply(m, "服务器去火星了，等会儿再试试吧")
		}
		b.Reply(m, fmt.Sprintf("\"%s\" —— %s", sentence.Text, sentence.From))
	})

	b.Start()
}
