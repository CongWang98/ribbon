package ribbon

import (
	"math"
	"fmt"
	"github.com/fogleman/ease"
	"github.com/fogleman/fauxgl"
	"../pdb"
)

func ellipseProfile(n int, w, h float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, n)
	for i := range result {
		t := float64(i) / float64(n)
		a := t*2*math.Pi + math.Pi/4
		x := math.Cos(a) * w / 2
		y := math.Sin(a) * h / 2
		result[i] = fauxgl.Vector{x, y, 0}
	}
	return result
}

func rectangleProfile(n int, w, h float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, 0, n)
	hw := w / 2
	hh := h / 2
	segments := [][2]fauxgl.Vector{
		{fauxgl.Vector{hw, hh, 0}, fauxgl.Vector{-hw, hh, 0}},
		{fauxgl.Vector{-hw, hh, 0}, fauxgl.Vector{-hw, -hh, 0}},
		{fauxgl.Vector{-hw, -hh, 0}, fauxgl.Vector{hw, -hh, 0}},
		{fauxgl.Vector{hw, -hh, 0}, fauxgl.Vector{hw, hh, 0}},
	}
	m := n / 4
	for _, s := range segments {
		for i := 0; i < m; i++ {
			t := float64(i) / float64(m)
			p := s[0].Lerp(s[1], t)
			result = append(result, p)
		}
	}
	return result
}

func roundedRectangleProfile(n int, w, h float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, 0, n)
	r := h / 2
	hw := w/2 - r
	hh := h / 2
	segments := [][2]fauxgl.Vector{
		{fauxgl.Vector{hw, hh, 0}, fauxgl.Vector{-hw, hh, 0}},
		{fauxgl.Vector{-hw, 0, 0}, fauxgl.Vector{}},
		{fauxgl.Vector{-hw, -hh, 0}, fauxgl.Vector{hw, -hh, 0}},
		{fauxgl.Vector{hw, 0, 0}, fauxgl.Vector{}},
	}
	m := n / 4
	for si, s := range segments {
		for i := 0; i < m; i++ {
			t := float64(i) / float64(m)
			var p fauxgl.Vector
			switch si {
			case 0, 2:
				p = s[0].Lerp(s[1], t)
			case 1:
				a := math.Pi/2 + math.Pi*t
				x := math.Cos(a) * r
				y := math.Sin(a) * r
				p = s[0].Add(fauxgl.Vector{x, y, 0})
			case 3:
				a := 3*math.Pi/2 + math.Pi*t
				x := math.Cos(a) * r
				y := math.Sin(a) * r
				p = s[0].Add(fauxgl.Vector{x, y, 0})
			}
			result = append(result, p)
		}
	}
	return result
}

func scaleProfile(p []fauxgl.Vector, s float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, len(p))
	for i := range result {
		result[i] = p[i].MulScalar(s)
	}
	return result
}

func translateProfile(p []fauxgl.Vector, dx, dy float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, len(p))
	for i := range result {
		result[i] = p[i].Add(fauxgl.Vector{dx, dy, 0})
	}
	return result
}

func segmentProfiles(pp1, pp2 *PeptidePlane, n int) (p1, p2 []fauxgl.Vector) {
	type0 := pp1.Residue1.Type
	type1, type2 := pp1.Transition()
	const ribbonWidth = 2
	const ribbonHeight = 0.125
	const ribbonOffset = 1.5
	const arrowHeadWidth = 3
	const arrowWidth = 2
	const arrowHeight = 0.5
	const tubeSize = 0.75
	offset1 := ribbonOffset
	offset2 := ribbonOffset
	if pp1.Flipped {
		offset1 = -offset1
	}
	if pp2.Flipped {
		offset2 = -offset2
	}
	switch type1 {
	case pdb.ResidueTypeHelix:
		if type0 == pdb.ResidueTypeStrand {
			p1 = roundedRectangleProfile(n, 0, 0)
		} else {
			p1 = roundedRectangleProfile(n, ribbonWidth, ribbonHeight)
		}
		p1 = translateProfile(p1, 0, offset1)
	case pdb.ResidueTypeStrand:
		if type2 == pdb.ResidueTypeStrand {
			p1 = rectangleProfile(n, arrowWidth, arrowHeight)
		} else {
			p1 = rectangleProfile(n, arrowHeadWidth, arrowHeight)
		}
	default:
		if type0 == pdb.ResidueTypeStrand {
			p1 = ellipseProfile(n, 0, 0)
		} else {
			p1 = ellipseProfile(n, tubeSize, tubeSize)
		}
	}
	switch type2 {
	case pdb.ResidueTypeHelix:
		p2 = roundedRectangleProfile(n, ribbonWidth, ribbonHeight)
		p2 = translateProfile(p2, 0, offset2)
	case pdb.ResidueTypeStrand:
		p2 = rectangleProfile(n, arrowWidth, arrowHeight)
	default:
		p2 = ellipseProfile(n, tubeSize, tubeSize)
	}
	if type1 == pdb.ResidueTypeStrand && type2 != pdb.ResidueTypeStrand {
		p2 = rectangleProfile(n, 0, arrowHeight)
	}
	return
}

func segmentColors(pp *PeptidePlane) (c1, c2 fauxgl.Color) {
	// const minTemp = 10
	// const maxTemp = 50
	// f1 := pp.Residue2.Atoms["CA"].TempFactor
	// f2 := pp.Residue3.Atoms["CA"].TempFactor
	// t1 := fauxgl.Clamp((f1-minTemp)/(maxTemp-minTemp), 0, 1)
	// t2 := fauxgl.Clamp((f2-minTemp)/(maxTemp-minTemp), 0, 1)
	// c1 = fauxgl.MakeColor(Viridis.Color(t1))
	// c2 = fauxgl.MakeColor(Viridis.Color(t2))
	// return
	type1, type2 := pp.Transition()
	switch type1 {
	case pdb.ResidueTypeHelix:
		c1 = fauxgl.HexColor("FFB733")
	case pdb.ResidueTypeStrand:
		c1 = fauxgl.HexColor("F57336")
	default:
		c1 = fauxgl.HexColor("047878")
	}
	switch type2 {
	case pdb.ResidueTypeHelix:
		c2 = fauxgl.HexColor("FFB733")
	case pdb.ResidueTypeStrand:
		c2 = fauxgl.HexColor("F57336")
	default:
		c2 = fauxgl.HexColor("047878")
	}
	if type1 == pdb.ResidueTypeStrand {
		c2 = c1
	}
	return
}

func createSegmentMesh(i, n int, pp1, pp2, pp3, pp4 *PeptidePlane) *fauxgl.Mesh {
	const splineSteps = 32
	const profileDetail = 16
	type0 := pp2.Residue1.Type
	type1, type2 := pp2.Transition()
	c1, c2 := segmentColors(pp2)
	profile1, profile2 := segmentProfiles(pp2, pp3, profileDetail)
	easeFunc := ease.Linear
	if !(type1 == pdb.ResidueTypeStrand && type2 != pdb.ResidueTypeStrand) {
		easeFunc = ease.InOutQuad
	}
	if type0 == pdb.ResidueTypeStrand && type1 != pdb.ResidueTypeStrand {
		easeFunc = ease.OutCirc
	}
	// if type1 != pdb.ResidueTypeStrand && type2 == pdb.ResidueTypeStrand {
	// 	easeFunc = ease.InOutSquare
	// }
	if i == 0 {
		profile1 = ellipseProfile(profileDetail, 0, 0)
		easeFunc = ease.OutCirc
	} else if i == n-1 {
		profile2 = ellipseProfile(profileDetail, 0, 0)
		easeFunc = ease.InCirc
	}
	splines1 := make([][]fauxgl.Vector, len(profile1))
	splines2 := make([][]fauxgl.Vector, len(profile2))
	for i := range splines1 {
		p1 := profile1[i]
		p2 := profile2[i]
		splines1[i] = splineForPlanes(pp1, pp2, pp3, pp4, splineSteps, p1.X, p1.Y)
		splines2[i] = splineForPlanes(pp1, pp2, pp3, pp4, splineSteps, p2.X, p2.Y)
	}
	var triangles []*fauxgl.Triangle
	var lines []*fauxgl.Line
	for i := 0; i < splineSteps; i++ {
		t0 := easeFunc(float64(i) / splineSteps)
		t1 := easeFunc(float64(i+1) / splineSteps)
		if i == 0 && type1 == pdb.ResidueTypeStrand && type2 != pdb.ResidueTypeStrand {
			p00 := splines1[0][i]
			p10 := splines1[profileDetail/4][i]
			p11 := splines1[2*profileDetail/4][i]
			p01 := splines1[3*profileDetail/4][i]
			triangles = triangulateQuad(triangles, p00, p01, p11, p10, c1, c1, c1, c1)
		}
		for j := 0; j < profileDetail; j++ {
			p100 := splines1[j][i]
			p101 := splines1[j][i+1]
			p110 := splines1[(j+1)%profileDetail][i]
			p111 := splines1[(j+1)%profileDetail][i+1]
			p200 := splines2[j][i]
			p201 := splines2[j][i+1]
			p210 := splines2[(j+1)%profileDetail][i]
			p211 := splines2[(j+1)%profileDetail][i+1]
			p00 := p100.Lerp(p200, t0)
			p01 := p101.Lerp(p201, t1)
			p10 := p110.Lerp(p210, t0)
			p11 := p111.Lerp(p211, t1)
			c00 := c1.Lerp(c2, t0)
			c01 := c1.Lerp(c2, t1)
			c10 := c1.Lerp(c2, t0)
			c11 := c1.Lerp(c2, t1)
			triangles = triangulateQuad(triangles, p10, p11, p01, p00, c10, c11, c01, c00)
		}
	}
	return fauxgl.NewMesh(triangles, lines)
}

func triangulateQuad(triangles []*fauxgl.Triangle, p1, p2, p3, p4 fauxgl.Vector, c1, c2, c3, c4 fauxgl.Color) []*fauxgl.Triangle {
	t1 := fauxgl.NewTriangleForPoints(p1, p2, p3)
	t1.V1.Color = c1
	t1.V2.Color = c2
	t1.V3.Color = c3
	t2 := fauxgl.NewTriangleForPoints(p1, p3, p4)
	t2.V1.Color = c1
	t2.V2.Color = c3
	t2.V3.Color = c4
	triangles = append(triangles, t1)
	triangles = append(triangles, t2)
	return triangles
}

func createChainMesh(chain *pdb.Chain) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	var planes []*PeptidePlane
	for i := 0; i < len(chain.Residues)-2; i++ {
		r1 := chain.Residues[i]
		r2 := chain.Residues[i+1]
		r3 := chain.Residues[i+2]
		plane := NewPeptidePlane(r1, r2, r3)
		if plane != nil {
			// TODO: better handling missing required atoms
			planes = append(planes, plane)
			fmt.Printf("%d", i)
			fmt.Printf("%d", i+1)
			fmt.Printf("%d\n", i+2)
		}
	}
	var previous fauxgl.Vector
	for i, p := range planes {
		if i > 0 && p.Side.Dot(previous) < 0 {
			p.Flip()
		}
		previous = p.Side
	}
	n := len(planes) - 3
	for i := 0; i < n; i++ {
		// TODO: handle ends better
		pp1 := planes[i]
		pp2 := planes[i+1]
		pp3 := planes[i+2]
		pp4 := planes[i+3]
		m := createSegmentMesh(i, n, pp1, pp2, pp3, pp4)
		fmt.Printf("%d", i)
		fmt.Printf("%d", i+1)
		fmt.Printf("%d", i+2)
		fmt.Printf("%d\n", i+3)
		mesh.Add(m)
	}
	return mesh
}
