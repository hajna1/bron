package lexer

import (
	"bufio"
	"fmt"
	"github.com/hajna1/bron/internal/sql/models"
	"github.com/hajna1/bron/internal/sql/query"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	sepatators = []byte{'(', ')', ',', ';'}

	backspaces = []byte{' ', '\n', '\t', '\r', '"'}

	keywords = []string{
		"create", "table", "alter", "index", "not", "null", "add", "column",
		"default", "primary", "key", "on", "unique", "true", "false",
		"foreign", "references", "drop", "exists", "constraint"}

	types = []string{"bigserial", "varchar", "timestamp", "real", "boolean"}
)

type Reader interface {
	ReadByte() (byte, error)
}

func New() *Lexer {
	mSep := make(map[byte]struct{})
	for _, sep := range sepatators {
		mSep[sep] = struct{}{}
	}

	mBackspaces := make(map[byte]struct{})
	for _, bs := range backspaces {
		mBackspaces[bs] = struct{}{}
	}
	mKeywords := make(map[string]struct{})
	for _, kw := range keywords {
		mKeywords[kw] = struct{}{}
	}

	mTypes := make(map[string]struct{})
	for _, tp := range types {
		mTypes[tp] = struct{}{}
	}

	l := Lexer{
		mapSeparators: mSep,
		mapBackspaces: mBackspaces,
		mapKeywords:   mKeywords,
		mapTypes:      mTypes,
	}
	return &l
}

type Lexer struct {
	mapSeparators map[byte]struct{}
	mapBackspaces map[byte]struct{}
	mapKeywords   map[string]struct{}
	mapTypes      map[string]struct{}
}

func (l *Lexer) ParseAllDirectory(directory string) ([]*query.Query, error) {
	result := make([]*query.Query, 0)
	fi, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, info := range fi {
		if info.IsDir() {
			continue
		}
		if strings.HasSuffix(info.Name(), "down.sql") {
			continue
		}
		f, err := os.Open(fmt.Sprintf("%s/%s", directory, info.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		queries, err := l.ReadAllQueries(reader)
		if err != nil {
			return nil, err
		}
		result = append(result, queries...)
	}
	return result, nil
}

func (l *Lexer) ReadAllQueries(r Reader) ([]*query.Query, error) {
	result := make([]*query.Query, 0)
	for {
		q, err := l.ReadLexems(r)
		if err != nil {
			if err == io.EOF {
				if len(q) != 0 {
					result = append(result, query.New(q))
				}
				return result, nil
			}
			return nil, err
		}
		result = append(result, query.New(q))
	}
}

func (l *Lexer) ReadLexems(r Reader) ([]models.Token, error) {
	buffer := make([]byte, 0, 1024)
	tokens := make([]models.Token, 0)
	temp := make([]byte, 2)
	b, err := r.ReadByte()
	if err != nil {
		return tokens, err
	}
	temp[0] = b
	for {
		b, err := r.ReadByte()
		if err != nil {
			if len(buffer) != 0 {
				tokens = append(tokens, l.handleWord(buffer))
			}
			return tokens, err
		}
		temp[1] = b

		switch {
		case l.IsBackspace(temp[0]):
			if len(buffer) != 0 {
				tokens = append(tokens, l.handleWord(buffer))
				buffer = make([]byte, 0, 1024)
			}
		case l.IsSeparator(temp[0]):
			if len(buffer) != 0 {
				tokens = append(tokens, l.handleWord(buffer))
				buffer = make([]byte, 0, 1024)
			}
			tokens = append(tokens, models.Token{
				Identifier: models.Separator,
				Value:      string([]byte{temp[0]}),
			})
			if temp[0] == ';' {
				return tokens, nil
			}
		case l.isSingleCommentary(string(temp)):
			if err := l.ignoreSingleCommentary(r); err != nil {
				return tokens, err
			}
			b, err := r.ReadByte()
			if err != nil {
				if len(buffer) != 0 {
					tokens = append(tokens, l.handleWord(buffer))
				}
				return tokens, err
			}
			temp[0] = b
			continue
		case l.isMultilineCommentary(string(temp)):
			if err := l.ignoreMultipleCommentary(r); err != nil {
				return tokens, err
			}
			b, err := r.ReadByte()
			if err != nil {
				if len(buffer) != 0 {
					tokens = append(tokens, l.handleWord(buffer))
				}
				return tokens, err
			}
			temp[0] = b
			continue
		default:
			buffer = append(buffer, temp[0])
		}
		temp[0] = temp[1]
	}
}

func (l *Lexer) handleWord(word []byte) models.Token {
	w := strings.ToLower(string(word))
	switch {
	case l.IsKeyword(w):
		return models.Token{
			Identifier: models.Keyword,
			Value:      w,
		}
	case l.IsType(w):
		return models.Token{
			Identifier: models.Type,
			Value:      w,
		}
	default:
		return models.Token{
			Identifier: models.Identifier,
			Value:      w,
		}
	}
}

func (l *Lexer) ignoreSingleCommentary(r Reader) error {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		if b == '\n' {
			return nil
		}
	}
}

func (l *Lexer) ignoreMultipleCommentary(r Reader) error {
	prev, err := r.ReadByte()
	if err != nil {
		return err
	}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return err
		}
		if prev == '*' && b == '/' {
			return nil
		}
		prev = b
	}
}

func (l *Lexer) isSingleCommentary(b string) bool {
	return b == "--"
}

func (l *Lexer) isMultilineCommentary(b string) bool {
	return b == "/*"

}

func (l *Lexer) IsBackspace(b byte) bool {
	_, has := l.mapBackspaces[b]
	return has
}

func (l *Lexer) IsSeparator(b byte) bool {
	_, has := l.mapSeparators[b]
	return has
}

func (l *Lexer) IsKeyword(s string) bool {
	_, has := l.mapKeywords[s]
	return has
}

func (l *Lexer) IsType(s string) bool {
	_, has := l.mapTypes[s]
	return has
}
