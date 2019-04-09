package repository

import (
	"business-app-reconciler-controller/pkg/model"
	"database/sql"
	"time"
)

const (
	CheckDuplicateCodebaseBranchActionLog = "select cb.id " +
		"from codebase_branch as cb " +
		"		left join codebase_branch_action_log cbal on cb.id = cbal.codebase_branch_id " +
		"		left join action_log al on cbal.action_log_id = al.id " +
		"WHERE cb.name = $1 " +
		"	AND al.event = $2 " +
		"	AND al.updated_at = $3 " +
		"order by al.updated_at desc " +
		"limit 1;"
	InsertCodebaseActionLog = "insert into action_log(event, detailed_message, username, updated_at) " +
		"VALUES($1, $2, $3, $4) returning id;"
	InsertCodebaseBranchActionLog = "insert into codebase_branch_action_log(codebase_branch_id, action_log_id) " +
		"values($1, $2);"
)

func CreateCodebaseBranchAction(txn sql.Tx, codebaseId int, codebaseActionId int) error {
	stmt, err := txn.Prepare(InsertCodebaseBranchActionLog)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(codebaseId, codebaseActionId)
	if err != nil {
		return err
	}
	return nil
}

func CreateCodebaseActionLog(txn sql.Tx, actionLog model.ActionLog) (*int, error) {
	stmt, err := txn.Prepare(InsertCodebaseActionLog)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(actionLog.Event, "", actionLog.Username, time.Unix(actionLog.UpdatedAt, 0).Format("2006-01-02 15:04:05")).Scan(&id)

	return &id, err
}

func GetLastIdCodebaseBranchActionLog(txn sql.Tx, codebaseBranch model.CodebaseBranch) (*int, error) {
	stmt, err := txn.Prepare(CheckDuplicateCodebaseBranchActionLog)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(codebaseBranch.Name, codebaseBranch.ActionLog.Event, time.Unix(codebaseBranch.ActionLog.UpdatedAt, 0).Format("2006-01-02 15:04:05")).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}