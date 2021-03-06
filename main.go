package main

import (
"fmt"
"strconv"
"time"
"math/rand"
"bufio"
"os"
"Go-2048/console_util"
)

const (
	FIELD_HEIGHT = 6
	FIELD_WIDTH  = 6
	DIR_TOP      = 0
	DIR_RIGHT    = 1
	DIR_BOTTOM   = 2
	DIR_LEFT     = 3
	UP           = 'w'
	RIGHT        = 'd'
	DOWN         = 's'
	LEFT         = 'a'
	KEY_EXIT     = 'q'
)

type Field struct {
	matrix [FIELD_HEIGHT][FIELD_WIDTH]int
	score  int
	moved  bool
	over   bool
}

func clearScreen() {
	console_util.CallClear()
}

func closeGame() {
	fmt.Println("Bye bye")
	os.Exit(3)
	console_util.ExitProgram()
}

func isEdge(row int, col int) bool {
	if row == 0 || row == (FIELD_HEIGHT - 1) || col == 0 || col == (FIELD_WIDTH - 1) {
		return true
	}
	return false
}
func canMoveLeft(field Field, row int, col int) bool {
	if field.matrix[row][col] == field.matrix[row][col - 1] && field.matrix[row][col] != 0 {
		return true
	}
	return false
}

func canMoveRight(field Field, row int, col int) bool {
	if field.matrix[row][col] == field.matrix[row][col + 1] && field.matrix[row][col] != 0 {
		return true
	}
	return false
}
func canMoveUp(field Field, row int, col int) bool {
	if field.matrix[row][col] == field.matrix[row - 1][col] && field.matrix[row][col] != 0 {
		return true
	}
	return false
}
func canMoveDown(field Field, row int, col int) bool {
	if field.matrix[row][col] == field.matrix[row + 1][col] && field.matrix[row][col] != 0 {
		return true
	}
	return false
}

func unmovable(field Field, row int, col int, current int) bool {
	if current == 0 ||
	field.matrix[row][col - 1] == current ||
	field.matrix[row - 1][col] == current ||
	field.matrix[row][col + 1] == current ||
	field.matrix[row + 1][col] == current {
		return false
	}

	return true
}

func show(field Field) {
	console_util.PrintlnColored("Go, go!", console_util.COLOR_CYAN)
	console_util.PrintlnColored("(w - up | s - down | a - left | d - right | q - exit)\n", console_util.COLOR_WHITE)
	for row := 0; row < FIELD_HEIGHT; row++ {
		for col := 0; col < FIELD_WIDTH; col++ {
			if isEdge(row, col) {
				console_util.PrintColored("##  ",console_util.COLOR_RED)
			} else if field.matrix[row][col] == 0 {
				console_util.PrintNormal("__  ")
			} else {
				console_util.PrintNormal(fmt.Sprintf("%2s  ", strconv.Itoa(field.matrix[row][col])))
			}
		}
		fmt.Println()
	}
	console_util.PrintNormal("\nYour score: ")
	console_util.PrintlnColored(strconv.Itoa(field.score), console_util.COLOR_YELLOW)

	console_util.PrintColored("Yout turn: ", console_util.COLOR_CYAN)
	console_util.PrintColored("", console_util.COLOR_GREEN)
}

func over(field Field) bool {
	for row := 1; row < FIELD_HEIGHT; row++ {
		for col := 1; col < FIELD_WIDTH; col++ {
			current := field.matrix[row][col]
			if unmovable(field, row, col, current) == false {
				return false
			}
		}
	}
	return true
}

func fill(field *Field) {
	for row := 0; row < 6; row++ {
		for col := 0; col < 6; col++ {
			if row == 0 || row == 5 || col == 0 || col == 5 {
				(*field).matrix[row][col] = -1
			}
		}
	}

	generate(field)
	generate(field)
}

func countFreePositions(field *Field) []int {
	freePos := make([]int, 0, 16)
	for row := 1; row < 5; row++ {
		for col := 1; col < 5; col++ {
			if (*field).matrix[row][col] == 0 {
				freePos = append(freePos, (row - 1) * 4 + col)
			}
		}
	}
	return freePos
}

func generateCoord(field *Field) (int, int) {
	freePos := countFreePositions(field)
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	position := random.Intn(len(freePos))
	coord := freePos[position]

	col := coord % 4
	if col == 0 {
		col = 4
	}

	row := ((coord - col) / 4) + 1
	return row, col

}

func generateValue() int {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	value := random.Intn(4)
	if value == 0 {
		return 4
	}
	return 2
}

func generate(field *Field) {
	row, col := generateCoord(field)
	(*field).matrix[row][col] = generateValue()
}

func (field *Field)left() {
	for row := 1; row < 5; row++ {
		for col := 2; col < 5; col++ {
			if field.matrix[row][col] > 0 {
				for field.matrix[row][col - 1] == 0 {
					field.matrix[row][col - 1] = field.matrix[row][col]
					field.matrix[row][col] = 0
					col--
					field.moved = true
				}
			}
		}
		for j := 2; j < 5; j++ {
			if canMoveLeft(*field, row, j) {
				field.matrix[row][j - 1] = 2 * field.matrix[row][j - 1]
				field.score += field.matrix[row][j - 1]
				field.moved = true
				for k := j; k < 4; k++ {
					field.matrix[row][k] = field.matrix[row][k + 1]
				}
				field.matrix[row][4] = 0
			}
		}
	}
}
func (field *Field)right() {
	for i := 1; i < 5; i++ {
		for j := 3; j > 0; j-- {
			if field.matrix[i][j] > 0 {
				for field.matrix[i][j + 1] == 0 {
					field.matrix[i][j + 1] = field.matrix[i][j]
					field.matrix[i][j] = 0
					j++
					field.moved = true
				}
			}
		}
		for j := 3; j > 0; j-- {
			if canMoveRight(*field, i, j) {
				field.matrix[i][j + 1] = 2 * field.matrix[i][j + 1]
				field.score += field.matrix[i][j + 1]
				field.moved = true
				for k := j; k > 1; k-- {
					field.matrix[i][k] = field.matrix[i][k - 1]
				}
				field.matrix[i][1] = 0
			}
		}
	}
}
func (field *Field)up() {
	for j := 1; j < 5; j++ {
		for i := 2; i < 5; i++ {
			if field.matrix[i][j] > 0 {
				for field.matrix[i - 1][j] == 0 {
					field.matrix[i - 1][j] = field.matrix[i][j]
					field.matrix[i][j] = 0
					i--
					field.moved = true
				}
			}
		}
		for i := 2; i < 5; i++ {
			if canMoveUp(*field, i, j) {
				field.matrix[i - 1][j] = 2 * field.matrix[i - 1][j]
				field.score += field.matrix[i - 1][j]
				field.moved = true
				for k := i; k < 4; k++ {
					field.matrix[k][j] = field.matrix[k + 1][j]
				}
				field.matrix[4][j] = 0
			}
		}
	}
}

func (field *Field)down() {
	for j := 1; j < 5; j++ {
		for i := 3; i > 0; i-- {
			if field.matrix[i][j] > 0 {
				for field.matrix[i + 1][j] == 0 {
					field.matrix[i + 1][j] = field.matrix[i][j]
					field.matrix[i][j] = 0
					i++
					field.moved = true
				}
			}
		}
		for i := 3; i > 0; i-- {
			if canMoveDown(*field, i, j) {
				field.matrix[i + 1][j] = 2 * field.matrix[i + 1][j]
				field.score += field.matrix[i + 1][j]
				field.moved = true
				for k := i; k > 1; k-- {
					field.matrix[k][j] = field.matrix[k - 1][j]
				}
				field.matrix[1][j] = 0
			}
		}
	}
}

func direct(field *Field, direction int) {
	switch direction {
	case DIR_TOP:
		field.up()
	case DIR_RIGHT:
		field.right()
	case DIR_BOTTOM:
		field.down()
	case DIR_LEFT:
		field.left()
	default:

	}

}

func selectDirection() int {
	//Works only with enter!
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	symb := rune([]byte(input)[0])
	switch symb {
	case UP:
		return DIR_TOP
	case RIGHT:
		return DIR_RIGHT
	case DOWN:
		return DIR_BOTTOM
	case LEFT:
		return DIR_LEFT
	case KEY_EXIT:
		closeGame()
		return -1
	default:
		return -1
	}
}

func main() {
	var x Field
	fill(&x)
	clearScreen()
	show(x)
	for x.over == false {
		direct(&x, selectDirection())
		if x.moved == true {
			generate(&x)
			x.moved = false
		}
		clearScreen()
		show(x)
	}
}
