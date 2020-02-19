package cli

// ComptcdCategories is a slice of *ComptcdCategory.
type ComptcdCategories []*ComptcdCategory

// ComptcdCategory is a category containing comptcds.
type ComptcdCategory struct {
	Name     string
	Comptcds Comptcds
}

func (c ComptcdCategories) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

func (c ComptcdCategories) Len() int {
	return len(c)
}

func (c ComptcdCategories) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// AddComptcd adds a comptcd to a category.
func (c ComptcdCategories) AddComptcd(category string, comptcd Comptcd) ComptcdCategories {
	for _, comptcdCategory := range c {
		if comptcdCategory.Name == category {
			comptcdCategory.Comptcds = append(comptcdCategory.Comptcds, comptcd)
			return c
		}
	}
	return append(c, &ComptcdCategory{Name: category, Comptcds: []Comptcd{comptcd}})
}

// VisibleComptcds returns a slice of the Comptcds with Hidden=false
func (c *ComptcdCategory) VisibleComptcds() []Comptcd {
	ret := []Comptcd{}
	for _, comptcd := range c.Comptcds {
		if !comptcd.Hidden {
			ret = append(ret, comptcd)
		}
	}
	return ret
}
