package db

import "database/sql"
import "log"

import _ "github.com/lib/pq"

import "github.com/moio/mgr-db-locks/dbconf"

// Transaction represents a running transaction in the database
type Transaction struct {
	Pid int32
	Sql string
}

// Block represents a dependency relationship between Transactions
type Block struct {
	Blocked  *Transaction
	Blocking *Transaction
}

// Blocks returns a slice of relationships between Transactions that block one another
func Blocks() []*Block {
	result := make([]*Block, 0)
	d := dbconf.New()

	db, err := sql.Open("postgres", d.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	sql := `
SELECT blocked_locks.pid AS blockedPid,
  blocked_activity.query AS blockedSql,
  blocking_locks.pid AS blockingPid,
  blocking_activity.query AS blockingSql
FROM pg_catalog.pg_locks blocked_locks
  JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
  JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.DATABASE IS NOT DISTINCT FROM blocked_locks.DATABASE
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
  JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.GRANTED;`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var blockedPid int32
		var blockedSql string
		var blockingPid int32
		var blockingSql string
		err := rows.Scan(&blockedPid, &blockedSql, &blockingPid, &blockingSql)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, &Block{&Transaction{blockedPid, blockedSql}, &Transaction{blockingPid, blockingSql}})
	}

	return result
}
