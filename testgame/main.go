package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 1000
	screenHeight = 380
)

var (
	running         = true
	backgroundColor = rl.NewColor(147, 211, 196, 255)

	grassSprite  rl.Texture2D
	playerSprite rl.Texture2D

	playerSrc  rl.Rectangle
	playerDest rl.Rectangle

	playerSpeed float32 = 3
)

func drawScene() {
	rl.DrawTexture(grassSprite, 100, 50, rl.White)
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerDest.Y -= playerSpeed
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerDest.Y += playerSpeed
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerDest.X -= playerSpeed
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerDest.X += playerSpeed
	}
}

func update() {
	running = !rl.WindowShouldClose()
}

func render() {
	// Begin drawing
	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)

	// Draw scene
	drawScene()

	// End drawing
	rl.EndDrawing()
}

func init() {
	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, "Benitolandia")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	// Load textures
	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")
	playerSprite = rl.LoadTexture("res/Characters/BasicCharakterSpritesheet.png")

	// Set player source rectangle
	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 100, 100)
}

func quit() {
	// Unload textures
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)

	// Close window
	rl.CloseWindow()
}

func main() {
	// Main loop
	for running {
		input()
		update()
		render()
	}

	quit()
}
