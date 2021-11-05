package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/danicat/simpleansi"
)

type spriter struct {
	row int
	col int
}

var player spriter
var ghosts []*spriter
var maze []string

var score int // 分数
var numDots int
var lives = 1

func main() {
	// initiallize the game
	initialize()
	defer cleanup()
	// load resources
	err := loadMaze("maze01.txt")
	if err != nil {
		log.Println("Failed to load maze: ", err)
		return
	}
	// game loop
	for {
		// update screen
		printScreen()
		// process input
		input, err := readInput()
		if err != nil {
			log.Println("Failed to read input: ", err)
			break
		}
		// process movement
		movePlayer(input)
		moveGhosts()

		// process collisions
		for _, g := range ghosts {
			if player == *g {
				lives = 0
			}
		}
		// check game over
		if input == "ESC" || numDots == 0 || lives <= 0 {
			break
		}
		// break
	}
}
func loadMaze(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}
	for row, line := range maze {
		for col, char := range line {
			switch char {
			case 'P':
				player = spriter{row, col} // 定位玩家的位置
			case 'G':
				ghosts = append(ghosts, &spriter{row, col}) // 鬼的位置
			case '.':
				numDots++
			}
		}
	}
	return nil
}

func printScreen() {
	simpleansi.ClearScreen()
	for _, line := range maze {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Printf("%c", chr)
				fallthrough
			case '.':
				fmt.Printf("%c", chr)
			default:
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	simpleansi.MoveCursor(player.row, player.col)
	fmt.Print("P")
	for _, g := range ghosts {
		simpleansi.MoveCursor(g.row, g.col)
		fmt.Print("G")
	}
	simpleansi.MoveCursor(len(maze)+1, 0)
	fmt.Println("Score:", score, "\tLives:", lives)
}

func initialize() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("Failed to set terminal to cbreak mode: ", err)
	}
}

func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("Unable to restore cooked mode: ", err)
	}
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}
	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}
	return "", nil
}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze) {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	}

	if maze[newRow][newCol] == '#' { // todo index 溢出问题
		newRow = oldRow
		newCol = oldCol
	}

	return
}

func movePlayer(dir string) {
	// 移动玩家位置
	player.row, player.col = makeMove(player.row, player.col, dir)
	switch maze[player.row][player.col] {
	case '.':
		numDots--
		score++
		maze[player.row] = maze[player.row][0:player.col] + " " + maze[player.row][player.col+1:]
	}
}

func drawDirection() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func moveGhosts() {
	for _, g := range ghosts {
		dir := drawDirection()
		g.row, g.col = makeMove(g.row, g.col, dir)
	}
}
