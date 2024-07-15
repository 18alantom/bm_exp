package main

type BenchConfig struct {
	Target string
	Apps   []App
}
type App struct {
	User string
	Repo string
}

func GetBenchConfig() BenchConfig {
	return BenchConfig{
		Target: "/Users/alan/Desktop/code/test_go/bm_poc/bench",
		Apps: []App{
			{"frappe", "erpnext"},
			{"frappe", "hrms"},
			{"frappe", "gameplan"},
			{"frappe", "builder"},
			{"frappe", "drive"},
		},
	}
}
