package main

// @param msg is the string to be hashed

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
)

// constants [§4.2.2]
var K [64]uint32 = [64]uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
	0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
	0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
	0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
	0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2}

// initial hash value [§5.3.3]
var H [8]uint32 = [8]uint32{
	0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19}

func Sha256(msg string) string {

	msgBSlice := preprocess(msg)
	hash := H

	for _, chunk := range msgBSlice {

		hash = Sha256_compress(chunk, hash)

	}

	return fmt.Sprintf("%08x%08x%08x%08x%08x%08x%08x%08x", hash[0], hash[1], hash[2], hash[3], hash[4], hash[5], hash[6], hash[7])

}

func Sha256_compress(chunk [16]uint32, iv [8]uint32) [8]uint32 {

	msgSchedule := createMessageSchedule(chunk)

	a := iv[0]
	b := iv[1]
	c := iv[2]
	d := iv[3]
	e := iv[4]
	f := iv[5]
	g := iv[6]
	h := iv[7]

	for t := 0; t < 64; t++ {

		T1 := h + Σ1(e) + Ch(e, f, g) + K[t] + msgSchedule[t]
		T2 := Σ0(a) + Maj(a, b, c)

		h = g
		g = f
		f = e
		e = (d + T1)
		d = c
		c = b
		b = a
		a = (T1 + T2)

	}

	iv[0] = (iv[0] + a)
	iv[1] = (iv[1] + b)
	iv[2] = (iv[2] + c)
	iv[3] = (iv[3] + d)
	iv[4] = (iv[4] + e)
	iv[5] = (iv[5] + f)
	iv[6] = (iv[6] + g)
	iv[7] = (iv[7] + h)

	return iv
}

func preprocess(msg string) [][16]uint32 {

	//Convert msg to bits
	msgBytes := stringToByteSlice(msg)

	//Get msg size
	msgSize := intToByteSlice(len(msgBytes)*8, 8)

	//Add trailing '1' bit
	msgBytes = append(msgBytes, 0x80)
	//Add traliing 0's
	for len(msgBytes)%64 != 56 {
		msgBytes = append(msgBytes, 0)
	}

	//Add msg Size
	msgBytes = append(msgBytes, msgSize...)

	msgUint32 := bytesToUint32(msgBytes)

	//Seperate into strings size 512
	msgUSlice := [][16]uint32{}
	for len(msgUint32) > 0 {
		var msgUint32_16 [16]uint32
		for i := 0; i < 16; i++ {
			msgUint32_16[i] = msgUint32[i]
		}

		msgUSlice = append(msgUSlice, msgUint32_16)
		msgUint32 = msgUint32[16:]
	}

	return msgUSlice
}

func stringToByteSlice(s string) []byte {

	val, err := hex.DecodeString(hex.EncodeToString([]byte(s)))
	if err != nil {
		log.Fatalln("Error Decoding string to hex slice", err)
	}

	return val
}

func intToByteSlice(num int, pad int) []byte {

	numH := strconv.FormatInt(int64(num), 16)
	if len(numH)%2 == 1 {
		numH = "0" + numH
	}
	msgBS, err := hex.DecodeString(numH)

	if err != nil {
		log.Fatalln("Error Decoding int to hex slice", err)
	}

	for i := len(msgBS); i < pad; i++ {
		msgBS = append([]byte{0}, msgBS...)
	}
	return msgBS
}

func bytesToUint32(msgBytes []byte) []uint32 {

	msgUint := []uint32{}

	for i := 0; i < len(msgBytes); i = i + 4 {
		word, err := strconv.ParseInt(fmt.Sprintf("%08b%08b%08b%08b", msgBytes[i], msgBytes[i+1], msgBytes[i+2], msgBytes[i+3]), 2, 64)
		if err != nil {
			log.Fatalln("Error converting bytes to uint32", err)
		}
		msgUint = append(msgUint, uint32(word))
	}
	return msgUint
}

func createMessageSchedule(chunk [16]uint32) [64]uint32 {

	var msgSchedule [64]uint32
	for j := range chunk {
		msgSchedule[j] = chunk[j]
	}

	for i := 16; i < 64; i++ {
		msgSchedule[i] = msgSchedule[i-16] + s0(msgSchedule[i-15]) + msgSchedule[i-7] + s1(msgSchedule[i-2])
	}

	return msgSchedule
}

func ROTR(x uint32, n uint32) uint32 {
	n = n % 32

	return x>>n | x<<(32-n)
}

func s0(x uint32) uint32 {

	return ROTR(x, 7) ^ ROTR(x, 18) ^ (x >> 3)
}

func s1(x uint32) uint32 {

	return ROTR(x, 17) ^ ROTR(x, 19) ^ (x >> 10)
}

func Σ0(x uint32) uint32 {
	return ROTR(x, 2) ^ ROTR(x, 13) ^ ROTR(x, 22)
}

func Σ1(x uint32) uint32 {
	return ROTR(x, 6) ^ ROTR(x, 11) ^ ROTR(x, 25)
}

func Ch(x uint32, y uint32, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

func Maj(x uint32, y uint32, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}
