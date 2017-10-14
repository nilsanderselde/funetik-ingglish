package main

import (
	"net/http"
)

var (
	defTilde = []string{"`", "~", "", ""}
	def1     = []string{"1", "!", "", ""}
	def2     = []string{"2", "@", "", ""}
	def3     = []string{"3", "#", "", ""}
	def4     = []string{"4", "$", "", ""}
	def5     = []string{"5", "%", "", ""}
	def6     = []string{"6", "^", "", ""}
	def7     = []string{"7", "&", "", ""}
	def8     = []string{"8", "*", "", ""}
	def9     = []string{"9", "(", "", ""}
	def0     = []string{"0", ")", "", ""}
	defDash  = []string{"-", "_", "", ""}
	defEq    = []string{"=", "+", "", ""}
	defQ     = []string{"ä", "Ä", "", ""}
	defW     = []string{"w", "W", "", wIPA}
	defE     = []string{"e", "E", "", ɛIPA}
	defR     = []string{"r", "R", "", rIPA}
	defT     = []string{"t", "T", "", tIPA}
	defY     = []string{"y", "Y", "", jIPA}
	defU     = []string{"u", "U", "", ʌIPA}
	defI     = []string{"i", "I", "", ɪIPA}
	defO     = []string{"o", "O", "", oIPA}
	defP     = []string{"p", "P", "", pIPA}
	defLb    = []string{"[", "{", "", ""}
	defRb    = []string{"]", "}", "", ""}
	defBksl  = []string{"\\", "|", "", ""}
	defA     = []string{"a", "A", "", ɑIPA}
	defS     = []string{"s", "S", "", sIPA}
	defD     = []string{"d", "D", "", dIPA}
	defF     = []string{"f", "F", "", fIPA}
	defG     = []string{"g", "G", "", ɡIPA}
	defH     = []string{"h", "H", "", hIPA}
	defJ     = []string{"ø", "Ø", "", dʒIPA}
	defK     = []string{"k", "K", "", kIPA}
	defL     = []string{"l", "L", "", lIPA}
	defSemi  = []string{";", ":", "", ""}
	defQt    = []string{"'", "\"", "", ""}
	defZ     = []string{"z", "Z", "", zIPA}
	defX     = []string{"x", "X", "", ""}
	defC     = []string{"c", "C", "", tʃIPA}
	defV     = []string{"v", "V", "", vIPA}
	defB     = []string{"b", "B", "", bIPA}
	defN     = []string{"n", "N", "", nIPA}
	defM     = []string{"m", "M", "", mIPA}
	defCom   = []string{",", "<", "", ""}
	defPer   = []string{".", ">", "", ""}
	defSl    = []string{"/", "?", "", ""}
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
				{"ä", "Ä", "final", æIPA},
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
				{"ʃ", "Ʃ", "old", ʃIPA},
				defK,
				{"л", "Л", "old", lIPA},
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ʒ", "Ʒ", "old", ʒIPA},
				{"ï", "Ï", "old", aɪIPA},
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
				{"æ", "Æ", "old", æIPA},
				defW,
				defE,
				defR,
				defT,
				defY,
				defU,
				defI,
				defO,
				defP,
				{"'", "\"", "old", ""},
				{"?", "[", "old", ""},
				{"/", "]", "old", ""},
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"ʃ", "Ʃ", "old", ʃIPA},
				defK,
				{"λ", "Λ", "old", lIPA},
				{"ā", "Ā", "old", eɪIPA},
				{".", ":", "old", ""},
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ē", "Ē", "old", iIPA},
				{"ī", "Ī", "old", aɪIPA},
				defV,
				defB,
				defN,
				defM,
				{"ū", "Ū", "old", uIPA},
				{"ʒ", "Ʒ", "old", ʒIPA},
				{",", ";", "old", ""},
			}
			return

		case "3":
			t.args.Kbd[0] = [][]string{
				{"ð", "Ð", "old", ðIPA},
				def1,
				{"2", "/", "old", ""},
				def3,
				def4,
				def5,
				{"6", ":", "old", ""},
				{"7", "?", "old", ""},
				def8,
				def9,
				def0,
				defDash,
				{"'", "\"", "old", ""},
			}
			t.args.Kbd[1] = [][]string{
				{"æ", "Æ", "old", æIPA},
				defW,
				defE,
				defR,
				defT,
				defY,
				defU,
				defI,
				defO,
				defP,
				{"ч", "Ч", "old", tʃIPA},
				{"θ", "Θ", "old", θIPA},
				{"=", "+", "old", ""},
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				{"j", "J", "", dʒIPA},
				defK,
				{"л", "Л", "old", lIPA},
				{"ā", "Ā", "old", eɪIPA},
				{"ʒ", "Ʒ", "old", ʒIPA},
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ē", "Ē", "old", iIPA},
				{"ī", "Ī", "old", aɪIPA},
				defV,
				defB,
				defN,
				defM,
				{"ū", "Ū", "old", uIPA},
				{"ʃ", "Ʃ", "old", ʃIPA},
				{".", ",", "old", ""},
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
				{"ä", "Ä", "final", æIPA},
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
				{"j", "J", "", dʒIPA},
				defK,
				defL,
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ĭ", "Ĭ", "old", ɪIPA},
				{"c", "C", "", tʃIPA},
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
				{"æ", "Æ", "old", æIPA},
				{"ū", "Ū", "old", uwIPA},
				defE,
				defR,
				defT,
				{"ē", "Ē", "old", ijIPA},
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
				{"ā", "Ā", "old", eɪIPA},
				defK,
				{"л", "Л", "old", lIPA},
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ʊ", "Ʊ", "old", ʊIPA},
				{"ī", "Ī", "old", aɪIPA},
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
				{"ō", "Ō", "old", oIPA},
				{"ū", "Ū", "old", uwIPA},
				defE,
				defR,
				defT,
				{"ē", "Ē", "old", ijIPA},
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
				{"ā", "Ā", "old", eɪIPA},
				defK,
				{"л", "Л", "old", lIPA},
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"ø", "Ø", "old", ʊIPA},
				{"ī", "Ī", "old", aɪIPA},
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
				{"ž", "Ž", "old", ʒIPA},
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
				{"ó", "Ó", "old", oIPA},
				{"ú", "Ú", "old", uwIPA},
				defE,
				defR,
				defT,
				{"é", "É", "old", ijIPA},
				defU,
				defI,
				defO,
				defP,
				{"ø", "Ø", "old", ʊIPA},
				{"š", "Š", "old", ʃIPA},
				defBksl,
			}
			t.args.Kbd[2] = [][]string{
				defA,
				defS,
				defD,
				defF,
				defG,
				defH,
				defJ,
				defK,
				defL,
				{"á", "Á", "old", eɪIPA},
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"í", "Í", "old", aɪIPA},
				defC,
				defV,
				defB,
				defN,
				defM,
				{",", ";", "old", ""},
				{".", ":", "old", ""},
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
				{"ó", "Ó", "old", oIPA},
				{"ú", "Ú", "old", uwIPA},
				defE,
				defR,
				defT,
				{"é", "É", "old", ijIPA},
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
				{"ø", "Ø", "final", ʊIPA},
				defK,
				defL,
				defSemi,
				defQt,
			}
			t.args.Kbd[3] = [][]string{
				defZ,
				{"í", "Í", "old", aɪIPA},
				{"á", "Á", "old", eɪIPA},
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
		t.args.CurrentPage = "?v=9"
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
			{"ä", "Ä", "final", æIPA},
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
			{"ø", "Ø", "final", ʊIPA},
			defK,
			defL,
			defSemi,
			defQt,
		}
		t.args.Kbd[3] = [][]string{
			defZ,
			{"ž", "Ž", "final", ʒIPA},
			{"š", "Š", "final", ʃIPA},
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
