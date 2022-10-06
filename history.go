package main

import (
	"encoding/json"
	"io/ioutil"
	"sort"
)

type History map[string]int

type sortEntriesByHistory struct {
	entries []string
	history History
}

func (h sortEntriesByHistory) Len() int {
	return len(h.entries)
}

func (h sortEntriesByHistory) Less(i, j int) bool {
	l := h.history[h.entries[i]]
	r := h.history[h.entries[j]]
	return l > r
}

func (h sortEntriesByHistory) Swap(i, j int) {
	h.entries[i], h.entries[j] = h.entries[j], h.entries[i]
}

func ReadHistory(path string) (History, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var history History
	err = json.Unmarshal(bytes, &history)
	return history, err
}

func WriteHistory(history History, path string) error {
	bytes, err := json.Marshal(history)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

func HistorySplit(history History, entries []string) (within []string, outside []string) {
	for _, entry := range entries {
		_, ok := history[entry]
		if ok {
			within = append(within, entry)
		} else {
			outside = append(outside, entry)
		}
	}
	return
}

func SortEntriesByHistory(history History, entries []string) {
	rule := sortEntriesByHistory{
		entries: entries,
		history: history,
	}
	sort.Sort(rule)
}
