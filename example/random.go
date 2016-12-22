package main

import (
	"bytes"
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

func newRand() *rand.Rand {
	return random.NewRand()
}

type alphabet struct {
	lower []rune
	upper []rune
}

var (
	abcEnglish = alphabet{
		lower: []rune("abcdefghijklmnopqrstuvwxyz"),
		upper: []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	abcRussian = alphabet{
		lower: []rune("абвгдеёжзийклмнопрстуфхчцшщъыьэюя"),
		upper: []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"),
	}

	abcUkrainian = alphabet{
		lower: []rune("абвгґдеєжзиіїйклмнопрстуфхцчшщьюя"),
		upper: []rune("АБВГҐДЕЄЖЗИІЇЙКЛМНОПРСТУФХЦЧШЩЬЮЯ"),
	}
)

var (
	digits       = []rune("0123456789")
	specialRunes = []rune("! @ # $ % ^ * ( ) _ + , . / ?")
	hieroglyphs  = []rune("座高楼个很老的人个非常好的朋友")
)

var randRunesSamples = [][]rune{
	digits,
	specialRunes,
	abcEnglish.lower,
	abcEnglish.upper,
	abcRussian.lower,
	abcRussian.upper,
	abcUkrainian.lower,
	abcUkrainian.upper,
	hieroglyphs,
}

func randWord(r *rand.Rand, n int) string {

	var rs []rune

	switch k := r.Intn(3); k {
	case 0:
		rs = abcEnglish.lower
	case 1:
		rs = abcRussian.lower
	case 2:
		rs = abcUkrainian.lower
	}

	word := make([]rune, n)

	for i := 0; i < n; i++ {
		word[i] = rs[r.Intn(len(rs))]
	}

	return string(word)
}

func randString(r *rand.Rand, maxLen int) string {

	if maxLen < 0 {
		maxLen = 0
	}

	var buffer bytes.Buffer

	n := random.Intn(r, randRange(r, 3, maxLen))

	for i := 0; i < n; i++ {

		if i > 0 {
			switch k := r.Intn(3); k {
			case 0:
				buffer.WriteRune(' ')
			case 1:
				buffer.WriteString(", ")
			case 2:
				buffer.WriteString(". ")
			}
		}

		word := randWord(r, randRange(r, 3, 8))

		buffer.WriteString(word)
	}

	return buffer.String()
}

func randRange(r *rand.Rand, min, max int) int {
	if min > max {
		min, max = max, min
	}
	return min + r.Intn(max-min)
}
