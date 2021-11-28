package main

import (
	"fmt"
	"strconv"
)

var players = [2]string{"x", "o"}

type Winner struct {
	found bool
	msg   string
}

func printBoard(values [9]string) {
	fmt.Print("\033[H\033[2J") // clears the screen

	for i, value := range values {
		fmt.Printf(" %-1s ", value)

		if i > 7 {
			fmt.Print("\n")
		} else if (i+1)%3 == 0 {
			fmt.Print("\n───┼───┼───\n")
		} else {
			fmt.Print("│")
		}
	}
}

func findWinnerXY(values [9]string, title string, posCalc func(int, int) int) <-chan Winner {
	wc := make(chan Winner)
	go func() {
		var winner string
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				value := values[posCalc(i, j)]
				if winner == "" {
					winner = value
				} else if winner != value {
					winner = ""
					break
				}
			}

			if winner != "" {
				break
			}
		}

		if winner != "" {
			wc <- Winner{msg: fmt.Sprintf("%s wins the game by %s", winner, title), found: true}
		}

		wc <- Winner{}
	}()

	return wc
}

func findWinnerDiagonal(values [9]string, title string, posCalc func(int) int) <-chan Winner {
	wc := make(chan Winner)
	go func() {
		var winner string
		for i := 0; i < 3; i++ {
			value := values[posCalc(i)]
			if winner == "" {
				winner = value
			} else if winner != value {
				winner = ""
				break
			}
		}

		if winner != "" {
			wc <- Winner{msg: fmt.Sprintf("%s wins the game by %s", winner, title), found: true}
		}

		wc <- Winner{}
	}()

	return wc
}

func findWinner(values [9]string) (bool, string) {
	finders := make(chan Winner)

	go func() {
		found := <-findWinnerXY(values, "row", func(i, j int) int {
			return (i * 3) + j
		})
		// fmt.Println("row", found)
		finders <- found
	}()

	go func() {
		found := <-findWinnerXY(values, "column", func(i, j int) int {
			return (j * 3) + i
		})
		// fmt.Println("column", found)
		finders <- found
	}()

	go func() {
		found := <-findWinnerDiagonal(values, "diagonal", func(i int) int {
			return (i * 3) + i
		})
		// fmt.Println("diagonal1", found)
		finders <- found
	}()

	go func() {
		found := <-findWinnerDiagonal(values, "diagonal", func(i int) int {
			return (i * 3) + (2 - i)
		})
		// fmt.Println("diagonal2", found)
		finders <- found
	}()

	winner := Winner{}
	for i := 0; i < 4; i++ {
		f := <-finders
		if f.found {
			winner = f
		}
	}

	return winner.found, winner.msg
}

func validateInput(in int, playedPositions [9]bool) (bool, string) {
	switch {
	case in < 0 || in > 8:
		return false, "Invalid position! Choose between 1 and 9"
	case playedPositions[in]:
		return false, "Position already taken!"
	default:
		return true, ""
	}
}

func main() {
	var playedPositions [9]bool
	var values [9]string
	for i := range values {
		values[i] = strconv.Itoa(i + 1)
	}

	printBoard(values)
	for i := 0; i < 9; i++ {
		player := players[i%2]
		fmt.Printf("Player %s chose a position: ", player)

		var position int
		fmt.Scanf("%d", &position)
		position-- // user inputs the position indexed by 1

		valid, errMsg := validateInput(position, playedPositions)
		if !valid {
			i--
			fmt.Println(errMsg, "Try again")
			continue
		}

		playedPositions[position] = true
		values[position] = player
		printBoard(values)

		found, winnerMsg := findWinner(values)
		if found {
			fmt.Printf(winnerMsg)
			return
		}
	}

	fmt.Println(
		`No one win,
		no one lost,
		at the end
		we're all winners
		who have lost`)
}
