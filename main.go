package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
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

// 関数doから関数subscribeへ処理を切り分ける
func do(sk string, pub string) error {
	ctx := context.Background()

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

	// TODO: このfor文内部の処理をmain関数に移しゴルーチンとチャンネルでいい感じにする
	for pev := range sub.Events {
		var inputHand string
		if re, err := regexp.Compile(`[RSP✊👊🤛🤜💪✌🤞🖐✋🤚🖖🫲🫱🫳🫴👋👐🤲🤗🤟🤝👌🤌🤏🤘🤙👈👉👆👇☝👍👏🙏]`); err != nil {
			return err
		} else {
			inputHand = re.FindString(pev.Content)
		}

		if inputHand == "" {
			continue
		}

		var content string
		playerHand := getPlayerHand(inputHand)
		yodogawaHand := biasJanken()
		result := doJanken(playerHand, yodogawaHand)
		if result == WIN || result == LOSE || result == DRAW {
			content = "Your hand: " + handNames[playerHand] + "\n" +
				"Yodogawa-san hand: " + handNames[yodogawaHand] + "\n" +
				outcomeNameMap[result]
		} else {
			content = "Your hand: " + handNames[playerHand] + "\n" +
				outcomeNameMap[result]
		}

		if err = postReply(sk, pub, pev, content); err != nil {
			return err
		}
	}

	return nil
}

func postReply(sk string, pub string, pev *nostr.Event, content string) error {
	ev := nostr.Event{}

	// Create a event
	ev.PubKey = pub
	ev.CreatedAt = nostr.Now()
	ev.Kind = nostr.KindTextNote // kind1
	ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"p", pev.PubKey})
	ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"e", pev.ID, "", "reply"})
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

	err = do(sk, pub)
	if err != nil {
		log.Fatal(err)
	}
}
