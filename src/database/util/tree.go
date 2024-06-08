package util

// define a tree structure for tracking the parent-Child relationship of each Window
type Child struct {
	Window *string // pointer to the Window in the Dictionary
	Weight float32 // Weight of relationship to parent
}
type Parent struct {
	Window string  // Window in the Dictionary
	Child  []Child // Children of the parent
}
type Tree struct {
	Dictionary []Parent  // Dictionary of Windows
	Weight     []float32 // Weight of each Window
}

func GenerateWeightedTree(seed []string) Tree {
	var tree Tree
	for i := 0; i < len(seed); i++ {
		// calculate the relationships of each seed string
		sapling := CalculateRelationships(seed[i], 1.0)
		// merge the sapling with the tree
		tree = MergeTrees(tree, sapling)
	}
	return tree
}

// The function calculateRelationships generates a tree structure
// from a seed string. Typically this should be a single word.
func CalculateRelationships(seed string, strength float32) Tree {
	var tree Tree
	// iterate over the seed string
	for i := 0; i < len(seed); i++ {
		prefix := seed[:i+1] // slice seed from 0 to i (includes i)
		suffix := seed[i+1:] // slice seed from i to end (excludes i)
		// check if the prefix is already in the Dictionary
		pFound := false
		sFound := false
		pIndex := len(tree.Dictionary)
		sIndex := len(tree.Dictionary)
		for j := 0; j < len(tree.Dictionary); j++ {
			if prefix == tree.Dictionary[j].Window && !pFound {
				// increase the Weight of the prefix based on the strength
				tree.Weight[j] += strength
				pFound = true
				pIndex = j
			}
			if suffix == tree.Dictionary[j].Window && !sFound {
				// increase the Weight of the suffix based on the strength
				tree.Weight[j] += strength
				sFound = true
				sIndex = j
			}
			if pFound && sFound {
				break
			}
		}
		// add the prefix to the Dictionary if it is not already there
		// and initialize its associated Weight to 1.0
		if !pFound && prefix != "" {
			tree.Dictionary = append(tree.Dictionary, Parent{Window: prefix})
			tree.Weight = append(tree.Weight, 1.0)
			pIndex = len(tree.Dictionary) - 1
		}
		// add the suffix to the Dictionary if it is not already there
		// and initialize its associated Weight to 1.0
		if !sFound && suffix != "" {
			tree.Dictionary = append(tree.Dictionary, Parent{Window: suffix})
			tree.Weight = append(tree.Weight, 1.0)
			sIndex = len(tree.Dictionary) - 1
		}
		// add suffix as a Child of prefix
		cFound := false
		for j := 0; j < len(tree.Dictionary[pIndex].Child); j++ {
			// check if the suffix is already a Child of the prefix
			if tree.Dictionary[pIndex].Child[j].Window == &tree.Dictionary[sIndex].Window {
				// increase the Weight of the relationship based on the strength
				tree.Dictionary[pIndex].Child[j].Weight += strength
				cFound = true
				break
			}
		}
		// add the suffix as a Child of the prefix if it is not already there
		if !cFound && pIndex != sIndex && suffix != "" {
			tree.Dictionary[pIndex].Child = append(tree.Dictionary[pIndex].Child, Child{Window: &tree.Dictionary[sIndex].Window, Weight: 1.0})
		}
	}
	return tree
}

func MergeTrees(tree Tree, sapling Tree) Tree {
	// iterate over the sapling Dictionary
	for i := 0; i < len(sapling.Dictionary); i++ {
		found := false
		// iterate over the tree Dictionary
		for j := 0; j < len(tree.Dictionary); j++ {
			// check if the Window is already in the tree Dictionary
			if sapling.Dictionary[i].Window == tree.Dictionary[j].Window {
				// increase the Weight of the Window based on the strength
				tree.Weight[j] += sapling.Weight[i]
				// merge the Children of the sapling with the tree
				for k := 0; k < len(sapling.Dictionary[i].Child); k++ {
					cFound := false
					for l := 0; l < len(tree.Dictionary[j].Child); l++ {
						// check if the Child is already in the tree Dictionary
						if sapling.Dictionary[i].Child[k].Window == tree.Dictionary[j].Child[l].Window {
							// increase the Weight of the relationship based on the strength
							tree.Dictionary[j].Child[l].Weight += sapling.Dictionary[i].Child[k].Weight
							cFound = true
							break
						}
					}
					// add the Child to the parent if it is not already there
					if !cFound {
						tree.Dictionary[j].Child = append(tree.Dictionary[j].Child, sapling.Dictionary[i].Child[k])
					}
				}
				found = true
				break
			}
		}
		// add the Window to the tree Dictionary if it is not already there
		// and initialize its associated Weight to 1.0
		if !found {
			tree.Dictionary = append(tree.Dictionary, sapling.Dictionary[i])
			tree.Weight = append(tree.Weight, sapling.Weight[i])
		}
	}
	return tree
}
