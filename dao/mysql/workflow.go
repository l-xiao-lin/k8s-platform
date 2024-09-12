package mysql

import (
	"database/sql"
	"k8s-platform/model"
)

var Workflow workflow

type workflow struct{}

// CreateWorkflow 创建workflow
func (w *workflow) CreateWorkflow(workflow *model.Workflow) (err error) {
	sqlStr := `insert into workflow(name,namespace,replicas,deployment,service,ingress,type) values(?,?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, workflow.Name, workflow.Namespace, workflow.Replicas, workflow.Deployment, workflow.Service, workflow.Ingress, workflow.Type)
	return
}

func (w *workflow) GetWorkflowById(id int) (workflow *model.Workflow, err error) {
	workflow = new(model.Workflow)
	sqlStr := `select id,name,namespace,replicas,deployment,service,ingress,type,create_time,update_time from workflow where id =?`
	err = db.Get(workflow, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidId
			return
		}
	}
	return
}

func (w *workflow) DeleteWorkflow(id int) (err error) {
	sqlStr := `delete from workflow where id = ?`
	_, err = db.Exec(sqlStr, id)
	return

}
