package parse

import (
	tok "../tokenizer"
	"fmt"
)

type Automaton struct {
	prods  []Production
	items  []Item
	states []*State
}

type Shift struct {
	Goto int
}

type Reduce struct {
	// Rule  int
	Rule        string
	Count       int
	Nonterminal string
	Production
}

type Accept struct{}

func (a *Automaton) getStateWithCore(core []int) int {
outer:
	for stateId, state := range a.states {
		if state.coreLen == len(core) {
			for i := 0; i < len(core); i++ {
				if core[i] != state.itemIds[i] {
					continue outer
				}
			}

			return stateId
		}
	}

	return -1
}

type Item struct {
	Production
	position int
}

func (it Item) String() string {
	str := it.Left.Name + " → "
	for i, s := range it.Right {
		if i == it.position {
			str += "."
		}

		if s.Terminal {
			str += "<" + s.Token.String() + ">"
		} else {
			str += s.Name + " "
		}
	}

	if it.position > len(it.Right) {
		str += "."
	}

	return str
}

func NewItem(prod Production) Item {
	return Item{prod, 0}
}

func (a *Automaton) Init(prods []Production) {
	a.prods = prods
	for _, prod := range prods {
		for i := range prod.Right {
			a.items = append(a.items, Item{prod, i})
		}

		a.items = append(a.items, Item{prod, len(prod.Right) + 1})
	}
}

func (a *Automaton) getClosure(i int) map[int]bool {

	return nil
}

type State struct {
	// itemIds map[int]bool
	id      int
	itemIds []int
	coreLen int
	neigh   []*State
	edges   map[ProductionElement]int
}

func (a *Automaton) BuildState(startItems []int) *State {
	s := &State{}
	s.coreLen = len(startItems)

	s.itemIds = append(s.itemIds, startItems...)
	itemSet := make(map[int]bool)

	for _, v := range startItems {
		itemSet[v] = true
	}

	s.edges = make(map[ProductionElement]int)

	repeat := true
	for repeat {
		repeat = false

		// For each item in the set
		for id, _ := range itemSet {
			if a.items[id].position > len(a.items[id].Right) {
				continue
			}

			// Check if the next token is a non-terminal
			t := a.items[id].Right[a.items[id].position]
			if !t.Terminal {
				// Add any item with t as Left hand side and position 0
				for i, item := range a.items {
					if item.Left == t && item.position == 0 {
						if _, ok := itemSet[i]; !ok {
							itemSet[i] = true
							repeat = true
							s.itemIds = append(s.itemIds, i)
						}
					}
				}
			} else if t.Token == tok.EPSILON {
				if _, ok := itemSet[id+1]; !ok {
					itemSet[id+1] = true
					repeat = true
					s.itemIds = append(s.itemIds, id+1)
				}

			}
		}
	}

	s.id = len(a.states)
	a.states = append(a.states, s)

	edges := make(map[ProductionElement][]int)

	// Construct the cores of the neighbours
	for id, _ := range itemSet {
		if a.items[id].position > len(a.items[id].Right) {
			continue
		}

		// Do not construct an edge if the next position is an epsilon
		// Do not construct an edge if the next position is an epsilon
		if a.items[id].Right[a.items[id].position].Token == tok.EPSILON {
			continue
		}

		edges[a.items[id].Right[a.items[id].position]] =
			append(edges[a.items[id].Right[a.items[id].position]], id+1)
	}

	// For each of the neighbor cores
	for e, v := range edges {
		if n := a.getStateWithCore(v); n == -1 {
			newState := a.BuildState(v)
			s.neigh = append(s.neigh, newState)
			s.edges[e] = newState.id
		} else {
			s.neigh = append(s.neigh, a.states[n])
			s.edges[e] = n
		}
	}

	return s
}

func (a *Automaton) BuildTable() (table map[int]map[tok.Terminal]interface{}, gotoTable map[int]map[string]int) {
	table = make(map[int]map[tok.Terminal]interface{})
	gotoTable = make(map[int]map[string]int)

	firstSet := First(a.prods)
	followSet := Follow(a.prods, firstSet)

	for stateId, state := range a.states {
		table[stateId] = make(map[tok.Terminal]interface{})
		gotoTable[stateId] = make(map[string]int)

		for _, id := range state.itemIds {
			// A reduction
			if a.items[id].position > len(a.items[id].Right) {

				for _, t := range followSet[a.items[id].Left.Name] {
					if _, notEmpty := table[stateId][t]; notEmpty {
						fmt.Println("Conflict")
					} else if a.items[id].Left.Name == "S'" {
						table[stateId][t] = Accept{}
					} else {
						// Epsilons are not present on the stack by definition,
						// avoid popping too much
						epsilonCount := 0
						for _, prodElem := range a.items[id].Right {
							if prodElem.Token == tok.EPSILON {
								epsilonCount++
							}
						}
						table[stateId][t] = Reduce{a.items[id].String(), len(a.items[id].Right) - epsilonCount, a.items[id].Left.Name, a.items[id].Production}
					}
				}

				continue
			}

			t := a.items[id].Right[a.items[id].position]

			if t.Terminal {
				if t.Token != tok.EPSILON {
					table[stateId][t.Token] = Shift{state.edges[t]}
				} else {

				}
			} else {
				gotoTable[stateId][t.Name] = state.edges[t]
			}
		}
	}

	return
}
