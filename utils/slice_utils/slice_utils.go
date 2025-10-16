package sliceutils

import (
	"maps"
)

func UniquesStringFamily[T ~string](slc []T) []T {

	uniques := make(map[T]any)

	for _, value := range slc {

		if _, exists := uniques[value]; !exists {
			uniques[value] = struct{}{}
		}

	}

	result := make([]T, 0, len(uniques))

	for ite := range maps.Keys(uniques) {

		result = append(result, ite)

	}

	return result
}
