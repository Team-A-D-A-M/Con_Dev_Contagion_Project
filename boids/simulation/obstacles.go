package simulation

import (
 "math"
 rl "github.com/gen2brain/raylib-go/raylib"
)

type Obstacle struct { Position rl.Vector2; Radius float32 }
type RectObstacle struct { Bounds rl.Rectangle }
var CurrentMap = 1
func ToggleMap() { if CurrentMap == 1 { CurrentMap = 2 } else { CurrentMap = 1 }; for _, b := range Boids { ResolveObstacleCollision(b, b.Position) } }
func GetMapName() string { if CurrentMap == 1 { return "MAP 1 - OBSTACLE FIELD" }; return "MAP 2 - MAZE" }
var Obstacles = []Obstacle{{rl.Vector2{X:230,Y:170},42},{rl.Vector2{X:430,Y:190},32},{rl.Vector2{X:830,Y:155},48},{rl.Vector2{X:930,Y:390},42},{rl.Vector2{X:300,Y:570},38},{rl.Vector2{X:710,Y:590},34}}
var RectObstacles = []RectObstacle{{rl.Rectangle{X:550,Y:100,Width:150,Height:35}},{rl.Rectangle{X:100,Y:360,Width:180,Height:35}},{rl.Rectangle{X:650,Y:420,Width:140,Height:40}}}
var MazeObstacles = []RectObstacle{{rl.Rectangle{X:150,Y:120,Width:330,Height:25}},{rl.Rectangle{X:560,Y:120,Width:330,Height:25}},{rl.Rectangle{X:150,Y:120,Width:25,Height:190}},{rl.Rectangle{X:465,Y:120,Width:25,Height:165}},{rl.Rectangle{X:875,Y:120,Width:25,Height:190}},{rl.Rectangle{X:275,Y:285,Width:300,Height:25}},{rl.Rectangle{X:635,Y:285,Width:240,Height:25}},{rl.Rectangle{X:275,Y:285,Width:25,Height:190}},{rl.Rectangle{X:600,Y:285,Width:25,Height:165}},{rl.Rectangle{X:850,Y:285,Width:25,Height:190}},{rl.Rectangle{X:275,Y:450,Width:325,Height:25}},{rl.Rectangle{X:675,Y:450,Width:200,Height:25}},{rl.Rectangle{X:450,Y:450,Width:25,Height:125}},{rl.Rectangle{X:725,Y:450,Width:25,Height:125}},{rl.Rectangle{X:150,Y:575,Width:300,Height:25}},{rl.Rectangle{X:500,Y:575,Width:350,Height:25}}}
func activeRects() []RectObstacle { if CurrentMap == 1 { return RectObstacles }; return MazeObstacles }
func DrawObstacles() { if CurrentMap == 1 { for _, o := range Obstacles { rl.DrawCircleV(o.Position,o.Radius,rl.Gray); rl.DrawCircleLines(int32(o.Position.X),int32(o.Position.Y),o.Radius,rl.DarkGray) } }; for _, o := range activeRects() { c:=rl.Gray; if CurrentMap==2 { c=rl.DarkGray }; rl.DrawRectangleRec(o.Bounds,c); r:=o.Bounds; rl.DrawRectangleLines(int32(r.X),int32(r.Y),int32(r.Width),int32(r.Height),rl.Black) } }
func clamp(x,lo,hi float32) float32 { if x<lo {return lo}; if x>hi{return hi}; return x }
func dot(a,b rl.Vector2) float32 {return a.X*b.X+a.Y*b.Y}
func sub(a,b rl.Vector2) rl.Vector2{return rl.Vector2Subtract(a,b)}
func add(a,b rl.Vector2) rl.Vector2{return rl.Vector2Add(a,b)}
func mul(a rl.Vector2,k float32) rl.Vector2{return rl.Vector2Scale(a,k)}
func length(a rl.Vector2) float32{return a.Length()}
func unit(a rl.Vector2) rl.Vector2 {if length(a)<0.0001{return rl.Vector2{X:1}};return rl.Vector2Normalize(a)}
func closestRect(p rl.Vector2,r rl.Rectangle) rl.Vector2 {return rl.Vector2{X:clamp(p.X,r.X,r.X+r.Width),Y:clamp(p.Y,r.Y,r.Y+r.Height)}}
// pointNormal returns an outward normal and the signed distance to the obstacle surface.
// Negative distances mean the point is inside the solid shape.
func circleSurface(p rl.Vector2,o Obstacle) (rl.Vector2,float32) {v:=sub(p,o.Position);d:=length(v);return unit(v),d-o.Radius}
func rectSurface(p rl.Vector2,r rl.Rectangle) (rl.Vector2,float32) {
 q:=closestRect(p,r);v:=sub(p,q);d:=length(v)
 if d>0.0001{return mul(v,1/d),d}
 left:=p.X-r.X;right:=r.X+r.Width-p.X;top:=p.Y-r.Y;bottom:=r.Y+r.Height-p.Y
 m:=left;n:=rl.Vector2{X:-1};if right<m {m=right;n=rl.Vector2{X:1}};if top<m {m=top;n=rl.Vector2{Y:-1}};if bottom<m {m=bottom;n=rl.Vector2{Y:1}};return n,-m
}
func nearestSurface(p rl.Vector2) (rl.Vector2,float32) {
 best:=float32(math.MaxFloat32);normal:=rl.Vector2{X:1}
 if CurrentMap==1 {for _,o:=range Obstacles {n,d:=circleSurface(p,o);if d<best {best=d;normal=n}}}
 for _,o:=range activeRects() {n,d:=rectSurface(p,o.Bounds);if d<best {best=d;normal=n}}
 return normal,best
}
// A short look-ahead keeps avoidance local, while solid collision handling prevents tunnelling.
func obstacleAvoidance(b *Boid) rl.Vector2 {
 speed:=length(b.Velocity);if speed<0.001 {return rl.Vector2Zero()}
 dir:=unit(b.Velocity);look:=clamp(speed*0.18,35,115)
 best:=float32(math.MaxFloat32);normal:=rl.Vector2Zero();at:=float32(0)
 for i:=0;i<=12;i++ {t:=float32(i)/12;point:=add(b.Position,mul(dir,look*t));n,d:=nearestSurface(point);if d<best {best=d;normal=n;at=t}}
 if best>22 {return rl.Vector2Zero()}
 // Choose a consistent tangent to slide past the surface, instead of steering backwards.
 tangent:=rl.Vector2{X:-normal.Y,Y:normal.X};if dot(tangent,dir)<0 {tangent=mul(tangent,-1)}
 urgency:=clamp((22-best)/42,0,1)
 desired:=add(mul(tangent,0.75),mul(normal,0.25+0.6*urgency))
 if at==0 && best<8 {desired=add(desired,mul(normal,1.2))}
 return mul(sub(mul(unit(desired),MaxSpeed),b.Velocity),urgency)
}
// ResolveObstacleCollision moves a boid outside solid geometry. previous is its prior position.
// Movement is swept in small steps, so high speed cannot jump across thin walls.
func ResolveObstacleCollision(b *Boid,previous rl.Vector2) {
 const margin float32=7
 delta:=sub(b.Position,previous);distance:=length(delta)
 steps:=int(math.Ceil(float64(distance/4)));if steps<1 {steps=1};if steps>256 {steps=256}
 pos:=previous
 // Handle map changes or spawn positions already inside a shape.
 for j:=0;j<8;j++ {n,d:=nearestSurface(pos);if d>=margin {break};pos=add(pos,mul(n,margin-d+0.1));vn:=dot(b.Velocity,n);if vn<0 {b.Velocity=sub(b.Velocity,mul(n,vn))}}
 step:=mul(delta,1/float32(steps))
 for i:=0;i<steps;i++ {
  candidate:=add(pos,step)
  for j:=0;j<4;j++ {n,d:=nearestSurface(candidate);if d>=margin {break};candidate=add(candidate,mul(n,margin-d+0.1));vn:=dot(b.Velocity,n);if vn<0 {b.Velocity=sub(b.Velocity,mul(n,vn))};into:=dot(step,n);if into<0 {step=sub(step,mul(n,into))}}
  pos=candidate
 }
 b.Position=pos
}
