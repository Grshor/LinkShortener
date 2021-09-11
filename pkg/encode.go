package pkg

import (
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_" // алфавит длиной в 63 символа
const alphabetLen = len(alphabet)

var sliceAlphabet = []rune(alphabet)

// DehydrateAndUpgrade - наверное медленный кодировщик числа в строку с алфавитом в 63 символа,
// сделан из-за ограничений алфавита-> нельзя использовать encoding/base64
func DehydrateAndUpgrade(id int) (shortString string) {
	var remainder int
	var sBuilder strings.Builder
	for id > 0 {
		remainder = id % alphabetLen
		sBuilder.WriteRune(sliceAlphabet[remainder])
		id /= alphabetLen
	}
	// мы получаем строку неизвестно какой длинны, а нужна строка длинны 10, поэтому мы её доращиваем до нужной длинны
	// думаю тут может быть риск коллизии, но что-то сделать надо, а решения лучше я не придумал
	taleRune := sliceAlphabet[alphabetLen-1]
	for len(sBuilder.String()) < 10 {
		sBuilder.WriteRune(taleRune)
	}

	shortString = sBuilder.String()
	return
}

// Неиспользуемый код, который не удаляю, чтобы не писать заново, если используемый не подойдёт.

// HashFnv хэширует данную строку в строку из чисел длинной в 10 символов
//func HashFnv(s string) string {
//	h := fnv.New32a()
//	h.Write([]byte(s))
//	return strconv.FormatUint(uint64(h.Sum32()), 10)
//}
//
//// IntToShortString я бы использовал, если бы знал, можно ли в результате использовать символ '-',
//// содержимый в алфавите рун для кодировки в пакете base64.
//func IntToShortString(id int) (shortString string) {
//	idS := strconv.Itoa(id)
//	uwEnc := base64.RawURLEncoding.EncodeToString([]byte(idS))
//	// дальше надо 'дорастить' uwEnc до 10 символов
//	return uwEnc
//}
