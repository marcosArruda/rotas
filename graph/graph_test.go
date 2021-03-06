package graph

import (
	"fmt"
	"testing"
)

var g = AirportsGraph{
	RoutesTable: map[string]*[]*Route{},
}

func TestAddOK(t *testing.T) {
	fillGraph()
}

func TestAddDuplicate(t *testing.T) {
	nGRU := Node{Name: "GRU"}
	g.AddNode(&nGRU)
	g.AddNode(&nGRU)
	want := true
	if got := g.ContainsStr("GRU"); got != want {
		t.Errorf("ContainsStr(...) = %v, want %v", got, want)
	}
}

func TestAllRoutesList(t *testing.T) {
	fillGraph()
	g.allRoutes("GRU", "ORL", map[string]bool{}, []string{})
	allRoutes := *g.RoutesTable["GRU-ORL"]
	fmt.Println(allRoutes[0].String(false))
}

func fillGraph() {
	/*
	  ORIGIN DESTINAION COST
	  GRU    BRC        10
	  BRC    SCL        5
	  GRU    CDG        75
	  GRU    SCL        20
	  GRU    ORL        56
	  ORL    CDG        5
	  SCL    ORL        20
	*/
	nGRU := Node{Name: "GRU"}
	nBRC := Node{Name: "BRC"}
	nSCL := Node{Name: "SCL"}
	nCDG := Node{Name: "CDG"}
	nORL := Node{Name: "ORL"}

	g.AddNode(&nGRU)
	g.AddNode(&nBRC)
	g.AddNode(&nSCL)
	g.AddNode(&nCDG)
	g.AddNode(&nORL)

	g.AddEdge(&nGRU, &nBRC, 10)
	g.AddEdge(&nBRC, &nSCL, 5)
	g.AddEdge(&nGRU, &nCDG, 75)
	g.AddEdge(&nGRU, &nSCL, 20)
	g.AddEdge(&nGRU, &nORL, 56)
	g.AddEdge(&nORL, &nCDG, 5)
	g.AddEdge(&nSCL, &nORL, 20)
}
