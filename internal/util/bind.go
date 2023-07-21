package util

import (
	"fmt"
	"strconv"

	"github.com/96Asch/mkvstage-server/internal/domain"
	"github.com/gin-gonic/gin"
)

func BindNamedParams(ctx *gin.Context, names ...string) (map[string]int64, error) {
	idMap := make(map[string]int64, len(names))

	for _, name := range names {
		field := ctx.Params.ByName(name)

		fieldID, err := strconv.Atoi(field)
		if err != nil {
			return map[string]int64{}, domain.NewBadRequestErr(fmt.Sprintf("Could not read %s", field))
		}

		idMap[name] = int64(fieldID)
	}

	return idMap, nil
}
