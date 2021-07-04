package ai

import (
	"curls/models"
)

var directions = []vertex{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

const (
	costBody  = 100_000
	costTail  = 200
	costSpace = 100
)

type vertex [2]int

type route struct {
	cost     int
	previous vertex
}

type matrix [][]int

func newMatrix(size int) matrix {
	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int, size)
	}
	for y := range matrix {
		for x := range matrix[0] {
			matrix[y][x] = costSpace
		}
	}

	return matrix
}

func (m *matrix) set(y, x, v int) {
	size := len(*m) - 1
	(*m)[size-y][x] = v
}

func (m *matrix) get(y, x int) int {
	size := len(*m) - 1
	return (*m)[size-y][x]
}

func (m *matrix) addSnake(snake models.Battlesnake) {
	for _, coord := range snake.Body {
		m.set(coord.Y, coord.X, costBody)
	}
	head := snake.Head
	m.set(head.Y, head.X, costBody)
}

func (m *matrix) addTail(you models.Battlesnake) {
	tail := you.Body[len(you.Body)-1]
	m.set(tail.Y, tail.X, costTail)
}

func NextMove(game models.GameRequest) models.MoveResponse {
	size := game.Board.Height
	boardMatrix := newMatrix(size)

	for _, snake := range game.Board.Snakes {
		boardMatrix.addSnake(snake)
	}
	you := game.You
	head := you.Head
	boardMatrix.addTail(you)

	tail := you.Body[len(you.Body)-1]
	tailVertex := vertex{tail.Y, tail.X}
	headVertex := vertex{head.Y, head.X}

	costTable := buildCostTable(boardMatrix, size, headVertex, tailVertex)

	moveVertex := moveTo(costTable, tailVertex)

	direction := getDirection(headVertex, moveVertex)
	return models.MoveResponse{Move: direction, Shout: ""}
}

func buildCostTable(boardMatrix matrix, size int, headVertex, tailVertex vertex) map[vertex]route {
	queue := []vertex{}
	currentVertex := headVertex

	costTable := map[vertex]route{}
	costTable[currentVertex] = route{cost: 101, previous: currentVertex}
	visited := map[vertex]bool{}

	queue = append(queue, currentVertex)

	firstMove := true
	for len(queue) > 0 {
		currentVertex, queue = queue[len(queue)-1], queue[:len(queue)-1]
		currentRoute := costTable[currentVertex]

		if visited[currentVertex] {
			continue
		}
		visited[currentVertex] = true

		for _, dir := range directions {
			nextVertex := vertex{currentVertex[0] + dir[0], currentVertex[1] + dir[1]}
			if (firstMove && nextVertex == tailVertex) || !valid(boardMatrix, size, nextVertex) {
				continue
			}

			nextCost := currentRoute.cost + boardMatrix.get(nextVertex[0], nextVertex[1])
			nextRoute, found := costTable[nextVertex]
			if !found {
				costTable[nextVertex] = route{cost: nextCost, previous: currentVertex}
			} else if nextCost < nextRoute.cost {
				nextRoute.cost = nextCost
				nextRoute.previous = currentVertex
				costTable[nextVertex] = nextRoute
			}
			queue = append(queue, nextVertex)
		}

		firstMove = false
	}

	return costTable
}

func moveTo(costTable map[vertex]route, v vertex) vertex {
	path := pathTo(costTable, v)
	return path[len(path)-2] // most recent element is the head
}

func pathTo(costTable map[vertex]route, v vertex) []vertex {
	current := v
	path := []vertex{current}
	for {
		route := costTable[current]
		if current == route.previous {
			break
		}
		path = append(path, route.previous)
		current = route.previous
	}

	return path
}

func valid(matrix [][]int, size int, nextVertex vertex) bool {
	if nextVertex[0] < 0 || nextVertex[1] < 0 || nextVertex[0] >= size || nextVertex[1] >= size {
		return false
	}

	if matrix[nextVertex[0]][nextVertex[1]] >= costBody {
		return false
	}

	return true
}

func getDirection(headVertex, moveVertex vertex) models.Move {
	if moveVertex[0] < headVertex[0] {
		return models.Down
	} else if moveVertex[0] > headVertex[0] {
		return models.Up
	} else if moveVertex[1] < headVertex[1] {
		return models.Left
	} else {
		return models.Right
	}
}
