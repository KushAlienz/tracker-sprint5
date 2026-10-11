package spentenergy

import (
	"errors"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

var (
	ErrSteps    = errors.New("количество шагов должно быть больше 0")
	ErrWeight   = errors.New("вес должен быть больше 0")
	ErrHeight   = errors.New("рост должен быть больше 0")
	ErrDuration = errors.New("продолжительность должна быть больше 0")
)

// Distance возвращает дистанцию в километрах.
func Distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient // длина шага в метрах
	distanceMeters := float64(steps) * stepLength
	return distanceMeters / mInKm
}

// MeanSpeed возвращает среднюю скорость в км/ч.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 || steps < 0 {
		return 0
	}
	return Distance(steps, height) / duration.Hours()
}

// validateInput проверяет корректность входных данных для расчёта калорий.
func validateInput(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return ErrSteps
	case weight <= 0:
		return ErrWeight
	case height <= 0:
		return ErrHeight
	case duration <= 0:
		return ErrDuration
	}
	return nil
}

// RunningSpentCalories возвращает калории, потраченные при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}
	meanSpeed := MeanSpeed(steps, height, duration)
	return weight * meanSpeed * duration.Minutes() / minInH, nil
}

// WalkingSpentCalories возвращает калории, потраченные при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}
	meanSpeed := MeanSpeed(steps, height, duration)
	calories := weight * meanSpeed * duration.Minutes() / minInH
	return calories * walkingCaloriesCoefficient, nil
}
