package main

import (
	"fmt"
	"time"
)

func main() {
	// Получаем текущую дату и время
	now := time.Now()

	// Шаг 1: Определяем круглое столетие
	year := now.Year()
	century := (year / 100) * 100 // Округляем вниз до ближайшего столетия

	// Шаг 2: Вычисляем количество дней с начала столетия
	centuryStart := time.Date(century, 1, 1, 0, 0, 0, 0, time.UTC)
	daysSinceCentury := int(now.Sub(centuryStart).Hours() / 24)

	// Шаг 3: Вычисляем минуты с начала дня с округлением вверх
	hours := now.Hour()
	minutes := now.Minute()
	seconds := now.Second()
	minutesSinceDay := hours*60 + minutes
	if seconds > 0 {
		minutesSinceDay++ // Округляем вверх, если есть секунды
	}

	// Шаг 4: Суммируем и ограничиваем до 65535
	buildID := (daysSinceCentury + minutesSinceDay) % 65536

	// Выводим результат
	fmt.Println(buildID)
}
