package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const MatrixSize = 1000
const Threads = 8

func generateInitialMatrix(coef int) [][]int {
	initialMatrix := make([][]int, MatrixSize)
	for i := 0; i < MatrixSize; i++ {
		initialMatrix[i] = make([]int, MatrixSize)
		for j := 0; j < MatrixSize; j++ {
			initialMatrix[i][j] = i*MatrixSize/coef + j*coef
		}
	}
	return initialMatrix
}

func multiplySquareMatrixNon(a, b [][]int) [][]int {
	result := make([][]int, MatrixSize)
	for i := 0; i < MatrixSize; i++ {
		result[i] = make([]int, MatrixSize)
		for j := 0; j < MatrixSize; j++ {
			for k := 0; k < MatrixSize; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func multiplySquareMatrix(a, b [][]int, threads int) [][]int {
	result := make([][]int, MatrixSize)
	//fmt.Println("running on", threads, "threads")

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
						result[i][j] += a[i][k] * b[k][j]
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

func checkMatrix(a, b [][]int) {
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
	a := generateInitialMatrix(1)
	b := generateInitialMatrix(2)
	start1 := time.Now()
	res := multiplySquareMatrixNon(a, b)
	fmt.Println("Non thread", time.Since(start1))
	for i := 1; i <= Threads; i++ {
		start := time.Now()
		res2 := multiplySquareMatrix(a, b, i)
		delta := time.Since(start)
		fmt.Println(i, "threads", delta)
		checkMatrix(res, res2)
	}
}
