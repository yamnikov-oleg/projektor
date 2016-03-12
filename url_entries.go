package main

func SearchUrlEntries(query string) LaunchEntriesList {
	if !IsUrl(query) {
		return nil
	}
	entry := NewUrlLaunchEntry(query)
	if IsInHistory(entry.Cmdline) {
		return nil
	}
	return LaunchEntriesList{entry}
}
