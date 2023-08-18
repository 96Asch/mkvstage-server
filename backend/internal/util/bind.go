package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/96Asch/mkvstage-server/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "is required"
	case "email":
		return "is an invalid email"
	}

	return ""
}

func BindModel(ctx *gin.Context, model interface{}) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	if err := ctx.BindJSON(model); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errstrings := make([]string, len(ve))
			for idx, fieldErr := range ve {
				errstrings[idx] = fmt.Sprintf(
					"field '%s' %s",
					fieldErr.Field(),
					msgForTag(fieldErr.Tag()))
			}

			return domain.NewBadRequestErr(strings.Join(errstrings, ","))
		}

		return domain.NewBadRequestErr(err.Error())
	}

	return nil
}
