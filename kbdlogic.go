// Nils Elde
// https://github.com/nilsanderselde/funetik-ingglish

package main

import (
	"net/http"
)

var (
	defTilde = []string{"`", "~", "", "", ""}
	def1     = []string{"1", "!", "", "", ""}
	def2     = []string{"2", "@", "", "", ""}
	def3     = []string{"3", "#", "", "", ""}
	def4     = []string{"4", "$", "", "", ""}
	def5     = []string{"5", "%", "", "", ""}
	def6     = []string{"6", "^", "", "", ""}
	def7     = []string{"7", "&", "", "", ""}
	def8     = []string{"8", "*", "", "", ""}
	def9     = []string{"9", "(", "", "", ""}
	def0     = []string{"0", ")", "", "", ""}
	defDash  = []string{"-", "_", "", "", ""}
	defEq    = []string{"=", "+", "", "", ""}
	defQ     = []string{"ä", "Ä", "", "", ""}
	defW     = []string{"w", "W", "", wIPA, "w"}
	defE     = []string{"e", "E", "", ɛIPA, "e"}
	defR     = []string{"r", "R", "", rIPA, "r"}
	defT     = []string{"t", "T", "", tIPA, "t"}
	defY     = []string{"y", "Y", "", jIPA, "y"}
	defU     = []string{"u", "U", "", ʌIPA, "u"}
	defI     = []string{"i", "I", "", ɪIPA, "i"}
	defO     = []string{"o", "O", "", oIPA, "o"}
	defP     = []string{"p", "P", "", pIPA, "p"}
	defLb    = []string{"[", "{", "", "", ""}
	defRb    = []string{"]", "}", "", "", ""}
	defBksl  = []string{"\\", "|", "", "", ""}
	defA     = []string{"a", "A", "", ɑIPA, "a"}
	defS     = []string{"s", "S", "", sIPA, "s"}
	defD     = []string{"d", "D", "", dIPA, "d"}
	defF     = []string{"f", "F", "", fIPA, "f"}
	defG     = []string{"g", "G", "", ɡIPA, "g"}
	defH     = []string{"h", "H", "", hIPA, "h"}
	defJ     = []string{"j", "J", "", dʒIPA, "dzh"}
	defK     = []string{"k", "K", "", kIPA, "k"}
	defL     = []string{"l", "L", "", lIPA, "l"}
	defSemi  = []string{";", ":", "", "", ""}
	defQt    = []string{"'", "\"", "", "", ""}
	defZ     = []string{"z", "Z", "", zIPA, "z"}
	defX     = []string{"x", "X", "", "", ""}
	defC     = []string{"c", "C", "", tʃIPA, "tsh"}
	defV     = []string{"v", "V", "", vIPA, "v"}
	defB     = []string{"b", "B", "", bIPA, "b"}
	defN     = []string{"n", "N", "", nIPA, "n"}
	defM     = []string{"m", "M", "", mIPA, "m"}
	defCom   = []string{",", "<", "", "", ""}
	defPer   = []string{".", ">", "", "", ""}
	defSl    = []string{"/", "?", "", "", ""}
)

var (
	ɑIPA  = "IPA: ɑ/ɔ/ɒ"
	æIPA  = "IPA: æ"
	ɛIPA  = "IPA: ɛ/e"
	ɪIPA  = "IPA: ɪ"
	iIPA  = "IPA: i"
	jIPA  = "IPA: j"
	ijIPA = "IPA: i/j"
	uIPA  = "IPA: u"
	wIPA  = "IPA: w"
	uwIPA = "IPA: u/w"
	ʌIPA  = "IPA: ʌ/ɐ"
	ʊIPA  = "IPA: ʊ"
	oIPA  = "IPA: oʊ/əʊ/əʉ"
	aɪIPA = "IPA: aɪ/ɑe"
	aʊIPA = "IPA: aʊ/æɔ"
	eɪIPA = "IPA: eɪ/æɪ"
	ɔɪIPA = "IPA: ɔɪ/oɪ"
	rIPA  = "IPA: ɹ̠"
	lIPA  = "IPA: l"
	nIPA  = "IPA: n"
	mIPA  = "IPA: m"
	bIPA  = "IPA: b"
	pIPA  = "IPA: p"
	vIPA  = "IPA: v"
	fIPA  = "IPA: f"
	ɡIPA  = "IPA: ɡ"
	kIPA  = "IPA: k"
	dIPA  = "IPA: d"
	tIPA  = "IPA: t"
	zIPA  = "IPA: z"
	sIPA  = "IPA: s"
	ʒIPA  = "IPA: ʒ"
	ʃIPA  = "IPA: ʃ"
	hIPA  = "IPA: h"
	ŋIPA  = "IPA: ŋ"
	ðIPA  = "IPA: ð"
	θIPA  = "IPA: θ"
	dʒIPA = "IPA: d͡ʒ"
	tʃIPA = "IPA: t͡ʃ"
)

// pickKeyboard determines which keyboard layout to show on keyboard page
func pickKeyboard(t *templateHandler, r *http.Request) {
	t.args.Kbd = make([][][]string, 4, 4)
	if r.URL.Query()["v"] != nil {
		v := &r.URL.Query()["v"][0]

		t.args.CurrentPage = "?v=" + *v
		t.args.KbdVer = *v

		switch *v {
		// "9" is default, to be executed regardless of the
		// existence of a version query, so if "9" is passed, skip to the
		// key data defined outside the query existence test
		case "9":
			break
		case "1":
			t.args.Kbd[0] = [][]string{
				defTilde,
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"ä", "Ä", "final", æIPA, "aa"},
				defW,
				defE,
				defR,
				defT,
				defY,
				defU,
				defI,
				defO,
				defP,
				defLb,
				defRb,
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"ʃ", "Ʃ", "old", ʃIPA, "sh"},
				defK,
				{"л", "Л", "old", lIPA, "l"},
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ʒ", "Ʒ", "old", ʒIPA, "zh"},
				{"ï", "Ï", "old", aɪIPA, "ai"},
				defV,
				defB,
				defN,
				defM,
				defCom,
				defPer,
				defSl,
			}
			return
		case "2":
			t.args.Kbd[0] = [][]string{
				defTilde,
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"æ", "Æ", "old", æIPA, "aa"},
				defW,
				defE,
				defR,
				defT,
				defY,
				defU,
				defI,
				defO,
				defP,
				{"'", "\"", "old", "", ""},
				{"?", "[", "old", "", ""},
				{"/", "]", "old", "", ""},
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"ʃ", "Ʃ", "old", ʃIPA, "sh"},
				defK,
				{"λ", "Λ", "old", lIPA, "l"},
				{"ā", "Ā", "old", eɪIPA, "ei"},
				{".", ":", "old", "", ""},
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ē", "Ē", "old", iIPA, "y"},
				{"ī", "Ī", "old", aɪIPA, "ai"},
				defV,
				defB,
				defN,
				defM,
				{"ū", "Ū", "old", uIPA, "w"},
				{"ʒ", "Ʒ", "old", ʒIPA, "zh"},
				{",", ";", "old", "", ""},
			}
			return

		case "3":
			t.args.Kbd[0] = [][]string{
				{"ð", "Ð", "old", ðIPA, "dh"},
				def1,
				{"2", "/", "old", "", ""},
				def3,
				def4,
				def5,
				{"6", ":", "old", "", ""},
				{"7", "?", "old", "", ""},
				def8,
				def9,
				def0,
				defDash,
				{"'", "\"", "old", "", ""},
			}
			t.args.Kbd[1] = [][]string{
				{"æ", "Æ", "old", æIPA, "aa"},
				defW,
				defE,
				defR,
				defT,
				defY,
				defU,
				defI,
				defO,
				defP,
				{"ч", "Ч", "old", tʃIPA, "tsh"},
				{"θ", "Θ", "old", θIPA, "th"},
				{"=", "+", "old", "", ""},
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"j", "J", "", dʒIPA, "dzh"},
				defK,
				{"л", "Л", "old", lIPA, "l"},
				{"ā", "Ā", "old", eɪIPA, "ei"},
				{"ʒ", "Ʒ", "old", ʒIPA, "zh"},
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ē", "Ē", "old", iIPA, "y"},
				{"ī", "Ī", "old", aɪIPA, "ai"},
				defV,
				defB,
				defN,
				defM,
				{"ū", "Ū", "old", uIPA, "w"},
				{"ʃ", "Ʃ", "old", ʃIPA, "sh"},
				{".", ",", "old", "", ""},
			}
			return
		case "4":
			t.args.Kbd[0] = [][]string{
				defTilde,
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"ä", "Ä", "final", æIPA, "aa"},
				defW,
				defE,
				defR,
				defT,
				defY,
				defU,
				{"i", "i", "", aɪIPA, "ai"},
				defO,
				defP,
				defLb,
				defRb,
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"j", "J", "", dʒIPA, "dzh"},
				defK,
				defL,
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ĭ", "Ĭ", "old", ɪIPA, "i"},
				{"c", "C", "", tʃIPA, "tsh"},
				defV,
				defB,
				defN,
				defM,
				defCom,
				defPer,
				defSl,
			}
			return
		case "5":
			t.args.Kbd[0] = [][]string{
				defTilde,
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"æ", "Æ", "old", æIPA, "aa"},
				{"ū", "Ū", "old", uwIPA, "w"},
				defE,
				defR,
				defT,
				{"ē", "Ē", "old", ijIPA, "y"},
				defU,
				defI,
				defO,
				defP,
				defLb,
				defRb,
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"ā", "Ā", "old", eɪIPA, "ei"},
				defK,
				{"л", "Л", "old", lIPA, "l"},
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ʊ", "Ʊ", "old", ʊIPA, "oo"},
				{"ī", "Ī", "old", aɪIPA, "ai"},
				defV,
				defB,
				defN,
				defM,
				defCom,
				defPer,
				defSl,
			}
			return
		case "6":
			t.args.Kbd[0] = [][]string{
				defTilde,
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"ō", "Ō", "old", oIPA, "o"},
				{"ū", "Ū", "old", uwIPA, "w"},
				defE,
				defR,
				defT,
				{"ē", "Ē", "old", ijIPA, "y"},
				defU,
				defI,
				{"o", "O", "", ɑIPA, "a"},
				defP,
				defLb,
				defRb,
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				{"a", "A", "", æIPA, "aa"},
				defS,
				defD,
				defF,
				defG,
				defH,
				{"ā", "Ā", "old", eɪIPA, "ei"},
				defK,
				{"л", "Л", "old", lIPA, "l"},
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ø", "Ø", "old", ʊIPA, "oo"},
				{"ī", "Ī", "old", aɪIPA, "ai"},
				defV,
				defB,
				defN,
				defM,
				defCom,
				defPer,
				defSl,
			}
			return
		case "7":
			t.args.Kbd[0] = [][]string{
				{"ž", "Ž", "old", ʒIPA, "zh"},
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"ó", "Ó", "old", oIPA, "o"},
				{"ú", "Ú", "old", uwIPA, "w"},
				defE,
				defR,
				defT,
				{"é", "É", "old", ijIPA, "y"},
				defU,
				defI,
				{"o", "O", "", ɑIPA, "a"},
				defP,
				{"ø", "Ø", "old", ʊIPA, "oo"},
				{"š", "Š", "old", ʃIPA, "sh"},
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				{"a", "A", "", æIPA, "aa"},
				defS,
				defD,
				defF,
				defG,
				defH,
				defJ,
				defK,
				defL,
				{"á", "Á", "old", eɪIPA, "ei"},
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"í", "Í", "old", aɪIPA, "ai"},
				defC,
				defV,
				defB,
				defN,
				defM,
				{",", ";", "old", "", ""},
				{".", ":", "old", "", ""},
				defSl,
			}
			return
		case "8":
			t.args.Kbd[0] = [][]string{
				defTilde,
				def1,
				def2,
				def3,
				def4,
				def5,
				def6,
				def7,
				def8,
				def9,
				def0,
				defDash,
				defEq,
			}
			t.args.Kbd[1] = [][]string{
				{"ó", "Ó", "old", oIPA, "o"},
				{"ú", "Ú", "old", uwIPA, "w"},
				defE,
				defR,
				defT,
				{"é", "É", "old", ijIPA, "y"},
				defU,
				defI,
				{"o", "O", "", ɑIPA, "a"},
				defP,
				defLb,
				defRb,
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				{"a", "A", "", æIPA, "aa"},
				defS,
				defD,
				defF,
				defG,
				defH,
				{"ø", "Ø", "final", ʊIPA, "oo"},
				defK,
				defL,
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"í", "Í", "old", aɪIPA, "ai"},
				{"á", "Á", "old", eɪIPA, "ei"},
				defV,
				defB,
				defN,
				defM,
				defCom,
				defPer,
				defSl,
			}
			return
		}
	}
	{
		// default (version 9) keyboard if no version query
		// exists or if version query is "9".
		// it's in its own block so it can be easily collapsed.
		t.args.CurrentPage = r.URL.Path
		t.args.KbdVer = "9"

		t.args.Kbd[0] = [][]string{
			defTilde,
			def1,
			def2,
			def3,
			def4,
			def5,
			def6,
			def7,
			def8,
			def9,
			def0,
			defDash,
			defEq,
		}
		t.args.Kbd[1] = [][]string{
			{"ä", "Ä", "final", æIPA, "aa"},
			{"w", "W", "", uwIPA, "w"},
			defE,
			defR,
			defT,
			{"y", "Y", "", ijIPA, "y"},
			defU,
			defI,
			defO,
			defP,
			defLb,
			defRb,
			defBksl,
		}
		t.args.Kbd[2] = [][]string{
			defA,
			defS,
			defD,
			defF,
			defG,
			defH,
			{"ø", "Ø", "final", ʊIPA, "oo"},
			defK,
			defL,
			defSemi,
			defQt,
		}
		t.args.Kbd[3] = [][]string{
			defZ,
			{"ž", "Ž", "final", ʒIPA, "zh"},
			{"š", "Š", "final", ʃIPA, "sh"},
			defV,
			defB,
			defN,
			defM,
			defCom,
			defPer,
			defSl,
		}
	}
}
