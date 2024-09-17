package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const MatrixSize = 1000
const Threads = 24

func generateRandomMatrix() [][]int {
	initialMatrix := make([][]int, MatrixSize)
	for i := 0; i < MatrixSize; i++ {
		initialMatrix[i] = make([]int, MatrixSize)
		for j := 0; j < MatrixSize; j++ {
			initialMatrix[i][j] = rand.Int()
		}
	}
	return initialMatrix
}

func multiplySquareMatrix(a, b [][]int) [][]int {
	result := make([][]int, MatrixSize)
	for i := 0; i < MatrixSize; i++ {
		result[i] = make([]int, MatrixSize)
		for j := 0; j < MatrixSize; j++ {
			for k := 0; k < MatrixSize; k++ {
				result[i][k] += a[i][j] * b[j][k]
			}
		}
	}
	return result
}

func AsyncMultiplySquareMatrix(a, b [][]int, threads int) [][]int {
	result := make([][]int, MatrixSize)

	rowsPerThread := MatrixSize / threads
	rowLimits := make([]int, 0)
	for i := 0; i < threads; i++ {
		rowLimits = append(rowLimits, i*rowsPerThread)
	}
	rowLimits = append(rowLimits, MatrixSize)

	wg := &sync.WaitGroup{}
	wg.Add(threads)
	for threadNum := 0; threadNum < threads; threadNum++ {
		go func(rowStart, rowFinish int) {
			defer wg.Done()
			for i := rowStart; i < rowFinish; i++ {
				result[i] = make([]int, MatrixSize)
				for j := 0; j < MatrixSize; j++ {
					for k := 0; k < MatrixSize; k++ {
						result[i][k] += a[i][j] * b[j][k]
					}
				}
			}
		}(rowLimits[threadNum], rowLimits[threadNum+1])
	}
	wg.Wait()
	return result
}

func printMatrix(a [][]int) {
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[i]); j++ {
			fmt.Print(a[i][j], ' ')
		}
		fmt.Println()
	}
}

func compareMatrices(a, b [][]int) {
	for i := 0; i < MatrixSize; i++ {
		for j := 0; j < MatrixSize; j++ {
			if a[i][j] != b[i][j] {
				fmt.Println("ACHTUNG!")
			}
		}
	}
}

func main() {
	fmt.Println("Available threads:", runtime.NumCPU())
	a := generateRandomMatrix()
	b := generateRandomMatrix()
	start1 := time.Now()
	res := multiplySquareMatrix(a, b)
	fmt.Println("Non thread", time.Since(start1))
	for i := 1; i <= Threads; i++ {
		start := time.Now()
		resAsync := AsyncMultiplySquareMatrix(a, b, i)
		delta := time.Since(start)
		fmt.Println(i, "threads", delta)
		compareMatrices(res, resAsync)
	}
}
