package main

import "test/bm_poc/bm"

func GetBenchConfig() bm.Config {
	return bm.Config{
		Apps: []bm.App{
			// {User: "frappe", Repo: "frappe", Branch: "develop"},
			{User: "frappe", Repo: "erpnext", Branch: "develop"},
			{User: "frappe", Repo: "hrms", Branch: "develop"},
			{User: "frappe", Repo: "gameplan", Branch: "main"},
			// {User: "frappe", Repo: "builder", Branch: "develop"},
			// {User: "frappe", Repo: "drive", Branch: "main"},
		},
	}
}
