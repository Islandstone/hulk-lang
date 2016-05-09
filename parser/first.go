package parse

import (
	// "fmt"
	tok "../tokenizer"
)

// First(A)

// A -> a = a in First(A)
// A -> B a = First(B) in First(A)
// A -> B a = a in First(A) if B is nullable
// A -> ε = ε in First(A)
// A -> B C .. Z = ε in First(A) if ε in First(B) and First(C) ... and First(Z)

func First(prods []Production) map[string][]tok.Terminal {
	firstSets := make(map[string]map[tok.Terminal]bool)
	firstMap := make(map[string][]tok.Terminal)

	for _, prod := range prods {
		firstSets[prod.Left.Name] = make(map[tok.Terminal]bool)
		firstMap[prod.Left.Name] = nil
	}

	isNullable := func(str string) bool {
		v, ok := firstSets[str][tok.EPSILON]
		return ok && v
	}

	repeat := true

	addToFirstSet := func(Left string, term tok.Terminal) {
		if _, exists := firstSets[Left][term]; !exists {
			// fmt.Println("Adding", term, "to", Left)
			firstSets[Left][term] = true
			repeat = true
		}
	}

	run := 1
	for repeat {
		// fmt.Println("Run", run)
		run += 1
		repeat = false

		for _, prod := range prods {
			if len(prod.Right) == 0 {
				continue
			}

			//if isTerminal(prod.Right[0]) {
			if prod.Right[0].Terminal {
				addToFirstSet(prod.Left.Name, prod.Right[0].Token)
			} else {
				allNullable := true

				for i, _ := range prod.Right {

					if _, ok := firstSets[prod.Right[i].Name]; ok {
						for token, _ := range firstSets[prod.Right[i].Name] {
							if token == tok.EPSILON {
								continue
							}

							addToFirstSet(prod.Left.Name, token)
						}
					}

					if !isNullable(prod.Right[i].Name) {
						allNullable = false
						break
					}

				}

				if allNullable {
					if v, exists := firstSets[prod.Left.Name][tok.EPSILON]; !exists || (exists && !v) {
						repeat = true
						firstSets[prod.Left.Name][tok.EPSILON] = true
						// fmt.Println("Adding", EPSILON, "to", prod.Left)
					}
				}
			}
		}
	}

	// fmt.Println("Done")

	// TODO: Update these in line
	for _, prod := range prods {
		for v, _ := range firstSets[prod.Left.Name] {
			firstMap[prod.Left.Name] = append(firstMap[prod.Left.Name], v)
			delete(firstSets[prod.Left.Name], v)
		}
	}

	// fmt.Println(firstMap)

	return firstMap
}
