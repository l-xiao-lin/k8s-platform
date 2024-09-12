package model

type Workflow struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Namespace  string `db:"namespace"`
	Replicas   int32  `db:"replicas"`
	Deployment string `db:"deployment"`
	Service    string `db:"service"`
	Ingress    string `db:"ingress"`
	Type       string `db:"type"`
	CreateTime string `db:"create_time"`
	UpdateTime string `db:"update_time"`
}
