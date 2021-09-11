package pkg

import (
	"strings"
	"testing"
)

func DehydrateAndUpgradeTest(t *testing.T) {
	// проверяем, подходит символ под наш base63 алфавит
	f := func(r rune) bool {
		return r < '0' || r > '9' && (r < 'A' || (r > 'z' && r != '_'))
	}

	for i := 0; i < 10000; i++ {
		if v := DehydrateAndUpgrade(i); len(v) != 10 || strings.IndexFunc(v, f) != -1 {
			t.Errorf("Ошибка с %v, получен %v.\n Длина: %v\n Лишние символы: %v",
				i, v, len(v), strings.IndexFunc(v, f) != -1)
		}
	}

}
