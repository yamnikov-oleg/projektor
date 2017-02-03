package main

func SearchUrlEntries(query string) LaunchEntriesList {
	isUrl := IsUrl(query)
	isHttpUrl := IsHttpUrl(query)

	if !isUrl && !isHttpUrl {
		return nil
	}

	if isHttpUrl {
		query = "http://" + query
	}

	entry := NewUrlLaunchEntry(query)
	if IsInHistory(entry.Cmdline) {
		return nil
	}

	return LaunchEntriesList{entry}
}
