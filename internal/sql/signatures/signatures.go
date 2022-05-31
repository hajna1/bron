package signatures

import "github.com/hajna1/bron/internal/sql/models"

var (
	CreateTableSignature = []models.Token{
		{
			Identifier: models.Keyword,
			Value:      "create",
		},
		{
			Identifier: models.Keyword,
			Value:      "table",
		},
	}
	AlterTableSignature = []models.Token{
		{
			Identifier: models.Keyword,
			Value:      "alter",
		},
		{
			Identifier: models.Keyword,
			Value:      "table",
		},
	}
)

var (
	PrimaryKeySignature = []models.Token{
		{
			Identifier: models.Keyword,
			Value:      "primary",
		},
		{
			Identifier: models.Keyword,
			Value:      "key",
		},
	}

	NotNullToken = []models.Token{
		{
			Identifier: models.Keyword,
			Value:      "not",
		},
		{
			Identifier: models.Keyword,
			Value:      "null",
		},
	}

	AddColumnSignature = []models.Token{
		{
			Identifier: models.Keyword,
			Value:      "add",
		},
		{
			Identifier: models.Keyword,
			Value:      "column",
		},
	}

	DropColumnSignature = []models.Token{
		{
			Identifier: models.Keyword,
			Value:      "drop",
		},
		{
			Identifier: models.Keyword,
			Value:      "column",
		},
	}
)

var (
	OpenBracket = models.Token{
		Identifier: models.Separator,
		Value:      "(",
	}

	CloseBracket = models.Token{
		Identifier: models.Separator,
		Value:      ")",
	}

	SemiColumn = models.Token{
		Identifier: models.Separator,
		Value:      ";",
	}

	Constraint = models.Token{
		Identifier: models.Keyword,
		Value:      "constraint",
	}

	Comma = models.Token{
		Identifier: models.Separator,
		Value:      ",",
	}
)
