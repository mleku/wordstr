# wordstr

`wordstr` is a tool to convert a nostr bech32 encoded `nsec` or hexadecimal
encoded secret key into a word mnemonic that can be stored in hard copy.

It uses the standard BIP-39 word list to represent 11 bit value elements of
the secret key but instead of hashing it to derive the actual secret key or
extended secret key as with bitcoin keys, it is the exact bits of the secret 
itself. (encrypting it would be a separate protocol to keep things simple).

usage:

to generate a word key:

```bash
wordstr from <hex/nsec nostr secret key>
```

to convert a word key back to hex and bech32 `nsec` format:

```bash
wordstr to <27 word mnemonic key>
```

the output of the `to` command is formatted as such for easy use in scripts
or invocation from other programs:

```
HSEC="<hex secret key>"
NSEC="<bech32 nsec secret key>"
```

## building/installing

`wordstr` can be run directly from a system with a configured Go
installation as follows:

```
go run wordstr.mleku.dev@latest <parameters>
```

or if your `GOBIN` refers to a location also present in your `PATH`, installed:

```bash
go install wordstr.mleku.dev@latest
```

## technical details of protocol

The implementation found in this repository uses big integers to derive the 
individual ciphers of the word key, and from them in reverse using a place 
table to reverse it. Possibly a more complex bit-rotating method could be 
implemented but it is a longer algorithm, and unnecessary for any 
non-embedded hardware. For embedded hardware, hexadecimal is preferable 
anyway. 

This encoding is to provide hard copy cold storage for users' nostr 
secret keys, as a method of entering the secret key from such cold storage 
into a signing device.

24 words of a 2048 word dictionary from BIP-34 produces 11 bits per word, 
and the nearest common length between these two is 264 bits, or 33 bytes, 
giving us 1 byte for check purposes.

The first bit must always be 1, to simplify the derivation of the ciphers of 
the mnemonic key, so we derive a check with the first byte of the SHA256 
hash of the 33 byte key, with the most significant bit set to 1.

- when generated, first byte of the hash of the key, bitwise OR 128 so the 
  MSB is always 1, and the remainder match the hash of the correct key

- when checked, the value is masked using bitwise AND 127, and checked 
  against the 7 least significant bits of the first byte of the hash of the 32 
  remaining bytes, so, the same operation as when generating, on the first 
  byte of the hash, and then compared to the first byte itself

This provides a protection against a 1/128 chance of a bit flip producing an 
also valid key versus the 7 bits of check, it would be nice if it was 
stronger but it will rarely flag as wrong without actually being wrong, and 
the 11 bit ciphers of this encoding mean more than 50% more words, and a 
much greater chance of error in transcribing in total, and a reduction in 
complexity of encoding by avoiding the need for another word and fiddly bit 
checking beyond this one OR masking on one byte.

For this reason, the error will return the decoded key from a word key that 
fails in case the bit flip is in that first word but not the key itself, 
this is unlikely, however, and deriving the npub will confirm this by the 
mismatch. This could be validated by fetching the relevant user metadata 
kind 0 event to show the user in such a case, and enable providing a 
correction to the word key. It is unlikely to recover the key but 1/128 of 
cases it could.