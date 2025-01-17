// Package json minifies JSON following the specifications at http://json.org/.
package json

import (
	"io"

	"github.com/ezoic/minify/v2"
	"github.com/ezoic/parse/json"
)

var (
	commaBytes     = []byte(",")
	colonBytes     = []byte(":")
	zeroBytes      = []byte("0")
	minusZeroBytes = []byte("-0")
)

////////////////////////////////////////////////////////////////

// DefaultMinifier is the default minifier.
var DefaultMinifier = &Minifier{}

// Minifier is a JSON minifier.
type Minifier struct {
	Precision int // number of significant digits
}

// Minify minifies JSON data, it reads from r and writes to w.
func Minify(m *minify.M, w io.Writer, r io.Reader, params map[string]string) error {
	return DefaultMinifier.Minify(m, w, r, params)
}

// Minify minifies JSON data, it reads from r and writes to w.
func (o *Minifier) Minify(_ *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
	skipComma := true

	p := json.NewParser(r)
	defer p.Restore()

	for {
		state := p.State()
		gt, text := p.Next()
		if gt == json.ErrorGrammar {
			if p.Err() != io.EOF {
				return p.Err()
			}
			return nil
		}

		if !skipComma && gt != json.EndObjectGrammar && gt != json.EndArrayGrammar {
			if state == json.ObjectKeyState || state == json.ArrayState {
				if _, err := w.Write(commaBytes); err != nil {
					return err
				}
			} else if state == json.ObjectValueState {
				if _, err := w.Write(colonBytes); err != nil {
					return err
				}
			}
		}
		skipComma = gt == json.StartObjectGrammar || gt == json.StartArrayGrammar

		if 0 < len(text) && ('0' <= text[0] && text[0] <= '9' || text[0] == '-') {
			text = minify.Number(text, o.Precision)
			if text[0] == '.' {
				w.Write(zeroBytes)
			} else if 1 < len(text) && text[0] == '-' && text[1] == '.' {
				text = text[1:]
				w.Write(minusZeroBytes)
			}
		}
		if _, err := w.Write(text); err != nil {
			return err
		}
	}
}
