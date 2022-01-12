package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CELL_SIZE = 10
	WIDTH     = 500
	HEIGHT    = 500
)

type Game struct {
	grid   [][]Cell
	buffer [][]Cell
}

type Cell struct {
	is_alive   uint8
	generation uint8
}

var (
	grid_width               = int(math.Round(WIDTH / CELL_SIZE))
	grid_height              = int(math.Round(HEIGHT / CELL_SIZE))
	white                    = color.RGBA{255, 255, 255, 255}
	black       color.RGBA   = color.RGBA{0, 0, 0, 255}
	red         color.RGBA   = color.RGBA{255, 0, 0, 255}
	orange      color.RGBA   = color.RGBA{255, 115, 28, 255}
	yellow      color.RGBA   = color.RGBA{251, 255, 0, 255}
	green       color.RGBA   = color.RGBA{0, 255, 0, 255}
	blue        color.RGBA   = color.RGBA{0, 0, 255, 255}
	purple      color.RGBA   = color.RGBA{255, 0, 255, 255}
	colors      []color.RGBA = []color.RGBA{white, red, orange, yellow, green, blue, purple}
	params                   = map[string]uint8{
		"born":          3,
		"min_neighbors": 2,
		"max_neighbors": 3,
	}
)

func (g *Game) Update() error {
	for i := 1; i < grid_width-1; i++ {
		for j := 1; j < grid_height-1; j++ {
			g.buffer[i][j].is_alive = 0

			neighbours := g.grid[i-1][j-1].is_alive + g.grid[i-1][j].is_alive + g.grid[i-1][j+1].is_alive + g.grid[i][j-1].is_alive + g.grid[i][j+1].is_alive + g.grid[i+1][j-1].is_alive + g.grid[i+1][j].is_alive + g.grid[i+1][j+1].is_alive

			if g.grid[i][j].is_alive == 0 && neighbours == params["born"] {
				g.buffer[i][j].is_alive = 1
				g.buffer[i][j].generation++
			} else if neighbours < params["min_neighbors"] || neighbours > params["max_neighbors"] {
				g.buffer[i][j].is_alive = 0
			} else {
				g.buffer[i][j] = g.grid[i][j]
			}
		}
	}
	temp := g.buffer
	g.buffer = g.grid
	g.grid = temp
	return nil
}

func get_color(time uint8) color.RGBA {
	index := math.Round(float64(time) / 255 * 7)
	return colors[int(math.Min(index, 6))]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(black)
	for i := 0; i < grid_width; i++ {
		for j := 0; j < grid_height; j++ {
			if g.grid[i][j].is_alive >= 1 {
				for x := 0; x < CELL_SIZE; x++ {
					for y := 0; y < CELL_SIZE; y++ {
						color := get_color(g.grid[i][j].generation)
						screen.Set((i*CELL_SIZE)+x, (j*CELL_SIZE)+y, color)
					}
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func (g *Game) init_game() {
	grid := make([][]Cell, grid_width)
	buffer := make([][]Cell, grid_width)
	for i := range grid {
		grid[i] = make([]Cell, grid_height)
		buffer[i] = make([]Cell, grid_height)
	}
	g.grid = grid
	g.buffer = buffer
}

func (g *Game) random_generation() {
	for x := 1; x < grid_width-1; x++ {
		for y := 1; y < grid_height-1; y++ {
			if rand.Float32() < 0.5 {
				g.grid[x][y].is_alive = 1
				g.grid[x][y].generation = 1
			} else {
				g.grid[x][y].is_alive = 0
				g.grid[x][y].generation = 0
			}
		}
	}
}

func main() {
	game := &Game{}
	game.init_game()
	game.random_generation()

	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Conway's Game of Life")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
