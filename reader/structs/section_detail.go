package structs

// SectionDetail represents all the events for a particular path
type SectionDetail struct {
	Section              string
	Hits                 int
	HitsLastTenSeconds   int
	Errors               int
	ErrorsLastTenSeconds int
	Events               []LogEvent
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
			vsf[index] = SectionDetail{
				Section: v.Section,
				Events:  append(vsf[index].Events, v),
				Hits:    len(vsf[index].Events) + 1,
			}
		} else {
			vsf = append(vsf, SectionDetail{
				Section: v.Section,
				Events:  append(make([]LogEvent, 0), v),
				Hits:    1,
			})
		}
	}
	return vsf
}
