package monitor

import (
	"math"
	"sort"
)

// calculateCoeffiecient takes an int-int map containing temp-% values and converts it to a sorted array containing coefficients per temp
func calculateCoefficient(powerCurve map[int]int) []coefMatch {
	// sort temps
	var sortedTemps = make([][2]int, 0, len(powerCurve))

	for k, v := range powerCurve {
		sortedTemps = append(sortedTemps, [2]int{k, v})
	}

	sort.Slice(sortedTemps, func(i, j int) bool {
		return sortedTemps[i][0] < sortedTemps[j][0]
	})

	// calculate coeffiecients

	var powerCurveCoefficients = make([]coefMatch, 0, len(sortedTemps))

	for e, i := range sortedTemps {
		// dont go out of bounds
		if e+1 != len(sortedTemps) {
			// dT/dp
			coef := float64(sortedTemps[e+1][1]-i[1]) / float64(sortedTemps[e+1][0]-i[0])
			powerCurveCoefficients = append(powerCurveCoefficients, coefMatch{temp: i[0], coef: coef})
		} else {
			// add coef of 0 for last value
			powerCurveCoefficients = append(powerCurveCoefficients, coefMatch{i[0], 0})
		}
	}

	return powerCurveCoefficients
}

// calculates the right speed in % based on the coefficients, powercurve and temperature
func calculateSpeedsFromTemp(temp float64, coefficients []coefMatch, powerCurve map[int]int) int {
	for i, match := range coefficients {
		if temp >= float64(match.temp) && (i+1 == len(coefficients) || temp < float64(coefficients[i+1].temp)) {
			diff := temp - float64(match.temp)
			increase := diff * match.coef

			// pumpspeed will be higher than 100 if temp is higher than highest set temp
			return min(int(math.Round(float64(powerCurve[match.temp])+increase)), 100)
		}
	}

	// temp is lower than lowest set value
	return powerCurve[coefficients[0].temp]
}

func min(one, two int) int {
	if one < two {
		return one
	}

	return two
}
