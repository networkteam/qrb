package fn

import "github.com/networkteam/qrb/builder"

func Lower(identer builder.Identer) builder.LetterCaseBuilder {
	return builder.LetterCase("lower", identer)
}

func Upper(identer builder.Identer) builder.LetterCaseBuilder {
	return builder.LetterCase("upper", identer)
}

func InitCap(identer builder.Identer) builder.LetterCaseBuilder {
	return builder.LetterCase("initcap", identer)
}
