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
			errors := vsf[index].Errors
			if v.Error {
				errors++
			}
			events := append(vsf[index].Events, v)
			trailingEvents := TrailingEvents(events, 10)
			vsf[index] = SectionDetail{
				Section:            v.Section,
				Events:             events,
				Hits:               len(events),
				HitsLastTenSeconds: len(trailingEvents),
				Errors:             errors,
			}
		} else {
			errors := 0
			if v.Error {
				errors++
			}
			events := append(make([]LogEvent, 0), v)
			trailingEvents := TrailingEvents(events, 10)
			vsf = append(vsf, SectionDetail{
				Section:            v.Section,
				Events:             events,
				Hits:               1,
				HitsLastTenSeconds: len(trailingEvents),
				Errors:             errors,
			})
		}
	}
	return vsf
}
