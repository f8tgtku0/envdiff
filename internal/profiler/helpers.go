package profiler

import "sort"

func sortStrings(s []string) {
	sort.Strings(s)
}

func sortProfiles(profiles []KeyProfile) {
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Key < profiles[j].Key
	})
}
