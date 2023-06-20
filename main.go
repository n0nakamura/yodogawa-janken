package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

const (
	name    = "yodogawa-janken"
	version = "0.3.0"
)

type Config struct {
	Nsec string `json:"nsec"`
}

var defaultRelays = []string{
	// "ws://172.17.0.1:7447",

	"wss://relay-jp.nostr.wirednet.jp",
	"wss://relay.nostr.wirednet.jp",
	"wss://nostr-relay.nokotaro.com",
	"wss://nostr.holybea.com",
	"wss://nostr.h3z.jp",

	// "wss://nos.lol",
	// "wss://nostr.mom",
	// "wss://nostr.oxtr.dev",
	// "wss://relay.nostr.bg",
	// "wss://nostr.bitcoiner.social",
	// "wss://relay.damus.io",
	// "wss://no.str.cr",
}

// TODO: エラーハンドリングをきちんと書く
// TODO: Contextを理解する

// TODO: string型の戻り値が無いことを示すのに "" でよいかどうか調べる
// TODO: そもそも関数calcHexを使う設計でよいかどうか
func calcHex(nsec string) (string, string, error) {
	var sk, pub string

	if _, s, err := nip19.Decode(nsec); err != nil {
		return "", "", err
	} else {
		sk = s.(string)
	}

	if p, err := nostr.GetPublicKey(sk); err == nil {
		if _, err := nip19.EncodePublicKey(pub); err == nil {
			pub = p
		} else {
			return "", "", err
		}
	} else {
		return "", "", err
	}

	return sk, pub, nil
}

func subscribeEvent(sk string, pub string, pevc chan *nostr.Event) error {
	ctx := context.Background()

	// TODO: 複数のリレーから取得できるようにする。
	relay, err := nostr.RelayConnect(ctx, defaultRelays[0])
	if err != nil {
		log.Fatal(err)
	}

	var filters nostr.Filters
	since := nostr.Now()
	filters = []nostr.Filter{{
		Kinds: []int{nostr.KindTextNote},
		Tags: nostr.TagMap{
			"p": {pub},
		},
		Since: &since,
	}}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		log.Fatal(err)
	}

	for pev := range sub.Events {
		pevc <- pev
	}

	return nil
}

func postReply(sk string, pub string, pevc chan *nostr.Event) error {
	pev := <-pevc

	var content string
	if c, err := generateContent(pev.Content); err == nil {
		content = c
	} else if err == ErrNoValuesIncluded {
		return nil
	} else {
		return err
	}

	// content = "Your hand: " + handNames[playerHand] + "\n" +
	// 	outcomeNameMap[result]

	// Create a event
	ev := nostr.Event{}
	ev.PubKey = pub
	ev.CreatedAt = nostr.Now()
	ev.Kind = nostr.KindTextNote // kind1
	ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"p", pev.PubKey})
	ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"e", pev.ID, "", ""})
	ev.Content = content
	ev.Sign(sk)

	// Post the event
	success := 0
	for _, url := range defaultRelays {
		relay, err := nostr.RelayConnect(context.Background(), url)
		if err != nil {
			continue
		}
		status, err := relay.Publish(context.Background(), ev)
		relay.Close()
		if err == nil && status != nostr.PublishStatusFailed {
			log.Println("published to", url, status)
			success++
		}
	}
	// Even if a relay connection fails, it does not immediately trigger an error.
	if success == 0 {
		return errors.New("failed to publish")
	}

	return nil
}

func main() {
	fmt.Println(name, version)

	var cfg Config
	if file, err := os.ReadFile("config.json"); err == nil {
		if err := json.Unmarshal(file, &cfg); err != nil {
			log.Fatal(err)
			return
		}
	} else {
		log.Fatal(err)
		return
	}

	sk, pub, err := calcHex(cfg.Nsec)
	if err != nil {
		log.Fatal(err)
		return
	}

	pevc := make(chan *nostr.Event)
	go func() {
		for {
			if err := postReply(sk, pub, pevc); err != nil {
				continue
			}
		}
	}()

	err = subscribeEvent(sk, pub, pevc)
	if err != nil {
		log.Fatal(err)
		return
	}
}
