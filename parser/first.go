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
		firstSets[prod.left.Name] = make(map[tok.Terminal]bool)
		firstMap[prod.left.Name] = nil
	}

	isNullable := func(str string) bool {
		v, ok := firstSets[str][tok.EPSILON]
		return ok && v
	}

	repeat := true

	addToFirstSet := func(left string, term tok.Terminal) {
		if _, exists := firstSets[left][term]; !exists {
			// fmt.Println("Adding", term, "to", left)
			firstSets[left][term] = true
			repeat = true
		}
	}

	run := 1
	for repeat {
		// fmt.Println("Run", run)
		run += 1
		repeat = false

		for _, prod := range prods {
			if len(prod.right) == 0 {
				continue
			}

			//if isTerminal(prod.right[0]) {
			if prod.right[0].Terminal {
				addToFirstSet(prod.left.Name, prod.right[0].Token)
			} else {
				allNullable := true

				for i, _ := range prod.right {

					if _, ok := firstSets[prod.right[i].Name]; ok {
						for token, _ := range firstSets[prod.right[i].Name] {
							if token == tok.EPSILON {
								continue
							}

							addToFirstSet(prod.left.Name, token)
						}
					}

					if !isNullable(prod.right[i].Name) {
						allNullable = false
						break
					}

				}

				if allNullable {
					if v, exists := firstSets[prod.left.Name][tok.EPSILON]; !exists || (exists && !v) {
						repeat = true
						firstSets[prod.left.Name][tok.EPSILON] = true
						// fmt.Println("Adding", EPSILON, "to", prod.left)
					}
				}
			}
		}
	}

	// fmt.Println("Done")

	// TODO: Update these in line
	for _, prod := range prods {
		for v, _ := range firstSets[prod.left.Name] {
			firstMap[prod.left.Name] = append(firstMap[prod.left.Name], v)
			delete(firstSets[prod.left.Name], v)
		}
	}

	// fmt.Println(firstMap)

	return firstMap
}
