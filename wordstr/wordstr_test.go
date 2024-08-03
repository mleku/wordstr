package wordstr_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"ec.mleku.dev/v2/secp256k1"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/bech32encoding"
	"wordstr.mleku.dev/wordstr"
)

func TestWordstr(t *testing.T) {
	var err error
	var sec *secp256k1.SecretKey
	var wordsx, wordsb, nsec, hsec, nsec2 string
	for _ = range 1000 {
		if sec, err = secp256k1.GenerateSecretKey(); err != nil {
			t.Fatal(err)
		}
		skb := sec.Serialize()
		skh := hex.EncodeToString(skb)
		// hex
		if wordsx, err = wordstr.FromNsec(skh); err != nil {
			t.Fatal(err)
		}
		// bech32
		if nsec, err = bech32encoding.HexToNsec(skh); err != nil {
			t.Fatal(err)
		}
		if wordsb, err = wordstr.FromNsec(nsec); err != nil {
			t.Fatal(err)
		}
		if wordsx != wordsb {
			t.Fatalf("words do not match hex %s vs nsec %s\n%s\n!=\n%s", skh, nsec, wordsx, wordsb)
		}
		split := strings.Split(wordsx, " ")
		places := wordstr.GetPlaces()
		if hsec, nsec2, err = wordstr.ToNsec(places, split); err != nil {
			t.Fatal(err)
		}
		if nsec2 != nsec {
			t.Fatalf("did not recover same nsec\n%s\n!=\n%s", nsec, nsec2)
		}
		if hsec != skh {
			t.Fatalf("did not recover same hsec\n%s\n!=\n%s", skh, hsec)
		}
	}
}
