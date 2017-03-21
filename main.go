package main

import "fmt"
import "strings"

import "github.com/moio/mgr-db-locks/db"
import "github.com/moio/mgr-db-locks/process"

func main() {
	blocks := db.Blocks()
	processes := process.InterestingProcesses()

	fmt.Println("digraph blocks {")

	transactions := make(map[int32]string)
	for _, block := range blocks {
		transactions[block.Blocked.Pid] = strings.Replace(block.Blocked.Sql, "\n", " ", -1)
		transactions[block.Blocking.Pid] = strings.Replace(block.Blocking.Sql, "\n", " ", -1)
	}

	for pid, sql := range transactions {
		fmt.Printf("  %d[label=\"%d (%s)\" tooltip=\"%s\"];\n", pid, pid, processes[pid].Client.Kind, sql)
	}

	for _, block := range blocks {
		fmt.Printf("  %d -> %d [label=\"waits on\"];\n", block.Blocked.Pid, block.Blocking.Pid)
	}
	fmt.Println("}")
}
