package fixtures

import (
	. "github.com/dai/go-ipfs/gxlibs/github.com/polydawn/refmt/tok"
)

var sequences_String = []Sequence{
	{"empty string",
		[]Token{
			TokStr(""),
		},
	},
	{"flat string",
		[]Token{
			TokStr("value"),
		},
	},
	{"strings needing escape",
		[]Token{
			TokStr("str\nbroken\ttabbed"),
		},
	},
}
