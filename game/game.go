package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/josephthomashines/TicTacGo/board"
)

// Errors
/*
type BoardOutOfBoundsError struct {
	received, limit int
	axis            string
}

func (e *BoardOutOfBoundsError) Error() string {
	out := fmt.Sprintf("Error: Index %d must be %s to %d on the %s axis", e.received, "%s", e.limit, e.axis)

	if e.received > e.limit {
		out = fmt.Sprintf(out, "<=")
	}

	if e.received < e.limit {
		out = fmt.Sprintf(out, ">=")
	}

	return out

}
*/

// Core
type Game struct {
	board board.Board
	turn  string
	move  int
}

func (g *Game) String() string {
	out := "\n"

	out += fmt.Sprintf("%s's Turn\n%s\n", g.turn, &g.board)

	return out
}

func (g *Game) SwitchTurn() {
	if g.turn == "X" {
		g.turn = "O"
	} else {
		g.turn = "X"
	}
}

func Init() Game {
	return Game{board.Init(3, 3), "X", 0}
}

func ThreeQuals(a, b, c string) bool {
	if a == b && b == c && a != board.PLACEHOLDER {
		return true
	}
	return false
}

func IsWin(g *Game) string {
	b := g.board.Get()

	// Horizontal
	for _, row := range b {
		if ThreeQuals(row[0], row[1], row[2]) {
			return row[0]
		}
	}

	// Vertical
	for j := range b[0] {
		if ThreeQuals(b[0][j], b[1][j], b[2][j]) {
			return b[0][j]
		}
	}

	// Diagonal 1

	if ThreeQuals(b[0][0], b[1][1], b[2][2]) {
		return b[0][0]
	}

	// Diagonal 2

	if ThreeQuals(b[0][2], b[1][1], b[2][0]) {
		return b[0][2]
	}

	draw := true
	for _, row := range b {
		for _, val := range row {
			if val == board.PLACEHOLDER {
				draw = false
			}
		}
	}

	if draw {
		return "-"
	}

	return board.PLACEHOLDER
}

func MapChoice(c int) (int, int) {
	if c == 9 {
		return 0, 2
	} else if c == 8 {
		return 0, 1
	} else if c == 7 {
		return 0, 0
	} else if c == 6 {
		return 1, 2
	} else if c == 5 {
		return 1, 1
	} else if c == 4 {
		return 1, 0
	} else if c == 3 {
		return 2, 2
	} else if c == 2 {
		return 2, 1
	} else if c == 1 {
		return 2, 0
	}
	return -1, -1
}

func Clear() {
	// Only works in Unix shells
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

}

// Read console input, let user play
func Play() {
	game := Init()
	gp := &game

	choice, x, y := -1, -1, -1

	Clear()
	fmt.Println()

	fmt.Println("TicTacGo: A go learning exercise by Joseph Hines\n")
	fmt.Println("To play, push a number followed by \nthe Enter key to select a spot \nto play the current character\n")
	fmt.Println("7 8 9\n4 5 6\n1 2 3")

	for IsWin(gp) == board.PLACEHOLDER {
		fmt.Print(gp)

		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
		}

		Clear()

		if int(char) < 58 && int(char) > 48 {
			choice = int(char) - 48
			x, y = MapChoice(choice)

			if gp.board.Get()[x][y] == board.PLACEHOLDER {
				gp.board.Set(x, y, gp.turn)
				gp.SwitchTurn()
				fmt.Println()
			} else {
				fmt.Println("Bad value: space already played")
			}
		} else {
			fmt.Println("Bad value: must be 1-9")
		}

	}

	gp.SwitchTurn()
	fmt.Print(gp)

	winner := IsWin(gp)

	if winner == "-" {
		fmt.Println("Draw! Everybody loses!")
	} else {
		fmt.Printf("%s wins!\n", winner)
	}
}

// Play n games randomly
// Report win rate of each player
// Average number of moves
type Results map[string]int

func PlaySim() Results {
	rand.Seed(time.Now().UTC().UnixNano())

	stats := make(Results)
	stats["O"] = 0
	stats["X"] = 0
	stats["D"] = 0
	stats["plays"] = 0

	game := Init()
	gp := &game

	choice, x, y, plays := -1, -1, -1, 0

	for IsWin(gp) == board.PLACEHOLDER {
		choice = rand.Intn(9) + 1
		x, y = MapChoice(choice)

		if gp.board.Get()[x][y] == board.PLACEHOLDER {
			gp.board.Set(x, y, gp.turn)
			gp.SwitchTurn()
			plays++
		}
	}

	gp.SwitchTurn()

	winner := IsWin(gp)

	if winner == "-" {
		stats["D"] = 1
	} else {
		stats[winner] = 1
	}

	stats["plays"] = plays

	return stats
}

func RunSims(n int) chan Results {
	c := make(chan Results)
	go func() {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				c <- PlaySim()
			}()
		}
		wg.Wait()
		close(c)
	}()
	return c
}

func Simulate(n int) {
	stats := make(Results)
	stats["O"] = 0
	stats["X"] = 0
	stats["D"] = 0
	stats["plays"] = 0

	start := time.Now()

	for v := range RunSims(n) {
		//		fmt.Println(v)
		stats["O"] += v["O"]
		stats["X"] += v["X"]
		stats["D"] += v["D"]
		stats["plays"] += v["plays"]
	}

	elapsed := time.Since(start)

	Clear()

	fmt.Printf("\nTicTacGo: Simulated %d Games (played randomly, X goes first)\n\n Took: %s\n\n", n, elapsed)

	fmt.Println("Totals:")
	fmt.Printf("  X wins: %d\n", stats["X"])
	fmt.Printf("  O wins: %d\n", stats["O"])
	fmt.Printf("  Draws: %d\n", stats["D"])
	fmt.Printf("  Plays: %d\n", stats["plays"])

	fmt.Println("Rates:")
	fmt.Printf("  X winrate: %.2f%%\n", float64(stats["X"])/float64(n)*100.0)
	fmt.Printf("  O winrate: %.2f%%\n", float64(stats["O"])/float64(n)*100.0)
	fmt.Printf("  Draw rate: %.2f%%\n", float64(stats["D"])/float64(n)*100.0)
	fmt.Printf("  Plays per game: %.2f\n", float64(stats["plays"])/float64(n))
}
