package crudgen

import "reflect"

type CRUDDescription struct {
	SrcFilePath string
	PackageName string
	TableName   string
	RepoType    reflect.Type
	ModelType   reflect.Type
}

type Attrib struct {
	Name   string
	Type   string
	Column string
	Order  int
}

type CRUDTemplateValues struct {
	Package       string
	Imports       map[string]string
	RepoTypeName  string
	ModelTypeName string
	InterfaceName string
	TableName     string
	PKeyFields    []Attrib
	InsertFields  []Attrib
	UpdateFields  []Attrib
}

type AttribList []Attrib

func (e AttribList) Len() int {
	return len(e)
}
func (e AttribList) Less(i, j int) bool {
	return e[i].Order < e[j].Order
}
func (e AttribList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
