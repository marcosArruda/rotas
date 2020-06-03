package graph

import (
    "testing"
)

var g AirportsGraph

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

func TestAdd(t *testing.T) {
  fillGraph()
  g.String()
}
