package fn

import "github.com/networkteam/qrb/builder"

func JsonBuildObject() builder.JsonBuildObjectBuilder {
	return builder.JsonBuildObject(false)
}

func JsonbBuildObject() builder.JsonBuildObjectBuilder {
	return builder.JsonBuildObject(true)
}
