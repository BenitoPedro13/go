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

	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerLeft, playerRight bool
	playerFrame                                   int

	frameCount int

	playerSpeed float32 = 3

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

func drawScene() {
	rl.DrawTexture(grassSprite, 100, 50, rl.White)
	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White)
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 1
		playerUp = true
	}
	if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 0
		playerDown = true
	}
	if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 2
		playerLeft = true
	}
	if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 3
		playerRight = true
	}
	if rl.IsKeyPressed(rl.KeyQ) {
		musicPaused = !musicPaused
		if musicPaused {
			rl.PauseMusicStream(music)
		} else {
			rl.ResumeMusicStream(music)
		}
	}
}

func update() {
	running = !rl.WindowShouldClose()

	// Reset movement flags at the beginning of the frame
	playerMoving = false
	playerUp = false
	playerDown = false
	playerLeft = false
	playerRight = false

	// Handle input after resetting flags
	input()

	playerSrc.X = playerSrc.Width * float32(playerFrame)

	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if frameCount%8 == 1 {
			playerFrame++
		}
	} else if frameCount%45 == 1 {
		// Reset to idle frame when not moving
		playerFrame++
	}

	frameCount++
	if playerFrame > 3 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 1 {
		playerFrame = 0
	}

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPaused {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-playerDest.Width/2), float32(playerDest.Y-playerDest.Height/2))
}

func render() {
	// Begin drawing
	rl.BeginDrawing()
	rl.ClearBackground(backgroundColor)
	rl.BeginMode2D(cam)

	// Draw scene
	drawScene()

	// End drawing
	rl.EndMode2D()
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

	// Load music
	rl.InitAudioDevice()
	music = rl.LoadMusicStream("res/sounds/game-music.mp3")
	musicPaused = false
	rl.PlayMusicStream(music)

	// Set camera
	cam = rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(playerDest.X-playerDest.Width/2, playerDest.Y-playerDest.Height/2),
		0.0,
		1.0,
	)
}

func quit() {
	// Unload textures
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)

	// Unload music
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()

	// Close window
	rl.CloseWindow()
}

func main() {
	// Main loop
	for running {
		update()
		render()
	}

	quit()
}
