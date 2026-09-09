package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuAction int
type Menu struct {
	pos        rl.Vector2
	background SpriteRenderer
	textbox    []*TextBox
	buttons    []*Button
	displayed  bool
	CustomDraw func()
}
type Button struct {
	OnClick        func()
	sprite         SpriteRenderer
	usesSprite     bool
	color          rl.Color
	rect           rl.Rectangle
	text           string
	textcolor      rl.Color
	showtext       bool
	hover          bool
	Showbutton     bool
	allowUserInput bool
}
type TextBox struct {
	pos         rl.Vector2
	contents    string
	size        float32
	showTextBox bool
}

func (menu *Menu) newTextBox(pos rl.Vector2, contents string, size float32) *TextBox {
	return &TextBox{
		pos:      rl.Vector2Add(menu.pos, pos),
		contents: contents,
		size:     size,
	}
}

func (t TextBox) displayText() {
	rl.DrawText(t.contents, int32(t.pos.X), int32(t.pos.Y), int32(t.size), rl.Black)
}
func (menu *Menu) newButton(OnClick func(), pos rl.Vector2, width, height float32, text string, spritenumber int, usesSprite bool) *Button {
	sprite := NewSpriteRenderer(spritenumber, rl.White, rl.Vector2Add(menu.pos, pos), 1, 0)
	rec := rl.Rectangle{
		X:      pos.X + menu.pos.X,
		Y:      pos.Y + menu.pos.Y,
		Width:  width,
		Height: height,
	}
	return &Button{
		OnClick:        OnClick,
		sprite:         sprite,
		usesSprite:     usesSprite,
		rect:           rec,
		color:          rl.Beige,
		text:           text,
		showtext:       true,
		hover:          false,
		textcolor:      rl.Black,
		Showbutton:     true,
		allowUserInput: false,
	}

}

/*
func (b *button) displayButton() {

		if b.hover {
			rl.DrawTexture(b.sprite, int32(b.pos.X), int32(b.pos.Y), rl.White)
		} else {
			rl.DrawTexture(b.sprite, int32(b.pos.X), int32(b.pos.Y), rl.Gray)
		}
	}
*/
func (b *Button) displayButton() { //implement button sprite
	if b.Showbutton {
		if b.usesSprite {
			b.sprite.Draw()
		} else {
			if !b.hover {
				rl.DrawRectangleRec(b.rect, b.color)
			} else {
				rl.DrawRectangleRec(b.rect, rl.Yellow)
			}
		}
		rl.DrawText(b.text, b.rect.ToInt32().X, b.rect.ToInt32().Y, 50, b.textcolor)
	}
}

func newMenu(pos rl.Vector2) Menu {
	return Menu{
		pos:       pos,
		displayed: true,
	}
}
func (menu *Menu) Draw() {
	menu.background.Draw()
	if menu.CustomDraw != nil {
		menu.CustomDraw()
	}

	for i := range menu.buttons {
		menu.buttons[i].displayButton()
	}

	for _, textbox := range menu.textbox {
		textbox.displayText()
	}
}
func (m *Menu) Update(gc *GameController) {
	mouse := rl.GetMousePosition()

	for i := range m.buttons {
		b := m.buttons[i]

		b.hover = rl.CheckCollisionPointRec(mouse, b.rect)

		if b.hover && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			b.OnClick()
		}
		if b.allowUserInput {
			b.getUserInput(gc)
		}

	}
}

func (b *Button) getUserInput(gc *GameController) {
	b.allowUserInput = true
	key := rl.GetCharPressed()
	if key > 0 {
		if key >= 32 && key <= 125 {
			b.text += string(rune(key))
		}
		key = rl.GetCharPressed()
	}
	if rl.IsKeyPressed(rl.KeyBackspace) && len(b.text) > 0 {
		b.text = b.text[:len(b.text)-1]
	}
	if rl.IsKeyPressed(rl.KeyEnter) {
		b.allowUserInput = false

		gc.NewSaveGame(b.text)
		
	}
}
