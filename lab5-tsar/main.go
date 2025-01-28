package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	rows       = 1000
	cols       = 1000
	numThreads = 12
	numSteps   = 10
)

type pair struct {
	first  int
	second int
}

func countLiveNeighbors(matrix [][]int, x, y, rows, cols int) int {
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	liveNeighbors := 0
	for _, dir := range directions {
		nx := (x + dir[0] + rows) % rows
		ny := (y + dir[1] + cols) % cols
		liveNeighbors += matrix[nx][ny]
	}
	return liveNeighbors
}

func evolveCell(matrix, newMatrix [][]int, x, y int) {
	liveNeighbors := countLiveNeighbors(matrix, x, y, rows, cols)
	if matrix[x][y] == 1 {
		if liveNeighbors < 2 || liveNeighbors > 3 {
			newMatrix[x][y] = 0
		} else {
			newMatrix[x][y] = 1
		}
	} else {
		if liveNeighbors == 3 {
			newMatrix[x][y] = 1
		} else {
			newMatrix[x][y] = 0
		}
	}
}

func countLiveNeighborsAs(matrix [][]int, x, y, rows, cols int, firstReqRow []int, secondReqRow []int) int {
	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	liveNeighbors := 0
	for _, dir := range directions {
		nx := x + dir[0]
		ny := (y + dir[1] + cols) % cols
		if nx == -1 {
			liveNeighbors += firstReqRow[ny]
		} else if nx == len(matrix) {
			liveNeighbors += secondReqRow[ny]
		} else {
			liveNeighbors += matrix[nx][ny]
		}
	}
	return liveNeighbors
}

func evolveCellAs(matrix, newMatrix [][]int, x, y int, firstReqRow []int, secondReqRow []int) {
	liveNeighbors := countLiveNeighborsAs(matrix, x, y, rows, cols, firstReqRow, secondReqRow)
	if matrix[x][y] == 1 {
		if liveNeighbors < 2 || liveNeighbors > 3 {
			newMatrix[x][y] = 0
		} else {
			newMatrix[x][y] = 1
		}
	} else {
		if liveNeighbors == 3 {
			newMatrix[x][y] = 1
		} else {
			newMatrix[x][y] = 0
		}
	}
}

func serveRows(threadID int, matrix [][]int, rowReceiver []chan []int, rowRequester []chan pair) {
	firstRequested := <-rowRequester[threadID]
	if firstRequested.second == 1 {
		rowReceiver[firstRequested.first] <- matrix[len(matrix)-1]
	} else {
		rowReceiver[firstRequested.first] <- matrix[0]
	}
	secondRequested := <-rowRequester[threadID]
	if secondRequested.second == 1 {
		rowReceiver[secondRequested.first] <- matrix[len(matrix)-1]
	} else {
		rowReceiver[secondRequested.first] <- matrix[0]
	}
}

func requestRows(threadID int, matrix [][]int, rowReceiver []chan []int, rowRequester []chan pair) ([]int, []int) {
	rowRequester[(threadID-1+numThreads)%numThreads] <- pair{threadID, 1}
	firstRequiredRow := <-rowReceiver[threadID]
	rowRequester[(threadID+1+numThreads)%numThreads] <- pair{threadID, 0}
	secondRequiredRow := <-rowReceiver[threadID]
	return firstRequiredRow, secondRequiredRow
}

func worker(threadID int, matrix, newMatrix [][]int, numThreads int, wg *sync.WaitGroup, rowReceiver []chan []int, rowRequester []chan pair) {
	defer wg.Done()
	var firstRequiredRow []int
	var secondRequiredRow []int
	if numThreads%2 == 0 {
		if threadID%2 == 0 {
			firstRequiredRow, secondRequiredRow = requestRows(threadID, matrix, rowReceiver, rowRequester)
			serveRows(threadID, matrix, rowReceiver, rowRequester)
		} else {
			serveRows(threadID, matrix, rowReceiver, rowRequester)
			firstRequiredRow, secondRequiredRow = requestRows(threadID, matrix, rowReceiver, rowRequester)
		}
	} else {
		if threadID%2 == 0 {
			firstRequiredRow, secondRequiredRow = requestRows(threadID, matrix, rowReceiver, rowRequester)
			serveRows(threadID, matrix, rowReceiver, rowRequester)
		} else {
			serveRows(threadID, matrix, rowReceiver, rowRequester)
			firstRequiredRow, secondRequiredRow = requestRows(threadID, matrix, rowReceiver, rowRequester)
		}
	}
	for x := 0; x < len(matrix); x++ {
		for y := 0; y < cols; y++ {
			evolveCellAs(matrix, newMatrix, x, y, firstRequiredRow, secondRequiredRow)
		}
	}
}

func main() {
	matrix := make([][]int, rows)
	newMatrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		newMatrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			matrix[i][j] = rand.Intn(2)
		}
	}

	wg := &sync.WaitGroup{}

	start := time.Now()
	for step := 0; step < numSteps; step++ {
		rowReceiver := make([]chan []int, numThreads)
		rowRequester := make([]chan pair, numThreads)
		wg.Add(numThreads)
		for i := 0; i < numThreads; i++ {
			rowReceiver[i] = make(chan []int)
			rowRequester[i] = make(chan pair, 2)
		}
		for i := 0; i < numThreads; i++ {

			startRow := i * (rows / numThreads)
			endRow := (i + 1) * (rows / numThreads)
			if i == numThreads-1 {
				endRow = rows
			}
			go worker(i, matrix[startRow:endRow], newMatrix, numThreads, wg, rowReceiver, rowRequester)
		}
		wg.Wait()
		matrix, newMatrix = newMatrix, matrix
	}
	elapsed := time.Since(start)
	fmt.Printf("Среднее время шага при выполнении аcинхронно на %d потоках: %.6f секунд\n", numThreads, elapsed.Seconds()/numSteps)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			matrix[i][j] = rand.Intn(2)
		}
	}

	start = time.Now()
	for step := 0; step < numSteps; step++ {
		for x := 0; x < rows; x++ {
			for y := 0; y < cols; y++ {
				evolveCell(matrix, newMatrix, x, y)
			}
		}
		matrix, newMatrix = newMatrix, matrix
	}
	elapsed = time.Since(start)
	fmt.Printf("Среднее время шага при выполнении cинхронно: %.6f секунд\n", elapsed.Seconds()/numSteps)
}
