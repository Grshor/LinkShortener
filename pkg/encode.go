package pkg

import (
	"strings"
)

// Алфавит кодировки длиной в 63 символа
const encodeRuneSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
const encodeRuneSetLength = len(encodeRuneSet)

var sliceAlphabet = []rune(encodeRuneSet) // В действительности меняться не будет

// EncodeAndUpgrade - наверное медленный кодировщик числа в строку с алфавитом в 63 символа,
// сделан из-за ограничений алфавита (по условию задания) (иначе я бы использовал encoding/base64)
func EncodeAndUpgrade(id int) (shortString string) {
	var remainder int
	var sBuilder strings.Builder

	for id > 0 {
		remainder = id % encodeRuneSetLength
		sBuilder.WriteRune(sliceAlphabet[remainder])
		id /= encodeRuneSetLength
	}
	shortString = upgradeLength(sBuilder)
	return
}

// мы получаем строку неизвестно какой длинны, а нужна строка длинны 10, поэтому мы её доращиваем до нужной длинны
// TODO - сделать более элегантное решение
func upgradeLength(sBuilder strings.Builder) string {
	taleRune := sliceAlphabet[encodeRuneSetLength-1]
	for len(sBuilder.String()) < 10 {
		sBuilder.WriteRune(taleRune)
	}
	return sBuilder.String()
}
