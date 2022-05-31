package generator

type Repository struct {
	PackageName    string
	ConstName      string
	ModelDirectory string
	ModelType      string
	OriginName     string
	GoName         string
	ID             Column
	CreatedAt      *Column
	UpdatedAt      *Column

	Columns []Column
}

type Params struct {
	OriginName string
	GoName     string
	Columns    []Column
}

type Column struct {
	TableName  string
	OriginName string
	ConstName  string
	GoName     string
	GoType     string
	IsNullable bool
}
