// Copyright ©️ 2021 oddstream.games

package maze

import (
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// StatePuckSettled when the Puck is not moving
	StatePuckSettled = iota
	// StatePuckMoving when Puck is lerping/chasing it's ball
	StatePuckMoving
)

// PuckState of this Puck
type PuckState int

// Puck defines the yellow blob/player avatar
type Puck struct {
	state PuckState

	tile *Tile // tile we are sitting on
	home *Tile // tile we started on

	puckImage *ebiten.Image

	x, y float64

	ball *Ball
}

// NewPuck creates a new Puck object
func NewPuck(start *Tile) *Puck {
	p := &Puck{home: start, tile: start, state: StatePuckSettled}

	dc := gg.NewContext(TileSize, TileSize)
	dc.SetRGB(1, 1, 0) // Yellow
	dc.DrawCircle(float64(TileSize/2), float64(TileSize/2), float64(TileSize/3))
	dc.Fill()
	dc.Stroke()
	p.puckImage = ebiten.NewImageFromImage(dc.Image())

	p.x, p.y = p.tile.Position()

	p.ball = NewBall(start)

	return p
}

// ThrowBallTo a target tile
func (p *Puck) ThrowBallTo(t *Tile) {
	p.ball.ThrowTo(t)
}

// BallTile getter for location of puck's ball
// func (p *Puck) BallTile() *Tile {
// 	return p.ball.Tile()
// }

/*
function Util.BFS(tStart, tDst)
  assert(tStart)
  assert(tDst)
  tStart.got:iterator(function(t) t.parent = nil end)
  local q = {tStart}        -- push onto queue
  tStart.parent = tStart   -- mark as itself
  while #q > 0 do
    local t = table.remove(q, 1)
    assert(t)
    if t == tDst then
      return
    end
    for dir = 1, 4 do
      local tn = t:neighbour(dir)
      if tn and (not t:isWall(dir)) and (tn.parent == nil) then
        tn.parent = t
        q[#q+1] = tn
        -- table.insert(q, tn) -- push to end of q
      end
    end
  end
  assert(false, 'BFS not found')
end
*/

// Update the state/position of the Puck
func (p *Puck) Update() error {
	p.ball.Update()
	if p.tile != p.ball.tile {

	}
	return nil
}

// Draw the Puck
func (p *Puck) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(p.puckImage, op)
	p.ball.Draw(screen)
}
