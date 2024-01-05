package luhn

import (
	"errors"

	"github.com/MWT-proger/go-loyalty-system/internal/logger"
)

const (
	asciiZero = 48
	asciiTen  = 57
)

// Validate(number string) bool принимает номер в виде строки
// и проверяет на соответствие алгортму Луна
func Validate(number string) bool {

	p := len(number) % 2
	sum, err := calculateLuhnSum(number, p)

	if err != nil {
		logger.Log.Error(err.Error())
		return false
	}

	// Если сумма по модулю 10 не равна 0, то число недействительно.
	if sum%10 != 0 {
		logger.Log.Error("invalid number")
		return false
	}

	return true
}

func calculateLuhnSum(number string, parity int) (int64, error) {
	var sum int64
	for i, d := range number {
		if d < asciiZero || d > asciiTen {
			return 0, errors.New("invalid digit")
		}

		d = d - asciiZero
		// Удваиваем значение каждой второй цифры.
		if i%2 == parity {

			d *= 2

			if d > 9 {
				d -= 9
			}
		}

		// Берем сумму всех цифр.
		sum += int64(d)
	}

	return sum, nil
}
