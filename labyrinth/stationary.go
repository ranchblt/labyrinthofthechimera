package labyrinth

import "github.com/hajimehoshi/ebiten"

type Stationary struct {
	Image   *ebiten.Image
	topLeft *coord
	scale   float64
}

func (s *Stationary) Len() int {
	return 1
}

func (s *Stationary) Dst(i int) (x0, y0, x1, y1 int) {
	return scaledDST(i, s.topLeft, s.Image, s.scale)
}

func (s *Stationary) Src(i int) (x0, y0, x1, y1 int) {
	w, h := s.Image.Size()
	return 0, 0, w, h
}

func (s *Stationary) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.Image, &ebiten.DrawImageOptions{
		ImageParts: s,
	})
}
