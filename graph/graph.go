package graph

import (
	"fmt"
	"sync"

	"bexs.marcosarruda.info/rotas/money"
)

// Node a single node that composes the tree
type Node struct {
	Name string
}

// Neighbor is a reference to a neighbor Node element and his distance cost
type Neighbor struct {
	Node *Node
	Cost money.USD
}

// AirportsGraph the Airports Graph
type AirportsGraph struct {
	Nodes []*Node
	edges map[string][]*Neighbor
	lock  sync.RWMutex
}

// Route is a calculated route with relative cost
type Route struct {
	Origin      string
	Destination string
	Map         string
	Cost        money.USD
}

// RouteTable is a table of calculated routes for some origin
type RouteTable struct {
	Origin string
	Routes *[]Route
	lock   sync.RWMutex
}

//BestRouteError error calcaulating best route
type BestRouteError struct {
	Text string
}

func (e BestRouteError) Error() string {
	return e.Text
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.Name)
}

// AddNode adds a node to the graph
func (g *AirportsGraph) AddNode(n *Node) {
	g.lock.Lock()
	if !g.Contains(n) {
		g.Nodes = append(g.Nodes, n)
	}
	g.lock.Unlock()
}

//Contains verify if the Graph table contains a Node
func (g *AirportsGraph) Contains(n *Node) bool {
	for _, x := range g.Nodes {
		if n.Name == x.Name {
			return true
		}
	}
	return false
}

//ContainsStr same as Contains but with string
func (g *AirportsGraph) ContainsStr(n string) bool {
	for _, x := range g.Nodes {
		if n == x.Name {
			return true
		}
	}
	return false
}

// AddEdge adds an edge to the graph
func (g *AirportsGraph) AddEdge(n1 *Node, n2 *Node, cost float64) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[string][]*Neighbor)
	}
	g.edges[n1.Name] = append(g.edges[n1.Name], &Neighbor{Node: n2, Cost: money.ToUSD(cost)})
	g.edges[n2.Name] = append(g.edges[n2.Name], &Neighbor{Node: n1, Cost: money.ToUSD(cost)})
	g.lock.Unlock()
}

//FindRoutes fuck
func FindRoutes(g *AirportsGraph, n *Node, t *Node) {
	if g.Contains(n) {
		//x = visinho
		for _, x := range g.edges[n.Name] {
			if x.Node.Name == t.Name {
				fmt.Printf("%v - %v > %v", n.Name, x.Node.Name, x.Cost.String())
			} else {
				for _, z := range g.edges[x.Node.Name] {
					if z.Node.Name == t.Name {
						fmt.Printf("%v - %v - %v > %v\n", n.Name, x.Node.Name, z.Node.Name, x.Cost.Sum(z.Cost).String())
					}
				}
			}
		}
	}
}

//MapRoutesFor return the cheapest best route
func (g *AirportsGraph) MapRoutesFor(n *Node) {
	/*
		g.lock.Lock()
		for _, x := range g.nodes {
			if n.Name == x.Name {
				rt := &RouteTable{
					Origin: x.Name,
					lock:   sync.RWMutex{},
				}
				*rt.Routes = append(*rt.Routes, Route{
					Origin:      x.Name,
					Destination: x.Name,
					Cost:        money.ToUSD(0),
					Map:         fmt.Sprintf("%v - %v > $%v", x.Name, x.Name, money.ToUSD(0).String()),
				})

				for _, neighbor := range g.edges[x.Name] {
					r := Route{
						Origin:      n.Name,
						Destination: neighbor.Node.Name,
						Cost:        neighbor.Cost,
						Map:         fmt.Sprintf("%v - %v > $%v", n.Name, neighbor.Node.Name, neighbor.Cost.String()),
					}
					*rt.Routes = append(*rt.Routes, r)

				}
				rt.lock.Unlock()
				break
			}

		}
		//...
		g.lock.Unlock()
	*/
}

//SendToNeighbors fuck
func (g *AirportsGraph) SendToNeighbors(rt *RouteTable, neighbor *Neighbor) {
	/*
		//g.neighbor.Node
		for _, x := range g.nodes {
			if neighbor.Node.Name == x.Name {
				for _, r := range *rt.Routes {
					r.
				}
				break
			}
		}
	*/
}

// String prints stuff
func (g *AirportsGraph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.Nodes); i++ {
		s += g.Nodes[i].String() + " -> "
		near := g.edges[g.Nodes[i].Name]
		for j := 0; j < len(near); j++ {
			s += near[j].Node.String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}


func printAllPathsUtil(u string, d string, isVisited []bool, localPathList []int) {
  // Mark the current node
  isVisited[u] = true;
  if (u.equals(d))
  {
    System.out.println(localPathList);
    // if match found then no need to traverse more till depth
    isVisited[u]= false;
    return;
  }
  // Recur for all the vertices
  // adjacent to current vertex
  for (Integer i : adjList[u]){
    if (!isVisited[i]){
      // store current node
      // in path[]
      localPathList.add(i);
      printAllPathsUtil(i, d, isVisited, localPathList);
      // remove current node
      // in path[]
      localPathList.remove(i);
    }
  }
  // Mark the current node
  isVisited[u] = false;
}

// find all routes to
func (g *AirportsGraph) AllRoutes(from string, to string) {
	fmt.Println("from:", from, "to:", to)

	g.lock.RLock()
	s := ""
	alreadySeen := make([]string, 0)

	s += from + " -> "
	alreadySeen = append(alreadySeen, from)
	near := g.edges[from]
	for j := 0; j < len(near); j++ {
		if Already(alreadySeen, near[j].Node.String()) == false {
			s += near[j].Node.String() + " - "
			alreadySeen = append(alreadySeen, near[j].Node.String())
		}
	}
	s += "\n"

	fmt.Println(s)
	g.lock.RUnlock()
}

//Already fucj
func Already(source []string, thing string) bool {
	for i := 0; i < len(source); i++ {
		if source[i] == thing {
			return true
		}
	}
	return false
}
