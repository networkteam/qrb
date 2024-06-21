package fn

import "github.com/networkteam/qrb/builder"

func Lower(identer builder.Exp) builder.ExpBase {
	return builder.FuncExp("lower", []builder.Exp{identer})
}

func Upper(identer builder.Exp) builder.ExpBase {
	return builder.FuncExp("upper", []builder.Exp{identer})
}

func Initcap(identer builder.Exp) builder.ExpBase {
	return builder.FuncExp("initcap", []builder.Exp{identer})
}
