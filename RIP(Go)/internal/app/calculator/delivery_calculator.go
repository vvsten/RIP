package calculator

import (
	"math"
	"strings"
	"rip-go-app/internal/app/ds"
)

// DeliveryCalculator - калькулятор доставки
type DeliveryCalculator struct{}

// NewDeliveryCalculator - создание нового калькулятора
func NewDeliveryCalculator() *DeliveryCalculator {
	return &DeliveryCalculator{}
}

// DeliveryResult - результат расчета доставки
type DeliveryResult struct {
	DeliveryDays int     `json:"delivery_days"`
	TotalCost    float64 `json:"total_cost"`
	Distance     float64 `json:"distance"`
	Volume       float64 `json:"volume"`
	IsValid      bool    `json:"is_valid"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

// CalculateDelivery - основной метод расчета доставки
func (dc *DeliveryCalculator) CalculateDelivery(service ds.Service, fromCity, toCity string, length, width, height, weight float64) DeliveryResult {
	result := DeliveryResult{
		IsValid: true,
	}

	// Проверяем ограничения
	if !dc.validateConstraints(service, length, width, height, weight) {
		result.IsValid = false
		result.ErrorMessage = "Груз не соответствует ограничениям выбранного типа транспорта"
		return result
	}

	// Рассчитываем объем
	result.Volume = length * width * height

	// Рассчитываем расстояние
	result.Distance = dc.calculateDistance(fromCity, toCity)

	// Рассчитываем сроки доставки
	result.DeliveryDays = dc.calculateDeliveryDays(service, result.Distance, result.Volume, weight)

	// Рассчитываем стоимость
	result.TotalCost = dc.calculateCost(service, result.Distance, result.Volume, weight)

	return result
}

// validateConstraints - проверка ограничений
func (dc *DeliveryCalculator) validateConstraints(service ds.Service, length, width, height, weight float64) bool {
	volume := length * width * height
	
	// Проверяем вес
	if weight > service.MaxWeight {
		return false
	}
	
	// Проверяем объем
	if volume > service.MaxVolume {
		return false
	}
	
	// Проверяем габариты (максимальные размеры для каждого типа транспорта)
	maxDimensions := dc.getMaxDimensions(service.ID)
	if length > maxDimensions.Length || width > maxDimensions.Width || height > maxDimensions.Height {
		return false
	}
	
	return true
}

// MaxDimensions - максимальные габариты
type MaxDimensions struct {
	Length float64
	Width  float64
	Height float64
}

// getMaxDimensions - получение максимальных габаритов для типа транспорта
func (dc *DeliveryCalculator) getMaxDimensions(serviceID int) MaxDimensions {
	switch serviceID {
	case 1: // Фура
		return MaxDimensions{Length: 13.6, Width: 2.5, Height: 2.7}
	case 2: // Малотоннажный грузовик
		return MaxDimensions{Length: 6.0, Width: 2.0, Height: 2.2}
	case 3: // Авиаперевозка
		return MaxDimensions{Length: 3.0, Width: 1.5, Height: 1.5}
	case 4: // Поезд
		return MaxDimensions{Length: 20.0, Width: 3.0, Height: 3.0}
	case 5: // Корабль
		return MaxDimensions{Length: 40.0, Width: 8.0, Height: 8.0}
	case 6: // Мультимодальные
		return MaxDimensions{Length: 13.6, Width: 2.5, Height: 2.7}
	default:
		return MaxDimensions{Length: 6.0, Width: 2.0, Height: 2.2}
	}
}

// calculateDeliveryDays - расчет сроков доставки
func (dc *DeliveryCalculator) calculateDeliveryDays(service ds.Service, distance, volume, weight float64) int {
	// Базовые сроки
	baseDays := service.DeliveryDays
	
	// Коэффициенты для разных типов транспорта
	coefficients := dc.getDeliveryCoefficients(service.ID)
	
	// Расчет по расстоянию
	distanceDays := math.Ceil(distance / coefficients.DistancePerDay)
	
	// Дополнительные дни за сложность груза
	complexityDays := dc.calculateComplexityDays(volume, weight, service.ID)
	
	// Итоговые сроки
	totalDays := baseDays + int(distanceDays) + complexityDays
	
	// Минимальные сроки для каждого типа транспорта
	minDays := dc.getMinDeliveryDays(service.ID)
	if totalDays < minDays {
		totalDays = minDays
	}
	
	return totalDays
}

// DeliveryCoefficients - коэффициенты доставки
type DeliveryCoefficients struct {
	DistancePerDay float64 // км в день
	WeightFactor   float64 // коэффициент веса
	VolumeFactor   float64 // коэффициент объема
}

// getDeliveryCoefficients - получение коэффициентов для типа транспорта
func (dc *DeliveryCalculator) getDeliveryCoefficients(serviceID int) DeliveryCoefficients {
	switch serviceID {
	case 1: // Фура
		return DeliveryCoefficients{DistancePerDay: 800, WeightFactor: 0.1, VolumeFactor: 0.05}
	case 2: // Малотоннажный грузовик
		return DeliveryCoefficients{DistancePerDay: 600, WeightFactor: 0.15, VolumeFactor: 0.1}
	case 3: // Авиаперевозка
		return DeliveryCoefficients{DistancePerDay: 2000, WeightFactor: 0.2, VolumeFactor: 0.3}
	case 4: // Поезд
		return DeliveryCoefficients{DistancePerDay: 1200, WeightFactor: 0.05, VolumeFactor: 0.02}
	case 5: // Корабль
		return DeliveryCoefficients{DistancePerDay: 500, WeightFactor: 0.02, VolumeFactor: 0.01}
	case 6: // Мультимодальные
		return DeliveryCoefficients{DistancePerDay: 700, WeightFactor: 0.12, VolumeFactor: 0.08}
	default:
		return DeliveryCoefficients{DistancePerDay: 600, WeightFactor: 0.1, VolumeFactor: 0.05}
	}
}

// calculateComplexityDays - расчет дополнительных дней за сложность груза
func (dc *DeliveryCalculator) calculateComplexityDays(volume, weight float64, serviceID int) int {
	// Дополнительные дни за большой объем
	volumeDays := 0
	if volume > 20 {
		volumeDays = int(volume / 20) // +1 день за каждые 20 м³
	}
	
	// Дополнительные дни за большой вес
	weightDays := 0
	if weight > 1000 {
		weightDays = int(weight / 1000) // +1 день за каждые 1000 кг
	}
	
	// Для авиаперевозки меньше дополнительных дней
	if serviceID == 3 {
		volumeDays = volumeDays / 2
		weightDays = weightDays / 2
	}
	
	return volumeDays + weightDays
}

// getMinDeliveryDays - минимальные сроки доставки
func (dc *DeliveryCalculator) getMinDeliveryDays(serviceID int) int {
	switch serviceID {
	case 1: // Фура
		return 1
	case 2: // Малотоннажный грузовик
		return 1
	case 3: // Авиаперевозка
		return 1
	case 4: // Поезд
		return 2
	case 5: // Корабль
		return 3
	case 6: // Мультимодальные
		return 2
	default:
		return 1
	}
}

// calculateCost - расчет стоимости доставки
func (dc *DeliveryCalculator) calculateCost(service ds.Service, distance, volume, weight float64) float64 {
	// Базовая стоимость
	baseCost := service.Price
	
	// Коэффициенты стоимости
	costCoeffs := dc.getCostCoefficients(service.ID)
	
	// Стоимость за расстояние
	distanceCost := distance * costCoeffs.DistanceRate
	
	// Стоимость за вес
	weightCost := weight * costCoeffs.WeightRate
	
	// Стоимость за объем
	volumeCost := volume * costCoeffs.VolumeRate
	
	// Дополнительные коэффициенты
	complexityMultiplier := dc.calculateComplexityMultiplier(volume, weight, service.ID)
	
	// Итоговая стоимость
	totalCost := (baseCost + distanceCost + weightCost + volumeCost) * complexityMultiplier
	
	// Минимальная стоимость
	if totalCost < baseCost {
		totalCost = baseCost
	}
	
	// Округляем до рублей
	return math.Round(totalCost*100) / 100
}

// CostCoefficients - коэффициенты стоимости
type CostCoefficients struct {
	DistanceRate float64 // руб/км
	WeightRate   float64 // руб/кг
	VolumeRate   float64 // руб/м³
}

// getCostCoefficients - получение коэффициентов стоимости
func (dc *DeliveryCalculator) getCostCoefficients(serviceID int) CostCoefficients {
	switch serviceID {
	case 1: // Фура
		return CostCoefficients{DistanceRate: 15, WeightRate: 2, VolumeRate: 50}
	case 2: // Малотоннажный грузовик
		return CostCoefficients{DistanceRate: 12, WeightRate: 3, VolumeRate: 60}
	case 3: // Авиаперевозка
		return CostCoefficients{DistanceRate: 25, WeightRate: 8, VolumeRate: 200}
	case 4: // Поезд
		return CostCoefficients{DistanceRate: 8, WeightRate: 1, VolumeRate: 30}
	case 5: // Корабль
		return CostCoefficients{DistanceRate: 5, WeightRate: 0.5, VolumeRate: 20}
	case 6: // Мультимодальные
		return CostCoefficients{DistanceRate: 18, WeightRate: 2.5, VolumeRate: 80}
	default:
		return CostCoefficients{DistanceRate: 12, WeightRate: 2, VolumeRate: 50}
	}
}

// calculateComplexityMultiplier - расчет коэффициента сложности
func (dc *DeliveryCalculator) calculateComplexityMultiplier(volume, weight float64, serviceID int) float64 {
	multiplier := 1.0
	
	// Коэффициент за большой объем
	if volume > 10 {
		volumeFactor := volume / 10
		multiplier += volumeFactor * 0.1 // +10% за каждые 10 м³
	}
	
	// Коэффициент за большой вес
	if weight > 500 {
		weightFactor := weight / 500
		multiplier += weightFactor * 0.05 // +5% за каждые 500 кг
	}
	
	// Для авиаперевозки больше коэффициент сложности
	if serviceID == 3 {
		multiplier *= 1.2
	}
	
	// Максимальный коэффициент
	if multiplier > 2.0 {
		multiplier = 2.0
	}
	
	return multiplier
}

// calculateDistance - расчет расстояния между городами (улучшенная версия)
func (dc *DeliveryCalculator) calculateDistance(fromCity, toCity string) float64 {
	// Приводим к нижнему регистру для сравнения
	from := strings.ToLower(strings.TrimSpace(fromCity))
	to := strings.ToLower(strings.TrimSpace(toCity))
	
	// Если города одинаковые
	if from == to {
		return 0
	}
	
	// Расширенная база данных расстояний между основными городами
	distances := map[string]map[string]float64{
		"москва": {
			"санкт-петербург": 635,
			"спб":             635,
			"екатеринбург":    1416,
			"новосибирск":     3354,
			"красноярск":      4205,
			"иркутск":         5152,
			"владивосток":     9100,
			"ростов-на-дону":  1070,
			"сочи":            1360,
			"казань":          820,
			"нижний новгород": 420,
			"самара":          1050,
			"волгоград":       970,
			"воронеж":         520,
			"саратов":         850,
			"пермь":           1380,
			"уфа":             1160,
			"челябинск":       1510,
			"омск":            2550,
			"тюмень":          1720,
			"краснодар":       1350,
			"ставрополь":      1400,
			"астрахань":       1400,
			"махачкала":       1800,
			"грозный":         1900,
			"элиста":          1200,
			"йошкар-ола":      650,
			"чебоксары":       650,
			"ижевск":          1200,
			"киров":           900,
			"сыктывкар":       1400,
			"архангельск":     1200,
			"мурманск":        1900,
			"петрозаводск":    1000,
			"калининград":     1200,
		},
		"санкт-петербург": {
			"спб":             0,
			"москва":          635,
			"екатеринбург":    1780,
			"новосибирск":     3720,
			"калининград":     550,
			"мурманск":        1050,
			"архангельск":     1130,
			"петрозаводск":    320,
			"великий новгород": 180,
			"псков":           280,
			"тверь":           480,
			"вологда":         700,
			"череповец":       650,
		},
		"екатеринбург": {
			"москва":          1416,
			"санкт-петербург": 1780,
			"спб":             1780,
			"новосибирск":     1940,
			"челябинск":       200,
			"пермь":           360,
			"тюмень":          320,
			"уфа":             520,
			"курган":          380,
			"оренбург":        800,
			"магнитогорск":    300,
		},
		"новосибирск": {
			"москва":          3354,
			"санкт-петербург": 3720,
			"спб":             3720,
			"екатеринбург":    1940,
			"омск":            650,
			"красноярск":      850,
			"томск":           270,
			"барнаул":         230,
			"кемерово":        260,
			"новокузнецк":     300,
			"бийск":           360,
			"горно-алтайск":   450,
		},
		"красноярск": {
			"москва":          4205,
			"санкт-петербург": 4570,
			"спб":             4570,
			"новосибирск":     850,
			"иркутск":         1060,
			"абакан":          410,
			"кызыл":           460,
			"норильск":        1500,
			"дудинка":         1600,
		},
		"иркутск": {
			"москва":          5152,
			"санкт-петербург": 5520,
			"спб":             5520,
			"красноярск":      1060,
			"улан-удэ":        450,
			"чита":            1100,
			"якутск":          2000,
			"магадан":         3000,
			"петропавловск-камчатский": 4000,
		},
		"владивосток": {
			"москва":          9100,
			"санкт-петербург": 9470,
			"спб":             9470,
			"хабаровск":       760,
			"южно-сахалинск":  1000,
			"благовещенск":    1100,
			"петропавловск-камчатский": 2000,
		},
	}
	
	// Ищем расстояние в базе данных
	if cityDistances, exists := distances[from]; exists {
		if distance, found := cityDistances[to]; found {
			return distance
		}
	}
	
	// Если расстояние не найдено, используем примерную оценку
	// Базовое расстояние для неизвестных маршрутов
	return 500.0
}
