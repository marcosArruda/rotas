package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bexs.marcosarruda.info/rotas/graph"
	"bexs.marcosarruda.info/rotas/rest"
)

var g *graph.AirportsGraph = &graph.AirportsGraph{
	RoutesTable: map[string]*[]*graph.Route{},
}

func main() {
	filePath := os.Args[1:][0]
	g.LoadFromDisk(filePath)
	go rest.API(filePath, g)
	time.Sleep(time.Second)
	parseOtherInputs()
}

func parseFirstInput(file *os.File) {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ",")
		origin := graph.Node{Name: res[0]}
		g.AddNode(&origin)
		destination := graph.Node{Name: res[1]}
		g.AddNode(&destination)
		cost, _ := strconv.ParseFloat(res[2], 64)
		g.AddEdge(&origin, &destination, cost)
	}
}

func parseOtherInputs() {
	r := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("please enter the route: ")
		r.Scan()
		inputSlice := strings.Split(r.Text(), "-")
		if len(inputSlice) < 2 {
			log.Fatal("Put an input in the format XXX-XXX\n")
		}
		originStr := inputSlice[0]
		destStr := inputSlice[1]
		contains := g.Contains(&graph.Node{Name: originStr}) && g.Contains(&graph.Node{Name: destStr})

		//best route: GRU - BRC - SCL - ORL - CDG > $40
		if !contains {
			log.Fatalf("Origin or Destination missing! %v", r.Text())
		}

		//calculate route;
		bestRoute, err := g.CalcBestRoute(originStr, destStr, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("best route:", bestRoute)
		//fmt.Println("best route: XXX - XXX - XXX - XXX > $40")
	}
}
