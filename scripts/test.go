package main

import (
	"log"
	"io/ioutil"
	"fmt"
	. "github.com/fogleman/fauxgl"
	"../pdb"
	"../ribbon"
	"github.com/nfnt/resize"
	"time"
	"strings"
)

const (
	size  = 2048
	scale = 4
)

var (
	eye    = V(0, 2, 4)
	center = V(0, 0, 0)
	up     = V(0, 1, 0).Normalize()
	light  = V(0.25, 0.25, 0.75).Normalize()
)

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

func plot(pdbfile string){
	//pdb_file := "zhongjitest.pdb"
	fmt.Printf(pdbfile)
	file_name := strings.Split(pdbfile,".")[0]
	fmt.Printf("&&&&&&\n")
	fmt.Printf(file_name)
	fmt.Printf("&&&&&&\n")
	data, err := ioutil.ReadFile(pdbfile)
	if err != nil {
		log.Fatal(err)
	}
	r := strings.NewReader(string(data))
	models, err := pdb.NewReader(r).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	model := models[0]
	fmt.Printf("atoms       = %d\n", len(model.Atoms))
	fmt.Printf("residues    = %d\n", len(model.Residues))
	fmt.Printf("chains      = %d\n", len(model.Chains))
	fmt.Printf("helixes     = %d\n", len(model.Helixes))
	fmt.Printf("strands     = %d\n", len(model.Strands))
	fmt.Printf("het-atoms   = %d\n", len(model.HetAtoms))
	fmt.Printf("connections = %d\n", len(model.Connections))
	var done func()
	done = timed("generating triangle mesh")
	mesh := ribbon.ModelMesh(model)
	done()

	fmt.Printf("triangles   = %d\n", len(mesh.Triangles))

	done = timed("transforming mesh")
	m := mesh.BiUnitCube()
	done()

	done = timed("finding ideal camera position")
	c := ribbon.PositionCamera(model, m)
	done()

	done = timed("writing mesh to disk")
	mesh.SaveSTL(fmt.Sprintf(file_name + ".stl"))
	done()

	// render
	done = timed("rendering image")
	context := NewContext(int(size*scale*c.Aspect), size*scale)
	context.ClearColorBufferWith(HexColor("#FFFFFF"))
	matrix := LookAt(c.Eye, c.Center, c.Up).Perspective(c.Fovy, c.Aspect, 1, 100)
	light := c.Eye.Sub(c.Center).Normalize()
	shader := NewPhongShader(matrix, light, c.Eye)
	shader.AmbientColor = Gray(0.3)
	shader.DiffuseColor = Gray(0.9)
	context.Shader = shader
	context.DrawTriangles(mesh.Triangles)
	done()

	// save image
	done = timed("downsampling image")
	image := context.Image()
	image = resize.Resize(uint(size*c.Aspect), size, image, resize.Bilinear)
	done()

	done = timed("writing image to disk")
	SavePNG(fmt.Sprintf(file_name + ".png"), image)
	done()

}


func main() {
	folderpath := "test"
	dir, _ := ioutil.ReadDir(folderpath)
	for _, fi := range dir {
		filetype := strings.Split(fi.Name(),".")[1]
		if filetype == "pdb"{
			fmt.Printf("*****\n")
			fmt.Printf(folderpath + "/" +fi.Name())
			fmt.Printf("*****\n")
			plot(folderpath + "/" +fi.Name())
		}
	}
}
