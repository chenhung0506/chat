package service

import "errors"

// func indexOf[T comparable](slice []T, value int) (T,error) {
// 	if len(slice) == 0 {
// 		return _, nil
// 	}
// 	for i, v := range slice {
// 		if i == value {
// 			return v, nil
// 		}
// 	}
// }
func indexOf[T comparable](slice []T, index int) (T, error) {
	// Default zero value for the generic type T
	var zeroValue T

	if len(slice) == 0 {
		return zeroValue, errors.New("slice is empty")
	}
	if index < 0 || index >= len(slice) {
		return zeroValue, errors.New("index out of range")
	}
	return slice[index], nil
}