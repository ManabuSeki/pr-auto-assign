package githubclient

import (
	"math/rand"
	"sort"
)

func (c *Client) sliceUnique(target []string) (unique []string) {
	m := map[string]bool{}

	for _, v := range target {
		if !m[v] {
			m[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}

func (c *Client) shuffle(data []string) []string {
	d := sort.StringSlice(data)
	rand.Shuffle(d.Len(), d.Swap)
	return d
}
