package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

var ErrUnknownTrainingType = errors.New("неизвестный тип тренировки")

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат данных: %q", datastring)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("ошибка парсинга шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть больше 0, получено %d", steps)
	}
	t.Steps = steps
	t.TrainingType = parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("ошибка парсинга длительности: %w", err)
	}
	if duration <= 0 {
		return fmt.Errorf("продолжительность должна быть больше 0, получено %s", duration)
	}
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	var spentCalories float64
	var err error

	switch t.TrainingType {
	case "Бег":
		spentCalories, err = spentenergy.RunningSpentCalories(
			t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		spentCalories, err = spentenergy.WalkingSpentCalories(
			t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrUnknownTrainingType
	}
	if err != nil {
		return "", err
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	info := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, meanSpeed, spentCalories,
	)
	return info, nil
}
