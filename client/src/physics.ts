// Physics types and constants

import { PositionComponent } from "./Components";

export class Vector2 {
  constructor(public x: number, public y: number) {}

  add(v: Vector2): Vector2 {
    return new Vector2(this.x + v.x, this.y + v.y);
  }

  subtract(v: Vector2): Vector2 {
    return new Vector2(this.x - v.x, this.y - v.y);
  }

  multiply(s: number): Vector2 {
    return new Vector2(this.x * s, this.y * s);
  }

  dot(v: Vector2): number {
    return this.x * v.x + this.y * v.y;
  }

  lengthSquared(): number {
    return this.dot(this);
  }

  length(): number {
    return Math.sqrt(this.lengthSquared());
  }

  norm(): Vector2 {
    const len = this.length();
    return len !== 0 ? this.multiply(1 / len) : new Vector2(0, 0);
  }
}

export class Circle {
  velocity = new Vector2(0, 0);

  constructor(
    public center: Vector2,
    public radius: number
  ) {}

  update(cenX: number, cenY: number, velX: number, velY: number) {
    this.center.x = cenX;
    this.center.y = cenY;
    this.velocity.x = velX;
    this.velocity.y = velY;
  }
}

export class Rect {
  constructor(
    public topLeft: Vector2,
    public width: number,
    public height: number
  ) {}
}

export class ShotData {
  constructor(
    public direction: Vector2,
    public power: number
  ) {}
}

// Constants
const dt = 1 / 120;
const drag = 0.989;
const stop2 = 1e-4;
const settleNeed = 5;
const maxSteps = 30 * 120;

export interface PhysicsState {
  activeCircle: Circle;
  circles: Circle[];
  shot: ShotData;
  walls: Rect[];
  step: number;
  slowFrames: number;
  onComplete: () => void;
}

// Physics functions
function applyImpulse(circle: Circle, shot: ShotData): void {
  circle.velocity = shot.direction.norm().multiply(shot.power);
}

function integrate(circles: Circle[]): void {
  for (const c of circles) {
    c.center = c.center.add(c.velocity.multiply(dt));
  }
}

function checkCircleCircle(c1: Circle, c2: Circle): boolean {
  return c1.center.subtract(c2.center).lengthSquared() < (c1.radius + c2.radius) ** 2;
}

function doPositionalCorrection(c1: Circle, c2: Circle): void {
  const delta = c1.center.subtract(c2.center);
  const dist = delta.length();
  const radSum = c1.radius + c2.radius;
  const penetration = radSum - dist;
  if (dist === 0) return;

  const normal = delta.multiply(1 / dist);
  const correction = normal.multiply(penetration * 0.5);

  c1.center = c1.center.add(correction);
  c2.center = c2.center.subtract(correction);
}

function resolveCircleCircle(circles: Circle[]): void {
  for (let i = 0; i < circles.length; i++) {
    for (let j = i + 1; j < circles.length; j++) {
      const c1 = circles[i];
      const c2 = circles[j];
      if (checkCircleCircle(c1, c2)) {
        doPositionalCorrection(c1, c2);

        const relVel = c1.velocity.subtract(c2.velocity);
        const delta = c1.center.subtract(c2.center);
        const dist = delta.length();
        if (dist === 0) continue;

        const normal = delta.multiply(1 / dist);
        const speed = relVel.dot(normal);
        if (speed >= 0) continue;

        const impulse = normal.multiply(-speed);
        c1.velocity = c1.velocity.add(impulse);
        c2.velocity = c2.velocity.subtract(impulse);
      }
    }
  }
}

function resolveCircleWallCollisions(circles: Circle[], walls: Rect[]): void {
  for (const circle of circles) {
    for (const wall of walls) {
      const left = wall.topLeft.x;
      const right = wall.topLeft.x + wall.width;
      const top = wall.topLeft.y;
      const bottom = wall.topLeft.y + wall.height;

      // Check left/right collision
      if (circle.center.x - circle.radius < left) {
        circle.center.x = left + circle.radius;
        circle.velocity.x *= -1;
      } else if (circle.center.x + circle.radius > right) {
        circle.center.x = right - circle.radius;
        circle.velocity.x *= -1;
      }

      // Check top/bottom collision
      if (circle.center.y - circle.radius < top) {
        circle.center.y = top + circle.radius;
        circle.velocity.y *= -1;
      } else if (circle.center.y + circle.radius > bottom) {
        circle.center.y = bottom - circle.radius;
        circle.velocity.y *= -1;
      }
    }
  }
}

function applyFriction(circles: Circle[]): void {
  for (const c of circles) {
    c.velocity = c.velocity.multiply(drag);
  }
}

function allStopped(circles: Circle[]): boolean {
  return circles.every(c => c.velocity.lengthSquared() < stop2);
}

function resetVelocities(circles: Circle[]): void {
  for (const c of circles) {
    c.velocity = new Vector2(0, 0);
  }
}

// Main physics resolver
export function startPhysicsSimulation(
  activeCircle: Circle,
  circles: Circle[],
  shot: ShotData,
  walls: Rect[],
  onComplete: () => void
): void {
  const state: PhysicsState = {
    activeCircle,
    circles,
    shot,
    walls,
    step: 0,
    slowFrames: 0,
    onComplete
  };

  applyImpulse(state.activeCircle, state.shot);
  requestAnimationFrame(() => physicsStep(state));
}

function physicsStep(state: PhysicsState): void {
  integrate(state.circles);
  resolveCircleWallCollisions(state.circles, state.walls);
  resolveCircleCircle(state.circles);
  applyFriction(state.circles);

  if (allStopped(state.circles)) {
    state.slowFrames++;
    if (state.slowFrames >= settleNeed) {
      resetVelocities(state.circles);
      state.onComplete();
      return;
    }
  } else {
    state.slowFrames = 0;
  }

  state.step++;
  if (state.step < maxSteps) {
    requestAnimationFrame(() => physicsStep(state));
  } else {
    state.onComplete(); // force exit
  }
}