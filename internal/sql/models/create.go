package models

type Create struct {
	TableName    string
	CreateColumn []CreateColumn
}

type UpdateType int

const (
	CreateType UpdateType = iota
	RemoveType
)

type Update struct {
	UpdateType   UpdateType
	TableName    string
	CreateColumn CreateColumn
}
type CreateColumn struct {
	ColumnName     string
	ColumnType     string
	ColumnTypeSize int
	IsPrimaryKey   bool
	IsNotNull      bool
}

type TokenType int

const (
	Identifier TokenType = iota
	Separator
	Keyword
	Type
)

type Token struct {
	Identifier TokenType
	Value      string
}
