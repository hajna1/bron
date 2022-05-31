package merger

import (
	"fmt"
	"github.com/hajna1/bron/internal/sql/models"
	"github.com/hajna1/bron/internal/sql/query"
	"github.com/hajna1/bron/internal/sql/signatures"
)

func ParseAlterTable(query *query.Query) (models.Update, error) {
	if err := query.TrimPrefix(signatures.AlterTableSignature); err != nil {
		return models.Update{}, err
	}
	tableName, err := query.Next()
	if err != nil {
		return models.Update{}, err
	}
	if tableName.Identifier != models.Identifier {
		return models.Update{}, fmt.Errorf("bad table name identifier")
	}
	result := models.Update{
		TableName: tableName.Value,
	}
	switch {
	case query.HasPrefix(signatures.AddColumnSignature):
		_ = query.TrimPrefix(signatures.AddColumnSignature)
		result.UpdateType = models.CreateType
		column, err := parseCreateColumn(query)
		if err != nil {
			return models.Update{}, err
		}
		result.CreateColumn = column
		return result, nil
	case query.HasPrefix(signatures.DropColumnSignature):
		_ = query.TrimPrefix(signatures.DropColumnSignature)
		result.UpdateType = models.RemoveType
		columnName, err := query.Next()
		if err != nil {
			return models.Update{}, err
		}
		if columnName.Identifier != models.Identifier {
			return models.Update{}, fmt.Errorf("bad column name identifier")
		}
		result.CreateColumn = models.CreateColumn{
			ColumnName: columnName.Value,
		}
		return result, nil
	default:
		return models.Update{}, fmt.Errorf("unknown alter table type")
	}
}
