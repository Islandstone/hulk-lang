package parse

import (
	tok "../tokenizer"
	// "fmt"
)

// Follow(A)

// http://www.jambe.co.nz/UNI/FirstAndFollowSets.html
// Rules for Follow Sets
//
// First put $ (the end of input marker) in Follow(S) (S is the start symbol)
// If there is a production A → aBb, (where a can be a whole string) then everything in FIRST(b) except for ε is placed in FOLLOW(B).
// If there is a production A → aB, then everything in FOLLOW(A) is in FOLLOW(B)
// If there is a production A → aBb, where FIRST(b) contains ε, then everything in FOLLOW(A) is in FOLLOW(B)

func Follow(prods []Production, first map[string][]tok.Terminal) map[string][]tok.Terminal {
	followSets := make(map[string]map[tok.Terminal]bool)
	followMap := make(map[string][]tok.Terminal)

	for _, prod := range prods {
		followSets[prod.Left.Name] = make(map[tok.Terminal]bool)
		followMap[prod.Left.Name] = nil
	}

	// All grammars have a S' which should have a $ in the follow set
	followSets["S'"][tok.EOF] = true

	repeat := true

	addToFollowSet := func(Left string, term tok.Terminal) {
		if _, exists := followSets[Left][term]; !exists {
			// fmt.Println("Adding", term, "to", Left)
			followSets[Left][term] = true
			if !repeat {
				repeat = true
			}
		}
	}

	isNullable := func(str string) bool {
		for _, e := range first[str] {
			if e == tok.EPSILON {
				return true
			}
		}
		return false
	}

	for repeat {
		repeat = false
		for _, prod := range prods {
			for i := len(prod.Right) - 1; i >= 0; i-- {
				last := prod.Right[i]

				if !last.Terminal {
					for key, _ := range followSets[prod.Left.Name] {
						addToFollowSet(last.Name, key)
					}

					if !isNullable(last.Name) {
						break
					}
				} else {
					break
				}
			}

			if len(prod.Right) <= 1 {
				continue
			}

			// All tokens except the last one in the production
			for i, this := range prod.Right[:len(prod.Right)-1] {
				if this.Terminal {
					continue
				}

				next := prod.Right[i+1]
				if next.Terminal {
					// If the next token is a terminal, just add it
					addToFollowSet(this.Name, next.Token)
				} else {
					// If the next token is a non-terminal, add its first set to
					// our follow set
					for _, f := range first[next.Name] {
						// ... except for epsilon
						if f == tok.EPSILON {
							continue
						}
						addToFollowSet(this.Name, f)
					}
				}
			}
		}
	}

	// TODO: Update these in line
	for _, prod := range prods {
		for v, _ := range followSets[prod.Left.Name] {
			followMap[prod.Left.Name] = append(followMap[prod.Left.Name], v)
			delete(followSets[prod.Left.Name], v)
		}
	}

	// fmt.Println("Follow:", followMap)

	return followMap
}
