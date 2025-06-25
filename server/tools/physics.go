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
	center   Vector2
	radius   float64
	velocity Vector2
}

func (c *Circle) DistanceSquared(c1 *Circle) float64 {
	return c.center.Subtract(c1.center).LengthSquared()
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

func PhysicsResolver(activePlayer *Circle, playerPositions []*Circle, walls []Rect, shotData ShotData) {
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
	circle.velocity = shotData.direction.Norm().Multiply(float64(shotData.power))
}

func Integrate(circles []*Circle) {
	for i := range circles {
		circles[i].center = circles[i].center.Add(circles[i].velocity.Multiply(dt))
	}
}

func CheckCircleCircleCollision(c1 Circle, c2 Circle) bool {
	return (c1.center.Subtract(c2.center).LengthSquared() < (c1.radius+c2.radius)*(c1.radius+c2.radius))
}

func CheckCircleWallCollision(c1 Circle, wall Rect) WallCollisionType {
	threshold := -0.05

	bottomCheck := (c1.center.Y + c1.radius) - (wall.topleft.Y - wall.height)
	topCheck := (wall.topleft.Y) - (c1.center.Y - c1.radius)
	leftCheck := (c1.center.X + c1.radius) - (wall.topleft.X)
	rightCheck := (wall.topleft.X + wall.width) - (c1.center.X - c1.radius)

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
	delta := c1.center.Subtract(c2.center)
	dist := delta.Length()
	radSum := c1.radius + c2.radius

	penetration := radSum - dist
	normal := delta.Multiply(1 / dist)

	c1.center = c1.center.Add(normal.Multiply(penetration * 0.5))
	c2.center = c2.center.Add(normal.Multiply(penetration * 0.5))
}

func FixOverlapX(dir WallCollisionType, circle *Circle, wall Rect) {
	if dir == LEFT {
		overlap := (circle.center.X + circle.radius) - wall.topleft.X
		if overlap > 0 {
			circle.center.X -= (overlap + 0.1)
		}
	} else if dir == RIGHT {
		overlap := (circle.center.X - circle.radius) - (wall.topleft.X + wall.width)
		if overlap < 0 {
			circle.center.X -= (overlap - 0.1)
		}
	}
}

func FixOverlapY(dir WallCollisionType, circle *Circle, wall Rect) {
	if dir == TOP {
		overlap := (circle.center.Y + circle.radius) - wall.topleft.Y
		if overlap > 0 {
			circle.center.Y -= (overlap + 0.1)
		}
	} else if dir == BOTTOM {
		overlap := (circle.center.Y - circle.radius) - (wall.topleft.Y + wall.height)
		if overlap < 0 {
			circle.center.Y -= (overlap - 0.1)
		}
	}
}

func ResolveCircleCircleCollisions(circles []*Circle) {
	//for all pairs, check collision and resolve if colliding
	for i := 0; i < len(circles); i++ {
		for j := i + 1; j < len(circles); j++ {
			if CheckCircleCircleCollision(*circles[i], *circles[j]) {
				DoPositionalCorrection(circles[i], circles[j])
				relVel := circles[i].velocity.Subtract(circles[j].velocity)
				delta := circles[i].center.Subtract(circles[j].center)
				dist := delta.Length()
				normal := delta.Multiply(1 / dist)
				speed := relVel.Dot(normal)
				if speed >= 0 {
					continue
				}

				impulse := normal.Multiply(-speed)
				circles[i].velocity = circles[i].velocity.Add(impulse)
				circles[j].velocity = circles[j].velocity.Subtract(impulse)
			}
		}
	}
}

func ResolveCircleWallCollisions(circles []*Circle, walls []Rect) {
	for i := range circles {
		for j := range walls {
			collisionDirection := CheckCircleWallCollision(*circles[i], walls[j])
			collisionDetected := false
			switch collisionDirection {
			case NONE:

			case TOP, BOTTOM:
				circles[i].velocity.Y *= -1
				FixOverlapY(collisionDirection, circles[i], walls[j])
				collisionDetected = true
			case LEFT, RIGHT:
				circles[i].velocity.X *= -1
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
		circles[i].velocity = circles[i].velocity.Multiply(drag)
	}
}

func AllStopped(circles []*Circle) bool {
	for i := range circles {
		if circles[i].velocity.LengthSquared() >= stop2 {
			return false
		}
	}

	return true
}

func ResetVelocities(circles []*Circle) {
	for i := range circles {
		circles[i].velocity = Vector2{0, 0}
	}
}
