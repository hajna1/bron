package merger

import (
	"fmt"
	"github.com/hajna1/bron/internal/sql/models"
	"github.com/hajna1/bron/internal/sql/query"
	"github.com/hajna1/bron/internal/sql/signatures"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ParseCreateTable(query *query.Query) (models.Create, error) {
	if err := query.TrimPrefix(signatures.CreateTableSignature); err != nil {
		return models.Create{}, err
	}
	nameToken, err := query.Next()
	if err != nil {
		return models.Create{}, err
	}
	if nameToken.Identifier != models.Identifier {
		return models.Create{}, fmt.Errorf("bad table name identifer")
	}
	table := models.Create{
		TableName:    nameToken.Value,
		CreateColumn: make([]models.CreateColumn, 0),
	}

	token, err := query.Next()
	if err != nil {
		return table, err
	}

	if token != signatures.OpenBracket {
		return table, fmt.Errorf("open bracket is missing")
	}

	for {
		token, err = query.Get()
		switch {
		case token == signatures.CloseBracket:
			return table, nil
		case token == signatures.Constraint:
			if err := parseCreateForeignKey(query); err != nil {
				return table, err
			}
		default:
			column, err := parseCreateColumn(query)
			if err != nil {
				return table, nil
			}
			table.CreateColumn = append(table.CreateColumn, column)
		}

	}
}

func parseCreateForeignKey(query *query.Query) error {
	//TODO: implement it
	for {
		token, err := query.Get()
		if err != nil {
			return err
		}
		if token == signatures.Comma {
			_, _ = query.Next()
			return nil
		}
		if token == signatures.CloseBracket {
			return nil
		}
		_, _ = query.Next()
	}
}

func parseCreateColumn(query *query.Query) (models.CreateColumn, error) {
	columnName, err := query.Next()
	if err != nil {
		return models.CreateColumn{}, err
	}

	if columnName.Identifier != models.Identifier {
		return models.CreateColumn{}, fmt.Errorf("bad column name")
	}
	result := models.CreateColumn{
		ColumnName: columnName.Value,
	}

	columnType, err := query.Next()
	if err != nil {
		return models.CreateColumn{}, err
	}
	if columnType.Identifier != models.Type {
		return models.CreateColumn{}, fmt.Errorf("bad column type")
	}
	result.ColumnType = columnType.Value

	if result.ColumnType == "varchar" {
		token, err := query.Next()
		if err != nil {
			return models.CreateColumn{}, err
		}
		if token == signatures.OpenBracket {
			varcharSize, err := query.Next()
			if err != nil {
				return models.CreateColumn{}, err
			}

			if varcharSize.Identifier != models.Identifier {
				return result, fmt.Errorf("varchar type has bad identifier")
			}
			v, err := strconv.Atoi(varcharSize.Value)
			if err != nil {
				return result, err
			}
			result.ColumnTypeSize = v

			token, err = query.Next()
			if err != nil {
				return models.CreateColumn{}, err
			}
			if token != signatures.CloseBracket {
				return result, fmt.Errorf("missing cloce bracket in varcher definition")
			}
		}
	}

	for {
		token, err := query.Get()
		if err != nil {
			return models.CreateColumn{}, err
		}
		if token == signatures.Comma {
			_, _ = query.Next()
			return result, nil
		}
		if token == signatures.CloseBracket {
			return result, nil
		}
		if token == signatures.SemiColumn {
			return result, nil
		}
		switch {
		case query.HasPrefix(signatures.PrimaryKeySignature):
			_ = query.TrimPrefix(signatures.PrimaryKeySignature)
			result.IsPrimaryKey = true
		case query.HasPrefix(signatures.NotNullToken):
			_ = query.TrimPrefix(signatures.NotNullToken)
			result.IsPrimaryKey = true
		case token == signatures.OpenBracket:
			if err := ignoreBrackets(query); err != nil {
				return models.CreateColumn{}, err
			}
		default:
			logrus.Debugf("bad or unimplemented parameter: %s\n", token.Value)
			_, _ = query.Next()
		}
	}
}

func ignoreBrackets(query *query.Query) error {
	_, _ = query.Next()
	numBrackets := 1
	for {
		token, err := query.Next()
		if err != nil {
			return err
		}
		switch token {
		case signatures.OpenBracket:
			numBrackets += 1
		case signatures.CloseBracket:
			numBrackets -= 1
			if numBrackets == 0 {
				return nil
			}
		}
	}
}
