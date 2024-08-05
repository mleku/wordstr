package wordstr

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"

	"github.com/Hubmakerlabs/replicatr/pkg/ec/secp256k1"
	"github.com/Hubmakerlabs/replicatr/pkg/nostr/bech32encoding"
	"github.com/minio/sha256-simd"
	"wordstr.mleku.dev/wordlists"
)

func GetPlaces() (places []*big.Int) {
	p := big.NewInt(1)
	for _ = range 24 {
		p1 := big.NewInt(0)
		p1.SetBytes(p.Bytes())
		places = append([]*big.Int{p1}, places...)
		p.Mul(p, big.NewInt(2048))
	}
	return
}

func FromNsec(nsec string) (words string, err error) {
	var sk []byte
	if len(nsec) == 2*secp256k1.SecKeyBytesLen {
		if sk, err = hex.DecodeString(nsec); err != nil {
			fmt.Fprintf(os.Stderr, "failed to decode nsec from hex form: %s\n", err)
			os.Exit(1)
		}
	} else {
		var prf string
		var val any
		if prf, val, err = bech32encoding.Decode(nsec); err != nil {
			err = fmt.Errorf("%s", err)
			fmt.Println(1, err)
			return
		}
		if prf != bech32encoding.NsecHRP {
			err = fmt.Errorf("nostr nsec must start with %s, not %s\n", bech32encoding.NsecHRP, prf)
			fmt.Println(2, err)
			return
		}
		if sk, err = hex.DecodeString(val.(string)); err != nil {
			err = fmt.Errorf("failed to decode nsec from hex form: %s\n", err)
			fmt.Println(3, err)
			return
		}
	}
	// left-pad in case the secret is smaller than 32 bytes long
	if len(sk) != secp256k1.SecKeyBytesLen {
		sk = append(make([]byte, secp256k1.SecKeyBytesLen+1-len(sk)), sk...)
	}
	h := sha256.Sum256(sk)
	// add 7 bits of check with the MSB forced to 1 (128)
	sec := append([]byte{h[0] | 128}, sk...)
	div, mod := big.NewInt(0), big.NewInt(0)
	div.SetBytes(sec)
	var ww []string
	for div.Cmp(big.NewInt(0)) > 0 {
		div.DivMod(div, big.NewInt(2048), mod)
		w := wordlists.English[mod.Int64()]
		ww = append([]string{w}, ww...)
	}
	words = strings.Join(ww, " ")
	return
}

func ToNsec(places []*big.Int, words []string) (hsec, nsec string, err error) {
	if len(words) != 24 {
		err = fmt.Errorf("not enough words in key, require 24, got %d", len(os.Args[2:]))
		return
	}
	var indexes []int64
	for _, word := range words {
		var found bool
		for i := range wordlists.English {
			if wordlists.English[i] == word {
				indexes = append(indexes, int64(i))
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("word not found in wordlist: %s", word)
			return
		}
	}
	a := math.MaxUint32
	_ = a
	key := big.NewInt(0)
	for i := range indexes {
		places[i].Mul(places[i], big.NewInt(indexes[i]))
	}
	for i := range places {
		key = key.Add(key, places[i])
	}
	skb := key.Bytes()
	if len(skb) < secp256k1.SecKeyBytesLen+1 {
		skb = append(make([]byte, secp256k1.SecKeyBytesLen+1-len(skb)), skb...)
	}
	k, check := skb[1:], skb[0]&127
	h := sha256.Sum256(k)
	// force the MSB to 1
	actual := h[0] & 127
	if check != actual {
		err = fmt.Errorf("key parity check failed, got %d, should have been %d - if key is intact, here it is: %0x",
			actual, check, k)
		return
	}
	hsec = hex.EncodeToString(k)
	if nsec, err = bech32encoding.HexToNsec(hsec); err != nil {
		err = fmt.Errorf("failed to encode nsec from hex form: %s", err)
		return
	}
	return
}
