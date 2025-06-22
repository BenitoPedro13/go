package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	screenWidth  = 1000
	screenHeight = 800
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

	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     map[int][2]int // Mapeamento dos índices para coordenadas [x,y] no tileset
	mapW, mapH int

	playerSpeed float32 = 3

	musicPaused bool
	music       rl.Music

	cam rl.Camera2D
)

// Função para obter o tipo de tile correto baseado nos vizinhos
func getAutotileIndex(x, y int) int {
	// Verificar se existe tile nas posições adjacentes
	hasTop := y > 0 && tileMap[(y-1)*mapW+x] != 0
	hasBottom := y < mapH-1 && tileMap[(y+1)*mapW+x] != 0
	hasLeft := x > 0 && tileMap[y*mapW+(x-1)] != 0
	hasRight := x < mapW-1 && tileMap[y*mapW+(x+1)] != 0

	// Sistema de autotiling baseado em bitmask
	// Baseado no tileset 3x3 padrão:
	// 1  2  3
	// 4  5  6
	// 7  8  9

	// Cantos individuais
	if !hasTop && !hasLeft && hasRight && hasBottom {
		return 1 // Canto superior esquerdo
	}
	if !hasTop && hasLeft && !hasRight && hasBottom {
		return 3 // Canto superior direito
	}
	if hasTop && !hasLeft && hasRight && !hasBottom {
		return 7 // Canto inferior esquerdo
	}
	if hasTop && hasLeft && !hasRight && !hasBottom {
		return 9 // Canto inferior direito
	}

	// Bordas
	if !hasTop && hasLeft && hasRight && hasBottom {
		return 2 // Borda superior
	}
	if hasTop && hasLeft && hasRight && !hasBottom {
		return 8 // Borda inferior
	}
	if hasTop && !hasLeft && hasRight && hasBottom {
		return 4 // Borda esquerda
	}
	if hasTop && hasLeft && !hasRight && hasBottom {
		return 6 // Borda direita
	}

	// Tile central (completamente cercado)
	if hasTop && hasBottom && hasLeft && hasRight {
		return 5
	}

	// Tile isolado
	if !hasTop && !hasBottom && !hasLeft && !hasRight {
		return 5 // Usar tile central para tile isolado
	}

	return 5 // Padrão - tile central
}

func drawScene() {
	// rl.DrawTexture(grassSprite, 100, 50, rl.White)

	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapW)
			tileDest.Y = tileDest.Height * float32(i/mapW)

			// Usar autotiling para obter o índice correto
			x := i % mapW
			y := i / mapW
			tileIndex := getAutotileIndex(x, y)

			// Usar srcMap para obter as coordenadas corretas no tileset
			if coords, exists := srcMap[tileIndex]; exists {
				tileSrc.X = tileSrc.Width * float32(coords[0])
				tileSrc.Y = tileSrc.Height * float32(coords[1])
			} else {
				// Fallback para tile central se não encontrar
				tileSrc.X = tileSrc.Width * 1  // Centro X
				tileSrc.Y = tileSrc.Height * 1 // Centro Y
			}

			rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
		}
	}

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

func loadMap() {
	mapW = 50
	mapH = 60

	// Inicializar com tiles vazios
	tileMap = make([]int, mapW*mapH)

	// Criar algumas áreas de grama com formatos interessantes
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			i := y*mapW + x

			// Área principal de grama (grande retângulo)
			if x >= 5 && x < 45 && y >= 5 && y < 55 {
				tileMap[i] = 1
			}

			// Remover algumas áreas para criar lagos/espaços vazios
			// Lago circular no centro
			centerX, centerY := 25, 30
			dx, dy := float32(x-centerX), float32(y-centerY)
			if dx*dx+dy*dy < 36 { // Raio de 6
				tileMap[i] = 0
			}

			// Pequeno lago no canto superior direito
			if x >= 35 && x < 42 && y >= 8 && y < 15 {
				tileMap[i] = 0
			}

			// Caminho/rio diagonal
			if x-y > 15 && x-y < 20 && y > 20 && y < 50 {
				tileMap[i] = 0
			}

			// Pequenas ilhas
			if x >= 15 && x < 20 && y >= 45 && y < 50 {
				tileMap[i] = 1
			}
			if x >= 30 && x < 35 && y >= 15 && y < 20 {
				tileMap[i] = 1
			}
		}
	}
}

func init() {
	// Initialize window
	rl.InitWindow(screenWidth, screenHeight, "Benitolandia")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	// Load textures
	grassSprite = rl.LoadTexture("res/Tilesets/Grass.png")

	// Set tile source and destination rectangles (aumentar o tamanho para melhor visualização)
	tileDest = rl.NewRectangle(0, 0, 32, 32)
	tileSrc = rl.NewRectangle(0, 0, 16, 16)

	// Inicializar srcMap com coordenadas corretas do tileset 3x3
	// Baseado no layout: 1=canto sup.esq, 2=borda sup, 3=canto sup.dir
	//                    4=borda esq,    5=centro,    6=borda dir
	//                    7=canto inf.esq, 8=borda inf, 9=canto inf.dir
	srcMap = map[int][2]int{
		1: {0, 0}, // Canto superior esquerdo
		2: {1, 0}, // Borda superior
		3: {2, 0}, // Canto superior direito
		4: {0, 1}, // Borda esquerda
		5: {1, 1}, // Centro
		6: {2, 1}, // Borda direita
		7: {0, 2}, // Canto inferior esquerdo
		8: {1, 2}, // Borda inferior
		9: {2, 2}, // Canto inferior direito
	}

	playerSprite = rl.LoadTexture("res/Characters/BasicCharakterSpritesheet.png")

	// Set player source rectangle
	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	// Posicionar o player no centro da área de grama
	playerDest = rl.NewRectangle(25*32, 30*32, 100, 100)

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
		1.3, // Reduzir o zoom para ver mais do mapa
	)

	// Load map
	loadMap()
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
