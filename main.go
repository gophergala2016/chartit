package main

import (
	"encoding/csv"
	"flag"
	"github.com/ajstarks/svgo"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// Chart is the simplest element only make up by label and value.
type Chart struct {
	Label string
	Value int
}

// Charts is slice of Chart
type Charts []Chart

// Sum of Charts,
func (c Charts) Sum() (s int) {
	for _, col := range c {
		s += col.Value
	}
	return
}

// Percentage of Charts,
func (c Charts) Percentage(label string) float64 {
	var numerator int
	for _, col := range c {
		if strings.Compare(col.Label, label) == 0 {
			numerator = col.Value
			break
		}
	}
	return float64(numerator) / float64(c.Sum())
}

// readCSV read from filename input by cmd flag, convert CSV into Charts struct and error if any file operand error happen.
func readCSV(filename string) (c Charts, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, item := range records {
		if v, err := strconv.Atoi(item[1]); err == nil {
			col := Chart{item[0], v}
			c = append(c, col)
		}
	}
	return
}

func Draw(c Charts, width, height int, w io.Writer) (err error) {
	var angle float64
	canvas := svg.New(w)
	r := width * 3 / 10
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, r, "fill:none;stroke:red")
	canvas.Line(width/2, height/2, width/2, (height/2)-r, "stroke:red")
	for _, col := range c {
		angle += float64(360) * c.Percentage(col.Label)
		radius := angleToRadius(angle)
		canvas.Line(width/2, height/2, (width/2)-int(math.Acos(radius))*r, (height/2)-int(math.Asin(radius))*r, "stroke:red")
	}
	canvas.End()
	return
}

func angleToRadius(angle float64) (radius float64) {
	if angle > 180 {
		angle = (angle - 180) * (-1)
	}
	radius = angle * math.Pi / 180
	return
}

func main() {
	var csvFile string
	var chartFile string
	var width, height int
	flag.StringVar(&csvFile, "csv", "input.csv", "CSV filename")
	flag.StringVar(&chartFile, "output", "output.svg", "OUTPUT filename")
	flag.IntVar(&width, "width", 1000, "OUTPUT file width")
	flag.IntVar(&height, "height", 800, "OUTPUT file height")
	flag.Parse()
	c, err := readCSV(csvFile)
	if err != nil {
		log.Fatal("Read csv file with some problem!")
	}
	out, err := os.Create(chartFile)
	if err != nil {
		log.Fatal("Create file with some problem!")
	}
	if err = Draw(c, width, height, out); err != nil {
		log.Fatal("Write canvas file with some problem!")
	}
}
