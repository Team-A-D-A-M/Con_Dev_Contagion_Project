package simulation

import rl "github.com/gen2brain/raylib-go/raylib"

type GridCell struct { X,Y int }
type Boid struct { Position rl.Vector2; Velocity rl.Vector2; GridCell GridCell }
type NearbyBoid struct { Boid *Boid; Distance float32 }

func boidCohesion(b *Boid,nearby []*NearbyBoid) rl.Vector2 {
 if len(nearby)==0{return rl.Vector2Zero()}
 center:=rl.Vector2Zero();for _,n:=range nearby {center=rl.Vector2Add(center,n.Boid.Position)}
 center=rl.Vector2Scale(center,1/float32(len(nearby)))
 desired:=rl.Vector2Subtract(center,b.Position)
 if desired.Length()>0 {desired=rl.Vector2Scale(rl.Vector2Normalize(desired),MaxSpeed)}
 return rl.Vector2Subtract(desired,b.Velocity)
}
func boidSeparation(b *Boid,nearby []*NearbyBoid) rl.Vector2 {
 steer:=rl.Vector2Zero();count:=0
 for _,n:=range nearby {d:=n.Distance;if d<=0.001||d>=DesiredSeparation {continue};v:=rl.Vector2Subtract(b.Position,n.Boid.Position);steer=rl.Vector2Add(steer,rl.Vector2Scale(rl.Vector2Normalize(v),1/d));count++}
 if count==0{return rl.Vector2Zero()};steer=rl.Vector2Scale(steer,1/float32(count));if steer.Length()==0{return rl.Vector2Zero()}
 return rl.Vector2Subtract(rl.Vector2Scale(rl.Vector2Normalize(steer),MaxSpeed),b.Velocity)
}
func boidAlignment(b *Boid,nearby []*NearbyBoid) rl.Vector2 {
 if len(nearby)==0{return rl.Vector2Zero()};avg:=rl.Vector2Zero();for _,n:=range nearby {avg=rl.Vector2Add(avg,n.Boid.Velocity)}
 avg=rl.Vector2Scale(avg,1/float32(len(nearby)));if avg.Length()==0{return rl.Vector2Zero()}
 return rl.Vector2Subtract(rl.Vector2Scale(rl.Vector2Normalize(avg),MaxSpeed),b.Velocity)
}
func BoidAction(b *Boid,dt float32) {
 nearby:=getNearbyBoids(b)
 cohesion:=boidCohesion(b,nearby);separation:=boidSeparation(b,nearby);alignment:=boidAlignment(b,nearby)
 flocking:=rl.Vector2Add(rl.Vector2Scale(cohesion,CohesionValue),rl.Vector2Add(rl.Vector2Scale(separation,SeparationValue),rl.Vector2Scale(alignment,AlignmentValue)))
 avoidance:=obstacleAvoidance(b)
 steering:=limitVector(flocking,MaxForce)
 if avoidance.Length()>0 {steering=limitVector(avoidance,MaxForce*7)}
 previous:=b.Position
 b.Velocity=limitVector(rl.Vector2Add(b.Velocity,rl.Vector2Scale(steering,dt)),MaxSpeed)
 b.Position=rl.Vector2Add(b.Position,rl.Vector2Scale(b.Velocity,dt))
 ResolveObstacleCollision(b,previous)
 b.Position=positionWrap(b.Position)
 if DoGridOptimize {updateBoidGrid(b)}
}
func updateBoidGrid(b *Boid) {b.GridCell.X=floorDiv(b.Position.X,DetectionRange);b.GridCell.Y=floorDiv(b.Position.Y,DetectionRange)}
func getNearbyBoids(b *Boid) []*NearbyBoid {
 nearby:=[]*NearbyBoid{}
 for _,other:=range Boids {if other==b {continue};dx:=b.GridCell.X-other.GridCell.X;dy:=b.GridCell.Y-other.GridCell.Y
 if DoGridOptimize&&(dx < -1||dx>1||dy < -1||dy>1){continue}
 d:=rl.Vector2Distance(b.Position,other.Position);if d<DetectionRange {nearby=append(nearby,&NearbyBoid{Boid:other,Distance:d})}}
 return nearby
}
var lineLength float32=24
func DrawBoid(b *Boid,highlight bool) {
 col:=rl.Blue;if highlight {col=rl.Red};vel:=rl.Vector2Normalize(b.Velocity)
 if highlight {rl.DrawCircleLines(int32(b.Position.X),int32(b.Position.Y),DetectionRange,rl.LightGray)}
 rl.DrawLine(int32(b.Position.X),int32(b.Position.Y),int32(b.Position.X+vel.X*lineLength),int32(b.Position.Y+vel.Y*lineLength),rl.LightGray)
 rl.DrawPoly(b.Position,3,6,rl.Vector2Angle(right2d,b.Velocity)*rl.Rad2deg,col)
}
func positionWrap(p rl.Vector2) rl.Vector2 {
 if p.X<0 {p.X+=Bounds.Width}else if p.X>Bounds.Width {p.X-=Bounds.Width}
 if p.Y<0 {p.Y+=Bounds.Height}else if p.Y>Bounds.Height {p.Y-=Bounds.Height}
 return p
}
