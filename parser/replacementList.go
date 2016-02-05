package parser

import (
)

type replaceObject struct {
	from        int
	to          int
	replacement string
}

type ReplacementList []replaceObject

func (dl ReplacementList) Len() int {
	return len(dl)
}

func (dl ReplacementList) Swap(i, j int) {
	dl[i], dl[j] = dl[j], dl[i]
}

func (dl ReplacementList) Less(i, j int) bool {
	return dl[i].from < dl[j].from
}
