package rest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"bexs.marcosarruda.info/rotas/graph"
)

var g *graph.AirportsGraph
var f string

//RouteWithCost is the Route in Json format. Received by the Rest API
type RouteWithCost struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Cost float64 `json:"cost"`
}

//RouteOnly is an input origin-destination to find the cheapest route
type RouteOnly struct {
	From string `json:"from"`
	To   string `json:"to"`
}

//BestRouteRespose is the response for a bestRoute request
type BestRouteRespose struct {
	Best string `json:"best"`
	Cost string `json:"cost"`
}

func registerRoute(w http.ResponseWriter, r *http.Request) {
	s := "Kindly enter the route with the 'from', 'to' and 'cost' only in order to register a route"
	var newRoute RouteWithCost
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, s)
	}

	errMarshalling := json.Unmarshal(reqBody, &newRoute)
	if errMarshalling != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, s)
	} else {
		newRoute.From = strings.ToUpper(newRoute.From)
		newRoute.To = strings.ToUpper(newRoute.To)
		updateSourceFile(&newRoute)
		g.ResetGraph()
		g.LoadFromDisk(f)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newRoute)
	}
}

func bestRoute(w http.ResponseWriter, r *http.Request) {
	from := mux.Vars(r)["from"]
	to := mux.Vars(r)["to"]
	best, errCalc := g.CalcBestRoute(strings.ToUpper(from), strings.ToUpper(to), true)
	if errCalc != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, errCalc.Error())
	} else {
		res := strings.Split(best, "|")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BestRouteRespose{
			Best: res[0],
			Cost: res[1],
		})
	}
}

//API this is the primary function to start the Http Server behind the Rest API
func API(fr string, gr *graph.AirportsGraph) {
	g = gr
	f = fr
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/route", registerRoute).Methods("POST", "PUT")
	router.HandleFunc("/route/{from}/{to}", bestRoute).Methods("GET")
	fmt.Println("call POST/PUT http://localhost:8080/route ::: to register a new route;")
	fmt.Println("call GET http://localhost:8080/route/{from}/{to} ::: to get the best route for a target and destination;")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func updateSourceFile(newRoute *RouteWithCost) {
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, ",")
		if res[0] == newRoute.From && res[1] == newRoute.To {
			replaceInSourceFile(res[0]+","+res[1], fmt.Sprintf("%.2f", newRoute.Cost/100))
			file.Close()
			return
		} else if res[1] == newRoute.From && res[0] == newRoute.To {
			replaceInSourceFile(res[1]+","+res[0], fmt.Sprintf("%.2f", newRoute.Cost/100))
			file.Close()
			return
		}
	}
	file.Close()
	addToSourceFile(newRoute.From + "," + newRoute.To + "," + fmt.Sprintf("%.2f", newRoute.Cost))
	return
}

func addToSourceFile(str string) {
	input, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	lines = append(lines, str)
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(f, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func replaceInSourceFile(str string, cost string) {
	input, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, str) {
			lines[i] = str + "," + cost
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(f, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
