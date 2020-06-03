package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"bexs.marcosarruda.info/rotas/graph"
)

var g graph.AirportsGraph

func main() {
	allArgs := os.Args[1:]
	//f, err := os.Open("input-file.txt")
	f, err := os.Open(allArgs[0])

	if err != nil {
		panic(err)
	}
	parseInput(f)
	g.AllRoutes(g.Nodes[0].Name, g.Nodes[2].Name)
}

func parseInput(file *os.File) {
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//neighbors := make(map[string][]neighbor)
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

func calcBestRoute(origin string, destination string) (string, error) {
	return "", graph.BestRouteError{Text: "Shit"}
}

func finalImpl() {

	r := bufio.NewScanner(os.Stdin)
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
	a, err := calcBestRoute(originStr, destStr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
	fmt.Println("best route: XXX - XXX - XXX - XXX > $40")
	//fmt.Println("OI!", inputSlice[0], " contains:", contains)

	//g.String()

	//r := bufio.NewScanner(os.Stdin)
	//fmt.Print("Which table would you like to watch update?: ")
	//r.Scan()
	//watchTable := r.Text()

	//Node{Name: "GRU"}
	//g.AddNode(&nGRU)
	//g.AddEdge(&nGRU, &nBRC, 10)

	//fmt.Println("oi")
	//node := graph.Node{Name: "brasil"}
	//fmt.Println(node.Name)
}
