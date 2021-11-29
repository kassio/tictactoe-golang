package main

import (
	"fmt"
	"strconv"
)

type Winner struct {
	found bool
	msg   string
}
type board [9]string

func printBoard(values board) {
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

func findWinnerXY(values board, title string, posCalc func(int, int) int) Winner {
	w := Winner{}
	var player string
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			value := values[posCalc(i, j)]
			if player == "" {
				player = value
			} else if player != value {
				player = ""
				break
			}
		}

		if player != "" {
			break
		}
	}

	if player != "" {
		w.msg = fmt.Sprintf("%s wins the game by %s", player, title)
		w.found = true
	}

	return w
}

func findWinnerDiagonal(values board, posCalc func(int) int) Winner {
	w := Winner{}
	var player string
	for i := 0; i < 3; i++ {
		value := values[posCalc(i)]
		if player == "" {
			player = value
		} else if player != value {
			player = ""
			break
		}
	}

	if player != "" {
		w.msg = fmt.Sprintf("%s wins the game by diagonal", player)
		w.found = true
	}

	return w
}

func findByRow(values board) Winner {
	return findWinnerXY(values, "row", func(i, j int) int {
		return (i * 3) + j
	})
}

func findByColumn(values board) Winner {
	return findWinnerXY(values, "column", func(i, j int) int {
		return (j * 3) + i
	})
}

func findByDiagonal1(values board) Winner {
	return findWinnerDiagonal(values, func(i int) int {
		return (i * 3) + i
	})
}
func findByDiagonal2(values board) Winner {
	return findWinnerDiagonal(values, func(i int) int {
		return (i * 3) + (2 - i)
	})
}

func findWinner(values board) Winner {
	finders := make(chan Winner)

	go func(values board) { finders <- findByRow(values) }(values)
	go func(values board) { finders <- findByColumn(values) }(values)
	go func(values board) { finders <- findByDiagonal1(values) }(values)
	go func(values board) { finders <- findByDiagonal2(values) }(values)

	winner := Winner{}
	for i := 0; i < 4; i++ {
		f := <-finders
		if f.found {
			winner = f
		}
	}

	return winner
}

func validateInput(in int, values board) (bool, string) {
	switch {
	case in < 0 || in > 8:
		return false, "Invalid position! Choose between 1 and 9"
	case values[in] == "x" || values[in] == "o":
		return false, "Position already taken!"
	default:
		return true, ""
	}
}

func main() {
	var players = [2]string{"x", "o"}

	var values board
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

		valid, errMsg := validateInput(position, values)
		if !valid {
			i--
			fmt.Println(errMsg, "Try again")
			continue
		}

		values[position] = player
		printBoard(values)

		// it needs at least five rounds to have a winner
		if i >= 4 {
			winner := findWinner(values)
			if winner.found {
				fmt.Printf(winner.msg)
				return
			}
		}
	}

	fmt.Println(
		`No one win,
		no one lost,
		at the end
		we're all winners
		who have lost`)
}
