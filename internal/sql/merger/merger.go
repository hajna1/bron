package merger

import (
	"github.com/hajna1/bron/internal/sql/models"
	"github.com/hajna1/bron/internal/sql/query"
	"github.com/hajna1/bron/internal/sql/signatures"
	"github.com/sirupsen/logrus"
)

func New() *Merger {
	return &Merger{}
}

type Merger struct {
}

func (m *Merger) IsCreateTable(query *query.Query) bool {
	return query.HasPrefix(signatures.CreateTableSignature)
}

func (m *Merger) IsAlterTable(query *query.Query) bool {
	return query.HasPrefix(signatures.AlterTableSignature)
}

func (m *Merger) ParseAll(qs ...*query.Query) map[string]models.Create {
	result := make(map[string]models.Create)
	for _, q := range qs {
		switch {
		case m.IsCreateTable(q):
			t, err := ParseCreateTable(q)
			if err != nil {
				logrus.Warnf("error: %s\n", err)
				continue
			}
			result[t.TableName] = t
		case m.IsAlterTable(q):
			command, err := ParseAlterTable(q)
			if err != nil {
				logrus.Warnf("error: %s\n", err)
				continue
			}
			switch command.UpdateType {
			case models.CreateType:
				t, has := result[command.TableName]
				if !has {
					logrus.Warnf("table name not found: %s", command.TableName)
					continue
				}
				t.CreateColumn = append(t.CreateColumn, command.CreateColumn)
				result[t.TableName] = t
			case models.RemoveType:
				t, has := result[command.TableName]
				if !has {
					logrus.Warnf("table name not found: %s", command.TableName)
					continue
				}
				columns := make([]models.CreateColumn, 0, len(t.CreateColumn))
				for _, c := range t.CreateColumn {
					if c.ColumnName != command.CreateColumn.ColumnName {
						columns = append(columns, c)
					}
				}
				t.CreateColumn = columns
				result[t.TableName] = t
			}
		default:
			logrus.Info("unknown query")
		}
	}
	return result
}
