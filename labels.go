package decl

import (
	"slices"
	"strings"
)

type Labels struct {
	labels []string
}

func (l *Labels) Append(labels string) *Labels {
	if labels == "" {
		return l
	}
	ls := strings.Split(labels, ",")
	for i := range ls {
		ls[i] = strings.TrimSpace(ls[i])
	}
	l.labels = slices.Compact(slices.Sorted(slices.Values(append(l.labels, ls...))))
	return l
}

func (l *Labels) Labeled(label string) bool {
	_, ok := slices.BinarySearch(l.labels, label)
	return ok
}
