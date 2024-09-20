package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const MatrixSize = 5000
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
	computeRows := func(wg *sync.WaitGroup, rowStart, rowFinish, routineIndex int) {
		//start := time.Now()
		defer wg.Done()
		for i := rowStart; i < rowFinish; i++ {
			result[i] = make([]int, MatrixSize)
			for j := 0; j < MatrixSize; j++ {
				for k := 0; k < MatrixSize; k++ {
					result[i][k] += a[i][j] * b[j][k]
				}
			}
		}
		//fmt.Println("\tRoutine", routineIndex, "finished", time.Since(start))
	}

	wg := &sync.WaitGroup{}
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		if i+1 == threads {
			go computeRows(wg, i*rowsPerThread, MatrixSize, i)
		} else {
			go computeRows(wg, i*rowsPerThread, (i+1)*rowsPerThread, i)
		}
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
				panic("ACHTUNG! Matrices do not match!\n")
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
	for i := 2; i <= Threads; i++ {
		start := time.Now()
		resAsync := AsyncMultiplySquareMatrix(a, b, i)
		delta := time.Since(start)
		fmt.Println(i, ",", delta)
		compareMatrices(res, resAsync)
	}
}
