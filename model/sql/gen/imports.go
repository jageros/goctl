package gen

import (
	"github.com/jageros/goctl/model/sql/template"
	"github.com/jageros/goctl/util"
	"github.com/jageros/goctl/util/pathx"
)

func genImports(table Table, withCache, timeImport, decimalImport bool) (string, error) {
	if withCache {
		text, err := pathx.LoadTemplate(category, importsTemplateFile, template.Imports)
		if err != nil {
			return "", err
		}

		buffer, err := util.With("import").Parse(text).Execute(map[string]any{
			"time":       timeImport,
			"containsPQ": table.ContainsPQ,
			"data":       table,
			"decimal":    decimalImport,
		})
		if err != nil {
			return "", err
		}

		return buffer.String(), nil
	}

	text, err := pathx.LoadTemplate(category, importsWithNoCacheTemplateFile, template.ImportsNoCache)
	if err != nil {
		return "", err
	}

	buffer, err := util.With("import").Parse(text).Execute(map[string]any{
		"time":       timeImport,
		"containsPQ": table.ContainsPQ,
		"data":       table,
		"decimal":    decimalImport,
	})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
