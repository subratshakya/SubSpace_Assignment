package stealth

import (
	"math"
	"math/rand"
	"time"
)

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// GenerateBezierPath creates a human-like mouse path using Cubic Bézier curves
// with overshoot and correction to mimic natural mouse movement
func GenerateBezierPath(start, end Point) []Point {
	rand.Seed(time.Now().UnixNano())
	
	// Calculate distance
	dx := end.X - start.X
	dy := end.Y - start.Y
	distance := math.Sqrt(dx*dx + dy*dy)
	
	// Number of control points based on distance (more points for longer distances)
	numPoints := int(distance/50) + 3
	if numPoints > 20 {
		numPoints = 20
	}
	if numPoints < 5 {
		numPoints = 5
	}
	
	// Generate control points with slight randomness
	// First control point (overshoot direction)
	cp1X := start.X + dx*0.3 + (rand.Float64()-0.5)*distance*0.2
	cp1Y := start.Y + dy*0.3 + (rand.Float64()-0.5)*distance*0.2
	
	// Second control point (correction towards end)
	cp2X := end.X - dx*0.2 + (rand.Float64()-0.5)*distance*0.15
	cp2Y := end.Y - dy*0.2 + (rand.Float64()-0.5)*distance*0.15
	
	// Generate path points using cubic Bézier curve
	path := make([]Point, 0, numPoints)
	for i := 0; i <= numPoints; i++ {
		t := float64(i) / float64(numPoints)
		
		// Cubic Bézier formula: B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃
		x := math.Pow(1-t, 3)*start.X +
			3*math.Pow(1-t, 2)*t*cp1X +
			3*(1-t)*math.Pow(t, 2)*cp2X +
			math.Pow(t, 3)*end.X
		
		y := math.Pow(1-t, 3)*start.Y +
			3*math.Pow(1-t, 2)*t*cp1Y +
			3*(1-t)*math.Pow(t, 2)*cp2Y +
			math.Pow(t, 3)*end.Y
		
		// Add slight micro-movements for realism
		x += (rand.Float64() - 0.5) * 0.5
		y += (rand.Float64() - 0.5) * 0.5
		
		path = append(path, Point{X: x, Y: y})
	}
	
	// Add final correction to ensure we end exactly at target
	path = append(path, end)
	
	return path
}

// GenerateOvershootPath creates a path that overshoots the target and corrects
func GenerateOvershootPath(start, end Point) []Point {
	rand.Seed(time.Now().UnixNano())
	
	dx := end.X - start.X
	dy := end.Y - start.Y
	
	// Overshoot amount (5-15% beyond target)
	overshootFactor := 0.05 + rand.Float64()*0.1
	
	// Overshoot point
	overshoot := Point{
		X: end.X + dx*overshootFactor,
		Y: end.Y + dy*overshootFactor,
	}
	
	// Create path: start -> overshoot -> end (with correction)
	path1 := GenerateBezierPath(start, overshoot)
	path2 := GenerateBezierPath(overshoot, end)
	
	// Combine paths
	fullPath := append(path1[:len(path1)-1], path2...)
	
	return fullPath
}

