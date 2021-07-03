package ai

import (
	"curls/models"
	"fmt"
	"math"
)

var directions = []vertex{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

type vertex [2]int

type route struct {
	cost     int
	previous vertex
}

func NextMove(game models.GameRequest) models.MoveResponse {
	size := game.Board.Height

	// previous + cost
	visited := map[vertex]bool{}

	costTable := map[vertex]route{}

	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int, size)
	}
	for y := range matrix {
		for x := range matrix[0] {
			matrix[y][x] = 100
		}
	}

	you := game.You
	for _, coord := range you.Body {
		matrix[coord.Y][coord.X] = math.MaxInt64
	}
	head := you.Head
	matrix[head.Y][head.X] = math.MaxInt64

	tail := you.Body[len(you.Body)-1]
	tailVertex := vertex{tail.Y, tail.X}
	matrix[tail.Y][tail.X] = 100

	snakes := game.Board.Snakes
	for _, snake := range snakes {
		if snake.ID == you.ID {
			continue
		}

		for _, coord := range snake.Body {
			matrix[coord.Y][coord.X] = math.MaxInt64
		}
		head := snake.Head
		matrix[head.Y][head.X] = math.MaxInt64
	}

	queue := []vertex{}
	currentVertex := vertex{head.Y, head.X}

	costTable[currentVertex] = route{cost: 101, previous: currentVertex}

	queue = append(queue, currentVertex)

	firstMove := true

	for len(queue) > 0 {
		currentVertex, queue = queue[len(queue)-1], queue[:len(queue)-1]
		currentRoute := costTable[currentVertex]
		visited[currentVertex] = true

		for _, dir := range directions {
			nextVertex := vertex{currentVertex[0] + dir[0], currentVertex[1] + dir[1]}
			if !valid(matrix, visited, size, nextVertex) {
				continue
			}
			if firstMove && nextVertex == tailVertex {
				continue
			}

			matrixCost := matrix[nextVertex[0]][nextVertex[1]]
			nextCost := currentRoute.cost + matrixCost
			if currentRoute.cost == 0 {
				fmt.Println(costTable)
				panic(2)
			}
			nextRoute, found := costTable[nextVertex]
			if !found {
				costTable[nextVertex] = route{cost: nextCost, previous: currentVertex}
			} else {

				if nextCost < nextRoute.cost {
					nextRoute.cost = nextCost
					nextRoute.previous = currentVertex
					costTable[nextVertex] = nextRoute

				} else {
					fmt.Println(costTable)
				}
			}

			queue = append(queue, nextVertex)
		}

		firstMove = false
	}

	// fmt.Println(matrix)
	for k, v := range costTable {
		fmt.Println(k, v)
	}

	headVertex := vertex{head.Y, head.X}
	path := []vertex{tailVertex}
	for tailVertex != headVertex {
		route, found := costTable[tailVertex]
		if !found {
			panic("not found")
		}
		path = append(path, route.previous)
		tailVertex = route.previous
	}
	fmt.Println("path from to", path, headVertex, tail)

	return models.MoveResponse{Move: models.Up, Shout: ""}
}

func valid(matrix [][]int, visited map[vertex]bool, size int, nextVertex vertex) bool {
	if nextVertex[0] < 0 || nextVertex[1] < 0 || nextVertex[0] >= size || nextVertex[1] >= size {
		return false
	}

	if visited[nextVertex] {
		return false
	}
	if matrix[nextVertex[0]][nextVertex[1]] == math.MaxInt64 {
		return false
	}

	return true
}
