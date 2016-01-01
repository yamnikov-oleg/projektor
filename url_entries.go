package main

func SearchUrlEntries(query string) LaunchEntriesList {
	if !IsUrl(query) {
		return nil
	}
	return LaunchEntriesList{NewUrlLaunchEntry(query)}
}
