package tools

import (
	"fmt"
	"math"
)

type Vector2 struct {
	X float64
	Y float64
}

func (v Vector2) Add(u Vector2) Vector2      { return Vector2{v.X + u.X, v.Y + u.Y} }
func (v Vector2) Subtract(u Vector2) Vector2 { return Vector2{v.X - u.X, v.Y - u.Y} }
func (v Vector2) Multiply(s float64) Vector2 { return Vector2{v.X * s, v.Y * s} }
func (v Vector2) Dot(u Vector2) float64      { return v.X*u.X + v.Y*u.Y }
func (v Vector2) LengthSquared() float64     { return v.Dot(v) }
func (v Vector2) Length() float64            { return math.Sqrt(v.LengthSquared()) }
func (v Vector2) Norm() Vector2 {
	if l := v.Length(); l != 0 {
		return v.Multiply(1 / l)
	}
	return Vector2{}
}

type Vector2Int struct {
	X int
	Y int
}

func (v Vector2Int) Add(u Vector2Int) Vector2Int      { return Vector2Int{v.X + u.X, v.Y + u.Y} }
func (v Vector2Int) Subtract(u Vector2Int) Vector2Int { return Vector2Int{v.X - u.X, v.Y - u.Y} }
func (v Vector2Int) DistanceSquared(u Vector2Int) int {
	return (v.X-u.X)*(v.X-u.X) + (v.Y-u.Y)*(v.Y-u.Y)
}

type WallCollisionType int

const (
	NONE WallCollisionType = iota
	BOTTOM
	TOP
	LEFT
	RIGHT
)

type Rect struct {
	topleft Vector2
	width   float64
	height  float64
}

type Circle struct {
	Center   Vector2
	Radius   float64
	Velocity Vector2
}

func (c *Circle) DistanceSquared(c1 *Circle) float64 {
	return c.Center.Subtract(c1.Center).LengthSquared()
}

type ShotData struct {
	direction Vector2
	power     int //should be a value from one to 5, could make this an enum.
}

const (
	dt         = 1.0 / 120.0
	drag       = 0.989
	stop2      = 1e-4
	settleNeed = 5
	maxSteps   = 30 * 120
)

func PhysicsResolver(activePlayer *Circle, playerPositions []*Circle, walls []*Rect, shotData ShotData) {
	ApplyImpulse(activePlayer, shotData)
	slowFrames := 0
	for step := 0; step < maxSteps; step++ {
		Integrate(playerPositions)
		ResolveCircleWallCollisions(playerPositions, walls)
		ResolveCircleCircleCollisions(playerPositions)
		ApplyFriction(playerPositions)

		if AllStopped(playerPositions) {
			slowFrames++
			if slowFrames >= settleNeed {
				ResetVelocities(playerPositions)
				break
			}
		} else {
			slowFrames = 0
		}
	}
}

// HELPER FUNCTIONS

func ApplyImpulse(circle *Circle, shotData ShotData) {
	circle.Velocity = shotData.direction.Norm().Multiply(float64(shotData.power))
}

func Integrate(circles []*Circle) {
	for i := range circles {
		circles[i].Center = circles[i].Center.Add(circles[i].Velocity.Multiply(dt))
	}
}

func CheckCircleCircleCollision(c1 Circle, c2 Circle) bool {
	return (c1.Center.Subtract(c2.Center).LengthSquared() < (c1.Radius+c2.Radius)*(c1.Radius+c2.Radius))
}

func CheckCircleWallCollision(c1 Circle, wall *Rect) WallCollisionType {
	threshold := -0.05

	bottomCheck := (c1.Center.Y + c1.Radius) - (wall.topleft.Y - wall.height)
	topCheck := (wall.topleft.Y) - (c1.Center.Y - c1.Radius)
	leftCheck := (c1.Center.X + c1.Radius) - (wall.topleft.X)
	rightCheck := (wall.topleft.X + wall.width) - (c1.Center.X - c1.Radius)

	if bottomCheck < threshold || topCheck < threshold || leftCheck < threshold || rightCheck < threshold {
		return NONE
	} else if bottomCheck > threshold && bottomCheck < -threshold {
		return BOTTOM
	} else if topCheck > threshold && topCheck < -threshold {
		return TOP
	} else if leftCheck > threshold && leftCheck < -threshold {
		return LEFT
	} else if rightCheck > threshold && rightCheck < -threshold {
		return RIGHT
	} else {
		fmt.Println("ERROR: Circle-Wall Collision failed -- CheckCircleWallCollision()")
		return NONE
	}
}

func DoPositionalCorrection(c1 *Circle, c2 *Circle) {
	delta := c1.Center.Subtract(c2.Center)
	dist := delta.Length()
	radSum := c1.Radius + c2.Radius

	penetration := radSum - dist
	normal := delta.Multiply(1 / dist)

	c1.Center = c1.Center.Add(normal.Multiply(penetration * 0.5))
	c2.Center = c2.Center.Add(normal.Multiply(penetration * 0.5))
}

func FixOverlapX(dir WallCollisionType, circle *Circle, wall *Rect) {
	if dir == LEFT {
		overlap := (circle.Center.X + circle.Radius) - wall.topleft.X
		if overlap > 0 {
			circle.Center.X -= (overlap + 0.1)
		}
	} else if dir == RIGHT {
		overlap := (circle.Center.X - circle.Radius) - (wall.topleft.X + wall.width)
		if overlap < 0 {
			circle.Center.X -= (overlap - 0.1)
		}
	}
}

func FixOverlapY(dir WallCollisionType, circle *Circle, wall *Rect) {
	if dir == TOP {
		overlap := (circle.Center.Y + circle.Radius) - wall.topleft.Y
		if overlap > 0 {
			circle.Center.Y -= (overlap + 0.1)
		}
	} else if dir == BOTTOM {
		overlap := (circle.Center.Y - circle.Radius) - (wall.topleft.Y + wall.height)
		if overlap < 0 {
			circle.Center.Y -= (overlap - 0.1)
		}
	}
}

func ResolveCircleCircleCollisions(circles []*Circle) {
	//for all pairs, check collision and resolve if colliding
	for i := 0; i < len(circles); i++ {
		for j := i + 1; j < len(circles); j++ {
			if CheckCircleCircleCollision(*circles[i], *circles[j]) {
				DoPositionalCorrection(circles[i], circles[j])
				relVel := circles[i].Velocity.Subtract(circles[j].Velocity)
				delta := circles[i].Center.Subtract(circles[j].Center)
				dist := delta.Length()
				normal := delta.Multiply(1 / dist)
				speed := relVel.Dot(normal)
				if speed >= 0 {
					continue
				}

				impulse := normal.Multiply(-speed)
				circles[i].Velocity = circles[i].Velocity.Add(impulse)
				circles[j].Velocity = circles[j].Velocity.Subtract(impulse)
			}
		}
	}
}

func ResolveCircleWallCollisions(circles []*Circle, walls []*Rect) {
	for i := range circles {
		for j := range walls {
			collisionDirection := CheckCircleWallCollision(*circles[i], walls[j])
			collisionDetected := false
			switch collisionDirection {
			case NONE:

			case TOP, BOTTOM:
				circles[i].Velocity.Y *= -1
				FixOverlapY(collisionDirection, circles[i], walls[j])
				collisionDetected = true
			case LEFT, RIGHT:
				circles[i].Velocity.X *= -1
				FixOverlapX(collisionDirection, circles[i], walls[j])
				collisionDetected = true
			}
			if collisionDetected {
				break
			}
		}
	}
}

func ApplyFriction(circles []*Circle) {
	for i := range circles {
		circles[i].Velocity = circles[i].Velocity.Multiply(drag)
	}
}

func AllStopped(circles []*Circle) bool {
	for i := range circles {
		if circles[i].Velocity.LengthSquared() >= stop2 {
			return false
		}
	}

	return true
}

func ResetVelocities(circles []*Circle) {
	for i := range circles {
		circles[i].Velocity = Vector2{0, 0}
	}
}
