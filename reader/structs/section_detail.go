package structs

import (
	"sort"
)

// SectionDetail represents all the events for a particular path
type SectionDetail struct {
	Section string
	Hits    int
	Errors  int
	Events  []LogEvent
}

// SectionDetailBy is the type of a "less" function that defines the ordering of its arguments.
type SectionDetailBy func(p1, p2 *SectionDetail) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by SectionDetailBy) Sort(details []SectionDetail) {
	es := &sectionDetailSorter{
		details: details,
		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(es)
}

type sectionDetailSorter struct {
	details []SectionDetail
	by      func(p1, p2 *SectionDetail) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *sectionDetailSorter) Len() int {
	return len(s.details)
}

// Swap is part of sort.Interface.
func (s *sectionDetailSorter) Swap(i, j int) {
	s.details[i], s.details[j] = s.details[j], s.details[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *sectionDetailSorter) Less(i, j int) bool {
	return s.by(&s.details[i], &s.details[j])
}

func findSectionDetail(details []SectionDetail, section string) int {
	for i, detail := range details {
		if detail.Section == section {
			return i
		}
	}
	return -1
}

// GroupBySection returns the slice of LogEvents matching filter
func GroupBySection(vs []LogEvent) []SectionDetail {
	vsf := make([]SectionDetail, 0)
	for _, v := range vs {
		index := findSectionDetail(vsf, v.Section)

		if index >= 0 {
			errors := vsf[index].Errors
			if v.Error {
				errors++
			}
			events := append(vsf[index].Events, v)

			vsf[index] = SectionDetail{
				Section: v.Section,
				Events:  events,
				Hits:    len(events),
				Errors:  errors,
			}
		} else {
			errors := 0
			if v.Error {
				errors++
			}
			events := append(make([]LogEvent, 0), v)
			vsf = append(vsf, SectionDetail{
				Section: v.Section,
				Events:  events,
				Hits:    1,
				Errors:  errors,
			})
		}
	}
	return vsf
}

// SortByHits returns the log events sorted by section
func SortByHits(details []SectionDetail) []SectionDetail {
	hits := func(p1, p2 *SectionDetail) bool {
		return p1.Hits < p2.Hits
	}

	SectionDetailBy(hits).Sort(details)
	return details
}

// SortByHitsDesc returns the log events sorted by section
func SortByHitsDesc(details []SectionDetail) []SectionDetail {
	hits := func(p1, p2 *SectionDetail) bool {
		return p1.Hits > p2.Hits
	}

	SectionDetailBy(hits).Sort(details)
	return details
}
