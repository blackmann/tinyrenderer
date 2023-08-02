package lib

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Model struct {
	Objects []Object
	Vertices []Vector3
}

type Object struct {
	Name     string
	Faces    [][]int64
}

func LoadModelWithScale(path string, scale int) (*Model, error) {
	var model Model
	var object Object

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	appendObj := func() {
		if len(object.Faces) > 0 {
			model.Objects = append(model.Objects, object)
		}
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.Split(row, " ")

		command := parts[0]

		switch command {
		case "o":
			appendObj()
			name := parts[1]
			object = Object{Name: name}

		case "v":
			x, _ := strconv.ParseFloat(parts[1], 64)
			y, _ := strconv.ParseFloat(parts[2], 64)
			z, _ := strconv.ParseFloat(parts[3], 64)

			// y up
			vertex := Vector3{X: x, Y: y * -1, Z: z}.Scale(float64(scale))
			model.Vertices = append(model.Vertices, vertex)

		case "f":
			var faces []int64
			for _, indexGroup := range parts[1:] {
				split := strings.Split(indexGroup, "/")
				index, _ :=  strconv.ParseInt(split[0], 10, 64)
				faces = append(faces, index)
			}

			object.Faces = append(object.Faces, faces)
		}
	}

	appendObj()

	return &model, nil
}

func LoadModel(path string) (*Model, error) {
	return LoadModelWithScale(path, 200)
}
