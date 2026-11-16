package main

import (
 "boids/simulation"
 "boids/ui"
 "math/rand/v2"
 rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
 rl.InitWindow(simulation.ScreenX+ui.SidebarWidth, simulation.ScreenY, "Boids - Concurrent Development")
 defer rl.CloseWindow()
 rl.SetTargetFPS(60)
 spawnBoids()
 target := rl.LoadRenderTexture(simulation.ScreenX, simulation.ScreenY)
 defer rl.UnloadRenderTexture(target)

 for !rl.WindowShouldClose() {
  if rl.IsKeyPressed(rl.KeyM) { simulation.ToggleMap() }
  if simulation.DetectionRange <= 0 { simulation.DoGridOptimize = false }
  dt := rl.GetFrameTime()
  for _, b := range simulation.Boids { simulation.BoidAction(b, dt) }

  // Render the simulation in its own 1100x700 coordinate system.
  rl.BeginTextureMode(target)
  rl.ClearBackground(rl.RayWhite)
  if simulation.DoGridOptimize { ui.DrawGrid() }
  simulation.DrawObstacles()
  for i, b := range simulation.Boids { simulation.DrawBoid(b, i == 0) }
  rl.EndTextureMode()

  // Draw the simulation to the right, then draw interactive controls
  // in screen coordinates on the left (not inside the render texture).
  rl.BeginDrawing()
  rl.ClearBackground(rl.RayWhite)
  rl.DrawTextureRec(target.Texture,
   rl.Rectangle{X:0,Y:0,Width:float32(target.Texture.Width),Height:-float32(target.Texture.Height)},
   rl.Vector2{X:float32(ui.SidebarWidth),Y:0}, rl.White)
  ui.InitGui()
  ui.DrawGUI(SetBoidsAmount)
  rl.DrawFPS(ui.SidebarWidth+simulation.ScreenX-120, 20)
  rl.EndDrawing()
 }
}

func SetBoidsAmount(newAmount int32) {
 simulation.BoidsAmount = newAmount
 spawnBoids()
}

func spawnBoids() {
 simulation.Boids = []*simulation.Boid{}
 for range simulation.BoidsAmount {
  b := &simulation.Boid{Position:rl.Vector2{
   X:rand.Float32()*150+simulation.Bounds.Width*0.5,
   Y:rand.Float32()*150+simulation.Bounds.Height*0.5,
  }}
  simulation.ResolveObstacleCollision(b,b.Position)
  simulation.Boids = append(simulation.Boids,b)
 }
}
