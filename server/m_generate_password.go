package server

import (
	"context"
	"errors"
	"log/slog"
	mr "math/rand"
	"strings"

	pb "github.com/kTowkA/GophKeeper/grpc"
)

func (s *Server) GeneratePassword(ctx context.Context, r *pb.GeneratePasswordRequest) (*pb.GeneratePasswordResponse, error) {
	password, err := generatePassword(int(r.Length))
	if err != nil {
		s.log.Error("генерация пароля", slog.String("ошибка", err.Error()))
		return nil, errors.New("ошибка генерации пароля")
	}
	return &pb.GeneratePasswordResponse{
		Password: password,
	}, nil
}

// generatePassword генератор псевдослучайных паролей длиной length (минимальное значение 4). возвращаеть пароль и nil (для возможных будущих изменений если все таки генерировать криптостойкий пароль)
// не стал использовать и искать какие-то готовые решения, т.к. надо в идеале проверять криптостойкость генерируемых паролей и написал простенький генератор для себя.
// генерируемый пароль в равной степени состоит из символов алфавита в нижнем и верхнем регистрах, числел и специальных символов
func generatePassword(length int) (string, error) {
	// шаблоны 4-х частей для будущего пароля
	abcLower := "qwertyuiopasdfghjklzxcvbnm"
	abcUpper := "QWERTYUIOPASDFGHJKLZXCVBNM"
	numbers := "1234567890"
	symbols := "!@#$%^&*()_+-|."

	// проверяем минимальную длину
	if length < 4 {
		length = 4
	}

	// смотрим размер каждой части в пароле и остаток, чтобы увеличить часть (приоритет: нижний регистр-верхний регистр-числа)
	countInPart := length / 4
	o := length % 4
	lengthPartAbcLower := countInPart
	lengthPartAbcUpper := countInPart
	lengthPartNumbers := countInPart
	lengthPartSymbols := countInPart
	if o > 0 {
		lengthPartAbcLower++
		o--
	}
	if o > 0 {
		lengthPartAbcUpper++
		o--
	}
	if o > 0 {
		lengthPartNumbers++
	}

	// разбиваем наши шаблоны и смешиваем (пвевдорандом)
	abcLowerArray := strings.Split(abcLower, "")
	abcUpperArray := strings.Split(abcUpper, "")
	numbersArray := strings.Split(numbers, "")
	symbolsArray := strings.Split(symbols, "")
	mr.Shuffle(len(abcLower), func(i, j int) {
		abcLowerArray[i], abcLowerArray[j] = abcLowerArray[j], abcLowerArray[i]
	})
	mr.Shuffle(len(abcUpperArray), func(i, j int) {
		abcUpperArray[i], abcUpperArray[j] = abcUpperArray[j], abcUpperArray[i]
	})
	mr.Shuffle(len(numbersArray), func(i, j int) {
		numbersArray[i], numbersArray[j] = numbersArray[j], numbersArray[i]
	})
	mr.Shuffle(len(symbolsArray), func(i, j int) {
		symbolsArray[i], symbolsArray[j] = symbolsArray[j], symbolsArray[i]
	})

	// создаем пароль нужной длины и добавляем в него первые элементы и шаблона нужной длины
	password := make([]string, 0, length)
	password = append(password, abcLowerArray[:lengthPartAbcLower]...)
	password = append(password, abcUpperArray[:lengthPartAbcUpper]...)
	password = append(password, numbersArray[:lengthPartNumbers]...)
	password = append(password, symbolsArray[:lengthPartSymbols]...)

	// снова смешиваем
	mr.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return strings.Join(password, ""), nil
}
