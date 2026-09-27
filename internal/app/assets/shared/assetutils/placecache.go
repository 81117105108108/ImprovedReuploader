package assetutils

import "github.com/kartFr/Asset-Reuploader/internal/atomicarray"

func MoveValueToTop[T comparable](arr *atomicarray.AtomicArray[T], value T) {
	arr.Update(func(currentArray []T) []T {
		if len(currentArray) == 0 {
			return nil
		}
		if currentArray[0] == value {
			return nil
		}

		for i, v := range currentArray {
			if v != value {
				continue
			}
			if i == 1 {
				currentArray[0], currentArray[1] = currentArray[1], currentArray[0]
				return currentArray
			}

			copy(currentArray[1:i+1], currentArray[0:i])
			currentArray[0] = value
			return currentArray
		}

		return nil
	})
}
