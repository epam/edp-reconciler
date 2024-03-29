package job_provisioning

import (
	"database/sql"
	"fmt"
)

const (
	SelectJobProvisioningSql = "select id from \"%v\".job_provisioning where name = $1 and scope = $2;"
	InsertJobProvisioningSql = "insert into \"%v\".job_provisioning(name, scope) values ($1, $2)"
)

func SelectJobProvision(txn *sql.Tx, name string, scope string, tenant string) (*int, error) {
	stmt, err := txn.Prepare(fmt.Sprintf(SelectJobProvisioningSql, tenant))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, scope).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &id, err
}

func CreateJobProvision(txn *sql.Tx, name string, scope string, tenant string) error {
	stmt, err := txn.Prepare(fmt.Sprintf(InsertJobProvisioningSql, tenant))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, scope)

	return err
}
