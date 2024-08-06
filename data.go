package main

import "test/bm_poc/bm"

func GetBenchConfig() bm.Config {
	return bm.Config{
		Apps: []bm.App{
			// {User: "frappe", Repo: "frappe", Branch: "develop"},
			// {User: "frappe", Repo: "erpnext", Branch: "develop"}, // doesn't build will support case later
			// {User: "frappe", Repo: "hrms", Branch: "develop"}, // requires common_site_config.json
			// {User: "frappe", Repo: "gameplan", Branch: "main"}, // requires common_site_config.json
			{User: "frappe", Repo: "builder", Branch: "develop"},
			{User: "frappe", Repo: "drive", Branch: "main"},
		},
	}
}
