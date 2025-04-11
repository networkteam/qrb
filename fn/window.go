package fn

import "github.com/networkteam/qrb/builder"

func RowNumber() builder.WindowFuncCallBuilder {
	return builder.WindowFuncCallBuilder{
		FuncCall: builder.FuncExp("row_number", nil),
	}
}
