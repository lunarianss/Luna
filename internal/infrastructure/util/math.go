package util

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func convertFloat32ToFloat64(input []float32) []float64 {
	output := make([]float64, len(input))

	for i, v := range input {
		output[i] = float64(v)
	}

	return output
}

func convertFloat64ToFloat32(input []float64) []float32 {
	output := make([]float32, len(input))

	for i, v := range input {
		output[i] = float32(v)
	}

	return output
}

func NormalizeVector(data []float32) ([]float32, error) {
	dataFLoat64 := convertFloat32ToFloat64(data)
	vec := mat.NewVecDense(len(data), dataFLoat64)
	norm := mat.Norm(vec, 2)

	if norm == 0 {
		return make([]float32, len(data)), fmt.Errorf("zero vector cannot be normalized")
	}

	normalized := mat.NewVecDense(vec.Len(), nil)
	normalized.ScaleVec(1/norm, vec)

	return convertFloat64ToFloat32(vec.RawVector().Data), nil
}
