package main

func SearchCmdEntries(query string) (list LaunchEntriesList) {
	if query != "" && !IsInHistory(query) {
		return LaunchEntriesList{NewEntryFromCommand(query)}
	}
	return nil
}
