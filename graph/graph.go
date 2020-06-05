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

type Route struct {
	PathList *[]string
	Cost     money.USD
}

// AirportsGraph the Airports Graph
type AirportsGraph struct {
	Nodes  []*Node
	edges  map[string][]*Neighbor
	Prices map[string]money.USD
	lock   sync.RWMutex
}

//BestRouteError error calcaulating best route
type BestRouteError struct {
	Text string
}

//RoutesTable is table of all routes
var RoutesTable map[string]*[]*Route

func (e BestRouteError) Error() string {
	return e.Text
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
	if g.Prices == nil {
		g.Prices = make(map[string]money.USD)
	}
	g.Prices[n1.Name+"-"+n2.Name] = money.ToUSD(cost)
	g.Prices[n2.Name+"-"+n1.Name] = money.ToUSD(cost)
	g.edges[n1.Name] = append(g.edges[n1.Name], &Neighbor{Node: n2, Cost: money.ToUSD(cost)})
	g.edges[n2.Name] = append(g.edges[n2.Name], &Neighbor{Node: n1, Cost: money.ToUSD(cost)})
	g.lock.Unlock()
}

// AllRoutes find all routes to destination
func (g *AirportsGraph) AllRoutes(from string, to string, alreadySeen map[string]bool, localPathList []string) {
	alreadySeen[from] = true

	if from == to { //found a route
		sx := ""
		localPathList = removeDup(localPathList)
		lastStop := ""
		var finalPrice money.USD = money.ToUSD(0)
		pathList := make([]string, 0)
		for _, x := range localPathList {
			if lastStop == "" {
				sx += x
				pathList = append(pathList, x)
			} else {
				sx += " -" + g.Prices[lastStop+x].String() + "-> " + x
				pathList = append(pathList, x)
				finalPrice = finalPrice.Sum(g.Prices[lastStop+x])
			}
			lastStop = x + "-"
		}
		route := Route{
			PathList: &pathList,
			Cost:     finalPrice,
		}
		// TODO: pegar este slice abaixo e adicionar e ordenar...
		RoutesTable[pathList[0]+"-"+pathList[len(pathList)-1]]
		fmt.Println(sx + " = " + finalPrice.String())
		alreadySeen[from] = false
		return
	}

	localPathList = append(localPathList, from)
	adjs := g.edges[from]
	for j := 0; j < len(adjs); j++ {
		curr := adjs[j].Node.Name
		if alreadySeen[curr] == false {
			localPathList = append(localPathList, curr)
			g.AllRoutes(curr, to, alreadySeen, localPathList)
			localPathList = RemoveStr(localPathList, curr)
		}
	}
	alreadySeen[from] = false
}

func (g *AirportsGraph) BestCostEffectiveRoute(from string, to string, alreadySeen map[string]bool, localPathList []string) {
	alreadySeen[from] = true
	if from == to {
		sx := ""
		localPathList = removeDup(localPathList)
		lastStop := ""
		var finalPrice money.USD = money.ToUSD(0)
		for _, x := range localPathList {
			if lastStop == "" {
				sx += x
			} else {
				sx += " -" + g.Prices[lastStop+x].String() + "-> " + x
				finalPrice = finalPrice.Sum(g.Prices[lastStop+x])
			}
			lastStop = x + "-"
		}
		fmt.Println(sx + " = " + finalPrice.String())
		alreadySeen[from] = false
		return
	}
	localPathList = append(localPathList, from)
	adjs := g.edges[from]
	for j := 0; j < len(adjs); j++ {
		curr := adjs[j].Node.Name
		if alreadySeen[curr] == false {
			localPathList = append(localPathList, curr)
			g.AllRoutes(curr, to, alreadySeen, localPathList)
			localPathList = RemoveStr(localPathList, curr)
		}
	}
	alreadySeen[from] = false
}

func removeDup(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

//RemoveStr removes string from slice
func RemoveStr(a []string, ix string) []string {
	for i := 0; i < len(a); i++ {
		if a[i] == ix {
			copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
			a[len(a)-1] = ""     // Erase last element (write zero value).
			a = a[:len(a)-1]
			return a
		}
	}
	return a
}
