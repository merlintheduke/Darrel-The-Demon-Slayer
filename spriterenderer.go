package main
import rl "github.com/gen2brain/raylib-go/raylib"

type SpriteRenderer struct {
	Sprite   rl.Texture2D
	Color    rl.Color
	Position rl.Vector2
	Angle    float32
	Scale    float32
	frameSize int
	totalFrames int
	currentFrame int
}

func NewSpriteRenderer(sprite int, newColor rl.Color, newPosition rl.Vector2,scale float32, angle float32) SpriteRenderer {
	sr := SpriteRenderer{
		Sprite:   textures[sprite],
		Color:    newColor,
		Position: newPosition,
		Angle:    angle,
		Scale:    scale,
		frameSize: 100,
		totalFrames: 1,
		currentFrame: 1,
	}
	return sr
}


func (sr SpriteRenderer) Draw() {
	frameHeight := float32(sr.Sprite.Height) / float32(sr.totalFrames)
	sourceRect := rl.NewRectangle(0, float32((sr.currentFrame-1))*frameHeight, float32(sr.Sprite.Width),frameHeight)
	destRect := rl.NewRectangle(sr.Position.X, sr.Position.Y, float32(sr.Sprite.Width)*sr.Scale, frameHeight*sr.Scale)
	origin := rl.Vector2Scale(
    rl.NewVector2(
        float32(sr.Sprite.Width)/2,
        frameHeight/2,
    ),
    sr.Scale,
)
	rl.DrawTexturePro(sr.Sprite, sourceRect,
		destRect,
		origin, sr.Angle, sr.Color)
}

func (sr *SpriteRenderer) nextFrame() {
	sr.currentFrame++
	if sr.currentFrame > sr.totalFrames {
		sr.currentFrame = 1
	}
}