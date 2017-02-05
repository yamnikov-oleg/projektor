package main

func MakeWebSearchEntry(query string) LaunchEntriesList {
	if query == "" {
		return nil
	}
	return LaunchEntriesList{
		NewWebSearchEntry(query),
	}
}
