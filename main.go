package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

// Player is the minimum representation that is needed to calculate teams
type Player struct {
	ID       int     `json:"id"`
	Strength float64 `json:"strength"`
}

// Team represents a group of Groups
type Team struct {
	Goalies  []Player `json:"goalies"`
	Centers  []Player `json:"centers"`
	Wings    []Player `json:"wings"`
	Forwards []Player `json:"forwards"`
}

// Teams holds 2 Team
type Teams struct {
	A Team `json:"team_a"`
	B Team `json:"team_b"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	team := dummyTeam()

	centers := team.Centers
	goalies := team.Goalies
	wings := team.Wings
	forwards := team.Forwards

	teams := make(chan Teams, 200)

	go func() {
		for cis := range GenCombinations(len(centers), len(centers)/2) {
			for gis := range GenCombinations(len(goalies), len(goalies)/2) {
				for wis := range GenCombinations(len(wings), len(wings)/2) {
					for fis := range GenCombinations(len(forwards), len(forwards)/2) {
						cA, cB := getPlayersAndComplementer(centers, cis)
						gA, gB := getPlayersAndComplementer(goalies, gis)
						wA, wB := getPlayersAndComplementer(wings, wis)
						fA, fB := getPlayersAndComplementer(forwards, fis)

						t1 := makeTeam(gA, cA, wA, fA)
						t2 := makeTeam(cB, gB, wB, fB)

						ts := Teams{t1, t2}

						teams <- ts
					}
				}
			}
		}
		close(teams)
	}()

	best := getBestBalanced(teams)

	b, err := json.Marshal(best)

	checkError(err)

	checkError(ioutil.WriteFile("output.json", b, 0644))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (team *Team) size() int {
	return len(team.Centers) + len(team.Goalies) + len(team.Wings) + len(team.Forwards)
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getBestBalanced(tch chan Teams) Teams {
	smallest := 10.0

	var best Teams
	for team := range tch {
		if diff := absInt(team.A.size() - team.B.size()); diff > 1 {
			tmp := team.A.Wings
			team.A.Wings = team.B.Wings
			team.B.Wings = tmp
		}
		if diff := math.Abs(team.A.averageStrenght() - team.B.averageStrenght()); diff < smallest {
			smallest = diff
			best = team
		}
	}
	return best
}

func (team *Team) printNicer() {
	fmt.Printf("%v\n%v\n\t\t%v\n\t\t%v\n", team.Forwards, team.Wings, team.Centers, team.Goalies)
}

func getPlayersAndComplementer(from []Player, selected []int) ([]Player, []Player) {
	groupA := make([]Player, len(selected))
	groupB := make([]Player, len(from)-len(selected))

	var idxSelected, idxGroupA, idxGroupB int

	for i, player := range from {
		if idxSelected >= len(selected) || i != selected[idxSelected] {
			groupB[idxGroupB] = player
			idxGroupB++
		} else {
			groupA[idxGroupA] = player
			idxGroupA++

			idxSelected++
		}
	}

	return groupA, groupB
}

func makeTeam(cs, gs, ws, fs []Player) Team {
	return Team{cs, gs, ws, fs}
}

func (team *Team) averageStrenght() float64 {
	return avgStrength(team.Centers) + avgStrength(team.Goalies) + avgStrength(team.Wings) + avgStrength(team.Forwards)
}

func avgStrength(players []Player) float64 {
	return sumStrength(players) / float64(len(players))
}

func sumStrength(players []Player) float64 {
	var sum float64
	for _, p := range players {
		sum += p.Strength
	}
	return sum
}

func dummyTeam() Team {
	var t Team

	p1 := Player{1, float64(rand.Intn(10) + 1)}
	p2 := Player{2, float64(rand.Intn(10) + 1)}
	p3 := Player{3, float64(rand.Intn(10) + 1)}
	p4 := Player{4, float64(rand.Intn(10) + 1)}
	p5 := Player{5, float64(rand.Intn(10) + 1)}
	p6 := Player{6, float64(rand.Intn(10) + 1)}
	p7 := Player{7, float64(rand.Intn(10) + 1)}
	p8 := Player{8, float64(rand.Intn(10) + 1)}
	p9 := Player{9, float64(rand.Intn(10) + 1)}
	p10 := Player{10, float64(rand.Intn(10) + 1)}
	p11 := Player{11, float64(rand.Intn(10) + 1)}
	p12 := Player{12, float64(rand.Intn(10) + 1)}
	p13 := Player{13, float64(rand.Intn(10) + 1)}
	p14 := Player{14, float64(rand.Intn(10) + 1)}
	p15 := Player{15, float64(rand.Intn(10) + 1)}
	p16 := Player{16, float64(rand.Intn(10) + 1)}
	p17 := Player{17, float64(rand.Intn(10) + 1)}
	// p18 := Player{18, float64(rand.Intn(10) + 1)}
	// p19 := Player{19, float64(rand.Intn(10) + 1)}
	p20 := Player{20, float64(rand.Intn(10) + 1)}

	t.Centers = []Player{p1, p2}
	t.Goalies = []Player{p3, p4}
	t.Wings = []Player{p5, p6, p7, p8, p13, p14, p20}
	t.Forwards = []Player{p9, p10, p11, p12, p15, p16, p17}

	return t
}
