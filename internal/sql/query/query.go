package query

import (
	"fmt"
	"github.com/hajna1/bron/internal/sql/models"
)

func New(tokens []models.Token) *Query {
	return &Query{
		Tokens: tokens,
		pos:    0,
	}
}

type Query struct {
	Tokens []models.Token
	pos    int
}

func (q *Query) Size() int {
	return len(q.Tokens) - q.pos
}

func (q *Query) Next() (models.Token, error) {
	if q.pos >= len(q.Tokens) {
		return models.Token{}, fmt.Errorf("end of query")
	}
	res := q.Tokens[q.pos]
	q.pos++
	return res, nil
}

func (q *Query) HasPrefix(prefix []models.Token) bool {
	if q.Size() < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if q.Tokens[q.pos+i] != prefix[i] {
			return false
		}
	}
	return true
}

func (q *Query) TrimPrefix(prefix []models.Token) error {
	if !q.HasPrefix(prefix) {
		return fmt.Errorf("query does not have such a prefix")
	}
	q.pos += len(prefix)
	return nil
}

func (q *Query) Get() (models.Token, error) {
	if q.pos >= len(q.Tokens) {
		return models.Token{}, fmt.Errorf("end of query")
	}
	res := q.Tokens[q.pos]
	return res, nil
}
