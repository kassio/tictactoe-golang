package main

import (
	"fmt"
	"strconv"
)

var players = [2]string{"x", "o"}

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

func winnerByRow(positions [9]string) (bool, string) {
	winner := ""
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			value := positions[(row*3)+col]
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

	return winner != "", winner
}

func winnerByCol(positions [9]string) (bool, string) {
	winner := ""
	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			value := positions[(row*3)+col]
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

	return winner != "", winner
}

func winnerByDiagonal(positions [9]string) (bool, string) {
	winner := ""
	for i := 0; i < 3; i++ {
		value := positions[(i*3)+i]
		if winner == "" {
			winner = value
		} else if winner != value {
			winner = ""
			break
		}
	}

	if winner != "" {
		return winner != "", winner
	}

	for i := 0; i < 3; i++ {
		value := positions[(i*3)+(2-i)]
		if winner == "" {
			winner = value
		} else if winner != value {
			winner = ""
			break
		}
	}

	return winner != "", winner
}

func main() {
	var playedPositions [9]bool
	var values [9]string
	for i := range values {
		values[i] = strconv.Itoa(i + 1)
	}
	fmt.Println(">", playedPositions)

	printBoard(values)
	for i := 0; i < 9; i++ {
		player := players[i%2]
		fmt.Printf("Player %s chose a position: ", player)

		var position int
		fmt.Scanf("%d", &position)
		position-- // user inputs the position indexed by 1

		if playedPositions[position] {
			i--
			fmt.Println("Position already taken! Try again")
			continue
		}

		playedPositions[position] = true
		values[position] = player
		printBoard(values)

		ended, winner := (winnerByRow(values))
		if ended {
			fmt.Printf("%s wins the game by row", winner)
			return
		}

		ended, winner = (winnerByCol(values))
		if ended {
			fmt.Printf("%s wins the game by column", winner)
			return
		}

		ended, winner = (winnerByDiagonal(values))
		if ended {
			fmt.Printf("%s wins the game by diagonal", winner)
			return
		}
	}

	fmt.Println(`No one win,
	no one lost,
	at the end
	we're all winners
	who have lost`)
}
