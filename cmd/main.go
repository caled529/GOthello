package main

import "fmt"

type Vec struct {
	x int
	y int
}

type Tile int

const (
	Empty Tile = iota
	Black
	White
	Legal
)

func main() {
	board := [8][8]Tile{}
	board[3][3] = White
	board[3][4] = Black
	board[4][3] = Black
	board[4][4] = White
	currentPlayer := Black
	populateLegalMoves(&board, currentPlayer)
	fmt.Println(boardToString(&board))
	currentPlayer = placeTile(Vec{3, 2}, &board, currentPlayer)
	populateLegalMoves(&board, currentPlayer)
	fmt.Println(boardToString(&board))
	currentPlayer = placeTile(Vec{2, 2}, &board, currentPlayer)
	populateLegalMoves(&board, currentPlayer)
	fmt.Println(boardToString(&board))
}

func populateLegalMoves(board *[8][8]Tile, currentPlayer Tile) {
	oppositePlayer := Black
	if currentPlayer == Black {
		oppositePlayer = White
	}
	var legalMoves []Vec
	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[x]); y++ {
			if board[x][y] != Empty {
				continue
			}
			pos := Vec{x, y}
			for _, move := range circleTileCheck(pos, board, oppositePlayer) {
				dir := Vec{move.x - x, move.y - y}
				if lineTileCheck(pos, dir, board, currentPlayer) {
					legalMoves = append(legalMoves, Vec{x, y})
				}
			}
		}
	}
	for _, move := range legalMoves {
		board[move.x][move.y] = Legal
	}
}

func circleTileCheck(pos Vec, board *[8][8]Tile, target Tile) []Vec {
	var matches []Vec
	for i := max(0, pos.x-1); i < min(len(board), pos.x+2); i++ {
		for j := max(0, pos.y-1); j < min(len(board[i]), pos.y+2); j++ {
			if board[i][j] == target {
				matches = append(matches, Vec{i, j})
			}
		}
	}
	return matches
}

func lineTileCheck(move Vec, dir Vec, board *[8][8]Tile, target Tile) bool {
	if dir == (Vec{0, 0}) || len(board) == 0 {
		return false
	}
	pos := Vec{move.x + dir.x, move.y + dir.y}
	// Doesn't count if tiles of same colour are adjacent
	if board[pos.x][pos.y] == target {
		return false
	}
	for {
		pos = Vec{pos.x + dir.x, pos.y + dir.y}
		if pos.x < 0 || pos.x >= len(board) || pos.y < 0 || pos.y >= len(board[0]) {
			return false
		}
		if board[pos.x][pos.y] == target {
			return true
		}
	}
}

func placeTile(move Vec, board *[8][8]Tile, currentPlayer Tile) Tile {
	illegalMove := true
	for i, row := range board {
		for j, tile := range row {
			if tile == Legal && move == (Vec{i, j}) {
				illegalMove = false
				break
			}
		}
	}
	if illegalMove {
		return currentPlayer
	}
	oppositePlayer := Black
	if currentPlayer == Black {
		oppositePlayer = White
	}
	board[move.x][move.y] = currentPlayer
	for _, opp := range circleTileCheck(move, board, oppositePlayer) {
		dir := Vec{opp.x - move.x, opp.y - move.y}
		if lineTileCheck(move, dir, board, currentPlayer) {
			flipWalk(move, dir, board, currentPlayer)
		}
	}
	for i, row := range board {
		for j, tile := range row {
			if tile == Legal {
				board[i][j] = Empty
			}
		}
	}
	return oppositePlayer
}

func flipWalk(move Vec, dir Vec, board *[8][8]Tile, currentPlayer Tile) {
	pos := move
	for {
		pos = Vec{pos.x + dir.x, pos.y + dir.y}
		if board[pos.x][pos.y] != currentPlayer {
			board[pos.x][pos.y] = currentPlayer
		} else {
			break
		}
	}
}

func boardToString(board *[8][8]Tile) string {
	var boardString string
	for j := 0; j < len(board[0]); j++ {
		for i := 0; i < len(board); i++ {
			switch board[i][j] {
			case Empty:
				boardString += "┼"
			case Black:
				boardString += "○"
			case White:
				boardString += "●"
			case Legal:
				boardString += "?"
			}
		}
		boardString += "\n"
	}
	return boardString
}
