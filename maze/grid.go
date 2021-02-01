// Copyright ©️ 2021 oddstream.games

package maze

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TilesAcross and TilesDown are package-level variables so they can be seen by Tile
var (
	TilesAcross int
	TilesDown   int
	TileSize    int
	CameraX     float64
	CameraY     float64
)

// Grid is an object representing the grid of tiles
type Grid struct {
	tiles           []*Tile // a slice (not array!) of pointers to Tile objects
	colorBackground color.RGBA
	colorWall       color.RGBA
	input           *Input
}

// NewGrid create a Grid object
func NewGrid(w, h int) *Grid {

	// var screenWidth, screenHeight int

	// if runtime.GOARCH == "wasm" {
	// 	screenWidth, screenHeight = WindowWidth, WindowHeight
	// } else {
	// 	screenWidth, screenHeight = ebiten.WindowSize()
	// }

	TileSize = 80
	TilesAcross, TilesDown = w, h

	g := &Grid{tiles: make([]*Tile, TilesAcross*TilesDown)}
	for i := range g.tiles {
		g.tiles[i] = NewTile(i%TilesAcross, i/TilesAcross)
	}

	// link the tiles together to avoid all that tedious 2d array stuff
	for _, t := range g.tiles {
		x := t.X
		y := t.Y
		t.edges[0] = g.findTile(x, y-1)
		t.edges[1] = g.findTile(x+1, y)
		t.edges[2] = g.findTile(x, y+1)
		t.edges[3] = g.findTile(x-1, y)
	}

	InitTile()

	g.CreateNextLevel()

	g.input = NewInput()

	return g
}

// Size returns the size of the grid in pixels
func (g *Grid) Size() (int, int) {
	return TilesAcross * TileSize, TilesDown * TileSize
}

func (g *Grid) findTile(x, y int) *Tile {

	// we re-order the tiles when dragging, to put the dragged tile at the top of the z-order
	// so we can't use i := x + (y * TilesAcross) to find index of tile in slice
	// except if this func is just used after tiles have been created

	// for _, t := range g.tiles {
	// 	if t.X == x && t.Y == y {
	// 		return t
	// 	}
	// }
	// return nil

	if x < 0 || x >= TilesAcross {
		return nil
	}
	if y < 0 || y >= TilesDown {
		return nil
	}
	i := x + (y * TilesAcross)
	if i < 0 || i > len(g.tiles) {
		log.Fatal("findTile index out of bounds")
	}
	return g.tiles[i]
}

// findTileAt finds the tile under the mouse click or touch
func (g *Grid) findTileAt(pt image.Point) *Tile {
	for _, t := range g.tiles {
		if InRect(pt, t.Rect) {
			return t
		}
	}
	return nil
}

func (g *Grid) randomTile() *Tile {
	i := rand.Intn(len(g.tiles))
	return g.tiles[i]
}

func (g *Grid) createRoom(x1, y1 int) {
	x2 := x1 + 4
	y2 := y1 + 4
	if x2 >= TilesAcross || y2 >= TilesDown {
		return
	}
	for x := x1; x < x2; x++ {
		for y := y1; y < y2; y++ {
			t := g.findTile(x, y)
			t.removeAllWalls()
		}
	}
	// for x := x1; x < x2; x++ {
	// 	for y := y1; y < y2; y++ {
	// 		t := g.findTile(x, y)
	// 		t.addWall(0)
	// 		t.addWall(1)
	// 		t.addWall(2)
	// 		t.addWall(3)
	// 	}
	// }
	for x := x1; x < x2; x++ {
		t := g.findTile(x, y1)
		t.addWall(0)
		t = g.findTile(x, y2-1)
		t.addWall(2)
	}
	for y := y1; y < y2; y++ {
		t := g.findTile(x1, y)
		t.addWall(3)
		t = g.findTile(x2-1, y)
		t.addWall(1)
	}
}

func (g *Grid) createRooms() {
	for x := 0; x < TilesAcross; x += 4 {
		for y := 0; y < TilesDown; y += 4 {
			if rand.Float64() < 0.5 {
				g.createRoom(x, y)
			}
		}
	}
}

func (g *Grid) carve() {
	t := g.randomTile()
	t.recursiveBacktracker()
}

func (g *Grid) fillCulDeSacs() {
	for _, t := range g.tiles {
		t.fillCulDeSac()
	}
}

// AllTiles applies a func to all tiles
func (g *Grid) AllTiles(fn func(*Tile)) {
	for _, t := range g.tiles {
		fn(t)
	}
}

// CreateNextLevel resets game data and moves the puzzle to the next level
func (g *Grid) CreateNextLevel() {
	for _, t := range g.tiles {
		t.Reset()
	}

	rand.Seed(time.Now().UnixNano())

	palette := Palettes[rand.Int()%len(Palettes)]
	g.colorBackground = CalcBackgroundColor(palette)
	g.colorWall = ExtendedColors[palette[0]]

	g.createRooms()

	g.carve()

	ThePuck = NewPuck(g.findTile(TilesAcross/2, TilesDown/2))

	for i := 0; i < 4; i++ {
		TheGhosts = append(TheGhosts, NewGhost(g.randomTile()))
	}
}

// Layout implements ebiten.Game's Layout.
func (g *Grid) Layout(outsideWidth, outsideHeight int) (int, int) {
	for _, t := range g.tiles {
		t.Layout()
	}
	// g.puck.Layout(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

// Update the board state (transitions, user input)
func (g *Grid) Update() error {

	g.input.Update()

	switch {
	case inpututil.IsKeyJustReleased(ebiten.KeyBackspace):
		GSM.Switch(NewMenu())
	case inpututil.IsKeyJustReleased(ebiten.KeyUp):
		ThePuck.tile.toggleWall(0)
	case inpututil.IsKeyJustReleased(ebiten.KeyRight):
		ThePuck.tile.toggleWall(1)
	case inpututil.IsKeyJustReleased(ebiten.KeyDown):
		ThePuck.tile.toggleWall(2)
	case inpututil.IsKeyJustReleased(ebiten.KeyLeft):
		ThePuck.tile.toggleWall(3)
	}

	for _, t := range g.tiles {
		t.Update()
	}

	if g.input.TouchX != 0 && g.input.TouchY != 0 {
		pt := image.Point{g.input.TouchX - int(CameraX), g.input.TouchY - int(CameraY)}
		t := g.findTileAt(pt)
		if t != nil {
			// println("input on tile", t.X, t.Y, t.wallCount())
			if t.wallCount() < 4 {
				g.AllTiles(func(t *Tile) { t.parent = nil; t.marked = false })
				ThePuck.ThrowBallTo(t)
			}
		}
	}

	for i := 0; i < len(TheGhosts); i++ {
		TheGhosts[i].Update()
	}
	ThePuck.Update()

	return nil
}

// Draw renders the grid into the gridImage
func (g *Grid) Draw(screen *ebiten.Image) {

	screen.Fill(g.colorBackground)

	for _, t := range g.tiles {
		t.Draw(screen)
	}

	for i := 0; i < len(TheGhosts); i++ {
		TheGhosts[i].Draw(screen)
	}

	ThePuck.Draw(screen)

	if DebugMode {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("grid:%d,%d camera:%v,%v puck:%v,%v", TilesAcross, TilesDown, CameraX, CameraY, ThePuck.tile.X, ThePuck.tile.Y))
	}
}
