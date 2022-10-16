package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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

func LookupHistory(options []string, path string) (History, []string) {
	var history History
	cacheDir, err := os.UserCacheDir()
	if err == nil {
		history, err = ReadHistory(cacheDir + "/" + path)
		if err == nil {
			in, out := HistorySplit(history, options)
			SortEntriesByHistory(history, in)
			options = append(in, out...)
		} else {
			//TODO: show error message
		}
	} else {
		//TODO: show another error message
	}

	if history == nil {
		history = make(History)
	}

	return history, options
}

func AppendHistory(history History, entry string, dir string) {
	//TODO: this branch can mayyybe be avoided
	count, ok := history[entry]
	if !ok {
		history[entry] = 1
	} else {
		history[entry] = count + 1
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		//TODO: show yet another error message
	}
	err = WriteHistory(history, cacheDir + "/" + dir)
	if err != nil {
		//TODO: show yet another error message
	}
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
