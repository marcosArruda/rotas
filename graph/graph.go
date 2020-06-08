package graph

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
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

//Route is a final Route found on the algoritm
type Route struct {
	PathList *[]string
	Cost     money.USD
}

func (r *Route) String(isWeb bool) string {
	s := ""
	for _, x := range *r.PathList {
		if s == "" {
			s += x
		} else {
			s += " - " + x
		}
	}
	if isWeb {
		return s + "|" + r.Cost.String()
	} else {
		return s + " > " + r.Cost.String()
	}

}

// AirportsGraph the Airports Graph
type AirportsGraph struct {
	Nodes       []*Node
	edges       map[string][]*Neighbor
	Prices      map[string]money.USD
	lock        sync.RWMutex
	RoutesTable map[string]*[]*Route
}

//BestRouteError error calcaulating best route
type BestRouteError struct {
	Text string
}

//RoutesTable is table of all routes
//var RoutesTable = map[string]*[]*Route{}

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

//ResetGraph remove all data to prepare for a new load from disk
func (g *AirportsGraph) ResetGraph() {
	g.lock.Lock()
	g.RoutesTable = map[string]*[]*Route{}
	g.Prices = map[string]money.USD{}
	g.edges = map[string][]*Neighbor{}
	g.Nodes = nil
	g.lock.Unlock()
}

//LoadFromDisk Load data from predefined file
func (g *AirportsGraph) LoadFromDisk(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ",")
		origin := Node{Name: res[0]}
		g.AddNode(&origin)
		destination := Node{Name: res[1]}
		g.AddNode(&destination)
		cost, _ := strconv.ParseFloat(res[2], 64)
		g.AddEdge(&origin, &destination, cost)
	}
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

func (g *AirportsGraph) CalcBestRoute(origin string, destination string, isWeb bool) (string, error) {
	contains := g.Contains(&Node{Name: origin}) && g.Contains(&Node{Name: destination})
	if !contains {
		return "", BestRouteError{Text: "Origin(" + origin + ") or Destination(" + destination + ") missing!"}
	}
	g.allRoutes(origin, destination, map[string]bool{}, []string{})
	allRoutes := *g.RoutesTable[origin+"-"+destination]
	return allRoutes[0].String(isWeb), nil
}

// AllRoutes find all routes to destination
func (g *AirportsGraph) allRoutes(from string, to string, alreadySeen map[string]bool, localPathList []string) {
	alreadySeen[from] = true
	if from == to { //found a route
		sx := ""
		localPathList = RemoveDup(localPathList)
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
		route := &Route{
			PathList: &pathList,
			Cost:     finalPrice,
		}
		routeTableName := pathList[0] + "-" + pathList[len(pathList)-1]
		// TODO: pegar este slice abaixo e adicionar e ordenar...
		if g.RoutesTable[routeTableName] == nil {
			newRouteSlc := &[]*Route{}
			g.RoutesTable[routeTableName] = newRouteSlc
		}
		routeList := g.RoutesTable[routeTableName]
		appendedSlice := append(*routeList, route)
		sort.Slice(appendedSlice, func(i, j int) bool {
			return appendedSlice[i].Cost < appendedSlice[j].Cost
		})
		*g.RoutesTable[routeTableName] = appendedSlice
		//fmt.Println(sx + " = " + finalPrice.String())
		alreadySeen[from] = false
		return
	}

	localPathList = append(localPathList, from)
	adjs := g.edges[from]
	for j := 0; j < len(adjs); j++ {
		curr := adjs[j].Node.Name
		if alreadySeen[curr] == false {
			localPathList = append(localPathList, curr)
			g.allRoutes(curr, to, alreadySeen, localPathList)
			localPathList = RemoveStr(localPathList, curr)
		}
	}
	alreadySeen[from] = false
}

//removeDup remove duplications from slice
func RemoveDup(elements []string) []string {
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
