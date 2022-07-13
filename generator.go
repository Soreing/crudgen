package crudgen

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Executes the descriptions to create CRUD functions
func Run(dscs []CRUDDescription) {
	for _, dsc := range dscs {
		vals := buildValues(dsc)
		generate(vals, dsc.SrcFilePath)
	}
}

// Parse Tags to create a map of properties with values
func parseTag(tag string) map[string]string {
	props := map[string]string{}
	keyValRe, _ := regexp.Compile(`^([A-Za-z0-9]*)(='([A-Za-z0-9_ ]*)')?$`)

	toks := strings.Split(tag, ",")
	for _, tok := range toks {
		grps := keyValRe.FindStringSubmatch(tok)

		switch len(grps) {
		case 2:
			props[grps[1]] = ""
		case 4:
			props[grps[1]] = grps[3]
		default:
			panic(fmt.Errorf("unexpected token key value format in \"%s\"", tok))
		}
	}

	return props
}

func getOrder(val string, ok bool) int {
	if !ok {
		return 0
	} else if num, err := strconv.Atoi(val); err != nil {
		panic(fmt.Errorf("failed to convert %s to int: %w", val, err))
	} else {
		return num
	}
}

// Builds pipeline values for the template from the description (input)
// Inspects the fields and tags of the model type to extract data for constructing queries
func buildValues(dsc CRUDDescription) CRUDTemplateValues {
	vals := CRUDTemplateValues{
		Package:       dsc.PackageName,
		Imports:       map[string]string{},
		RepoTypeName:  dsc.RepoType.Name(),
		ModelTypeName: dsc.ModelType.String(),
		InterfaceName: dsc.ModelType.Name() + "CRUD",
		TableName:     dsc.TableName,
		PKeyFields:    []Attrib{},
		InsertFields:  []Attrib{},
		UpdateFields:  []Attrib{},
	}

	for i := 0; i < dsc.ModelType.NumField(); i++ {
		field := dsc.ModelType.Field(i)
		props := parseTag(field.Tag.Get("gen"))

		// Get column name that is the alias of the field in the db
		column, ex := props["name"]
		if !ex || column == "" {
			column = strings.ToLower(field.Name)
		}

		// Get order that tells the order of apearance in parameter lists and placeholders
		order := 0
		if val, ok := props["order"]; ok {
			if num, err := strconv.Atoi(val); err != nil {
				panic(fmt.Errorf("failed to convert %s to int: %w", val, err))
			} else {
				order = num
			}
		}

		atr := Attrib{
			Name:   field.Name,
			Type:   field.Type.String(),
			Column: column,
			Order:  order,
		}

		// If the field is marked as pk, add it to privatek eys
		if _, ok := props["pk"]; ok {
			vals.PKeyFields = append(vals.PKeyFields, atr)
			if field.Type.PkgPath() != "" {
				vals.Imports[field.Type.PkgPath()] = field.Type.PkgPath()
			}
		}
		// If the field is NOT omitted from update, add it to updates
		if _, ok := props["omitup"]; !ok {
			vals.UpdateFields = append(vals.UpdateFields, atr)
		}
		// If the field is NOT omitted from insert, add it to insert
		if _, ok := props["omitin"]; !ok {
			vals.InsertFields = append(vals.InsertFields, atr)
		}
	}

	// Check PK Orders for correctness and uniqueness
	sort.Sort(AttribList(vals.PKeyFields))
	for i, e := range vals.PKeyFields {
		if e.Order == 0 {
			panic(fmt.Errorf("%s is a primary key and must have an order > 0", e.Name))
		} else if i+1 != e.Order {
			panic(fmt.Errorf("order of %s is not unique (%d)", e.Name, e.Order))
		}
	}

	if len(vals.PKeyFields) == 0 {
		panic(fmt.Errorf("%s model must have a primary key", dsc.ModelType.String()))
	}

	// Assign order to update and insert attribute lists
	for i := 0; i < len(vals.InsertFields); i++ {
		vals.InsertFields[i].Order = i + 1
	}
	for i := 0; i < len(vals.UpdateFields); i++ {
		vals.UpdateFields[i].Order = i + 1 + len(vals.PKeyFields)
	}

	vals.Imports[dsc.ModelType.PkgPath()] = dsc.ModelType.PkgPath()

	return vals
}

// Generates the source file
func generate(vals CRUDTemplateValues, filename string) {
	buf := new(bytes.Buffer)
	err := pkgTemplate.Execute(buf, vals)

	if err != nil {
		panic(fmt.Errorf("failed to execute template: %w", err))
	} else {
		if f, err := os.Create(filename); err != nil {
			panic(fmt.Errorf("failed to open file with name \"%s\": %w", filename, err))
		} else {
			defer f.Close()
			buf.WriteTo(f)
		}
	}
}
