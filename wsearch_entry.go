package main

func MakeWebSearchEntry(query string) LaunchEntriesList {
	return LaunchEntriesList{
		NewWebSearchEntry(query),
	}
}
