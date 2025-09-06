package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func main() {
	for _, r := range decode("ok©世界😀", true) {
		fmt.Println(r)
	}
}

/*
	| Rune range         | Bytes needed | Leading byte pattern | Continuation byte pattern |
	|--------------------+--------------+----------------------+---------------------------|
	| U+0000  – U+007F   |      1       |      0xxxxxxx        |            N/A            |
	| U+0080  – U+07FF   |      2       |      110xxxxx        |         10xxxxxx          |
	| U+0800  – U+FFFF   |      3       |      1110xxxx        |         10xxxxxx          |
	| U+10000 – U+10FFFF |      4       |      11110xxx        |         10xxxxxx          |
*/

type CodePoint struct {
	rawBytes     []byte
	payloadBytes []byte
	payload      rune
}

func (cp CodePoint) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("U+%.4x - %c", cp.payload, cp.payload))

	if len(cp.payloadBytes) > 0 {
		var contByte string

		sb.WriteString(fmt.Sprintf("\n[requires %d bytes]\n", len(cp.payloadBytes)))
		for i, p := range cp.payloadBytes {
			if i > 0 {
				contByte = " cont. byte"
			}
			sb.WriteString(fmt.Sprintf("   byte %d [0x%.2x - %.8b], masked [0x%.2x %.8b]%s\n", i, cp.rawBytes[i], cp.rawBytes[i], p, p, contByte))
		}
	}

	return sb.String()
}

func decode(data string, verbose bool) []CodePoint {
	cps := make([]CodePoint, 0, utf8.RuneCountInString(data))

	var j int

	for i := 0; i < len(data); i++ {
		b := data[i]

		switch {
		// 11110xxx - 4 bytes
		case b&0xF0 == 0xF0:
			firstBytePayload := data[i] ^ 0b11110000
			secondBytePayload := data[i+1] ^ 0b10000000
			thirdBytePayload := data[i+2] ^ 0b10000000
			fourthBytePayload := data[i+3] ^ 0b10000000
			codePoint := CodePoint{
				payload: toCodePoint(rune(firstBytePayload), rune(secondBytePayload), rune(thirdBytePayload), rune(fourthBytePayload)),
			}
			if verbose {
				codePoint.rawBytes = []byte{data[i], data[i+1], data[i+2], data[i+3]}
				codePoint.payloadBytes = []byte{firstBytePayload, secondBytePayload, thirdBytePayload, fourthBytePayload}
			}

			cps = append(cps, codePoint)

			i += 3
			j++

		// 1110xxxx - 3 bytes
		case b&0xE0 == 0xE0:
			firstBytePayload := data[i] ^ 0b11100000
			secondBytePayload := data[i+1] ^ 0b10000000
			thirdBytePayload := data[i+2] ^ 0b10000000
			codePoint := CodePoint{
				payload: toCodePoint(rune(firstBytePayload), rune(secondBytePayload), rune(thirdBytePayload)),
			}
			if verbose {
				codePoint.rawBytes = []byte{data[i], data[i+1], data[i+2]}
				codePoint.payloadBytes = []byte{firstBytePayload, secondBytePayload, thirdBytePayload}
			}

			cps = append(cps, codePoint)

			i += 2
			j++

		// 110xxxxx - 2 bytes
		case b&0xC0 == 0xC0:
			firstBytePayload := data[i] ^ 0b11000000
			secondBytePayload := data[i+1] ^ 0b10000000
			codePoint := CodePoint{
				payload: toCodePoint(rune(firstBytePayload), rune(secondBytePayload)),
			}
			if verbose {
				codePoint.rawBytes = []byte{data[i], data[i+1]}
				codePoint.payloadBytes = []byte{firstBytePayload, secondBytePayload}
			}

			cps = append(cps, codePoint)

			i++
			j++

		// 0xxxxxxx - 1 byte
		case b&0x80 == 0:
			firstBytePayload := data[i]
			codePoint := CodePoint{
				payload: toCodePoint(rune(firstBytePayload)),
			}
			if verbose {
				codePoint.rawBytes = []byte{data[i]}
				codePoint.payloadBytes = []byte{firstBytePayload}
			}

			cps = append(cps, codePoint)

			j++

		default:
			fmt.Println("unrecognised bitpattern")
		}
	}

	return cps
}

func toCodePoint(rs ...rune) rune {
	if len(rs) == 0 {
		return 0
	}

	var res rune

	for i, j := 0, len(rs)-1; i < len(rs) && j >= 0; i, j = i+1, j-1 {
		res |= rs[i] << (j * 6)
	}

	return res
}
