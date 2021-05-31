package repository

import (
	"database/sql"
	"fmt"
	"github.com/epam/edp-reconciler/v2/pkg/model"
	"github.com/epam/edp-reconciler/v2/pkg/model/cdpipeline"
)

const (
	insertCDPipeline             = "insert into \"%v\".cd_pipeline(name, deployment_type, status) VALUES ($1, $2, $3) returning id, name, deployment_type, status;"
	selectCDPipeline             = "select * from \"%v\".cd_pipeline cdp where cdp.name = $1 ;"
	updateCDPipelineStatusQuery  = "update \"%v\".cd_pipeline set status = $1 where id = $2 ;"
	insertCDPipelineDockerStream = "insert into \"%v\".cd_pipeline_docker_stream(cd_pipeline_id, codebase_docker_stream_id) VALUES ($1, $2);"
	deleteAllDockerStreams       = "delete from \"%v\".cd_pipeline_docker_stream cpds  where cpds.cd_pipeline_id = $1 ;"
	deleteCDPipeline             = "delete from \"%v\".cd_pipeline where name = $1 ;"
)

func CreateCDPipeline(txn sql.Tx, cdPipeline cdpipeline.CDPipeline, status, schema string) (*model.CDPipelineDTO, error) {
	stmt, err := txn.Prepare(fmt.Sprintf(insertCDPipeline, schema))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var cdPipelineDto model.CDPipelineDTO
	err = stmt.QueryRow(cdPipeline.Name, cdPipeline.DeploymentType, status).
		Scan(&cdPipelineDto.Id, &cdPipelineDto.Name, &cdPipelineDto.DeploymentType, &cdPipelineDto.Status)
	if err != nil {
		return nil, err
	}
	return &cdPipelineDto, nil
}

func GetCDPipeline(txn sql.Tx, cdPipelineName string, schemaName string) (*model.CDPipelineDTO, error) {
	stmt, err := txn.Prepare(fmt.Sprintf(selectCDPipeline, schemaName))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var cdPipeline model.CDPipelineDTO
	err = stmt.QueryRow(cdPipelineName).
		Scan(&cdPipeline.Id, &cdPipeline.Name, &cdPipeline.DeploymentType, &cdPipeline.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &cdPipeline, nil
}

func UpdateCDPipelineStatus(txn sql.Tx, pipelineId int, cdPipelineStatus string, schemaName string) error {
	stmt, err := txn.Prepare(fmt.Sprintf(updateCDPipelineStatusQuery, schemaName))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(cdPipelineStatus, pipelineId)
	return err
}

func CreateCDPipelineDockerStream(txn sql.Tx, pipelineId int, dockerStreamId int, schemaName string) error {
	stmt, err := txn.Prepare(fmt.Sprintf(insertCDPipelineDockerStream, schemaName))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pipelineId, dockerStreamId)
	return err
}

func DeleteCDPipelineDockerStreams(txn sql.Tx, pipelineId int, schemaName string) error {
	stmt, err := txn.Prepare(fmt.Sprintf(deleteAllDockerStreams, schemaName))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(pipelineId)
	return err
}

func DeleteCDPipeline(txn sql.Tx, pipeName, schema string) error {
	if _, err := txn.Exec(fmt.Sprintf(deleteCDPipeline, schema), pipeName); err != nil {
		return err
	}
	return nil
}
