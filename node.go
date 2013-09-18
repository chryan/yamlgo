package yaml

import (
	"reflect"
)

type Node struct {
	Tag string
	Value interface{}

	nodeType reflect.Type
}

func (n *Node) IsNull() bool {
	return n.Value == nil;
}

func (n *Node) IsScalar() bool {
	if n.nodeType == nil {
		return false
	}
	switch n.nodeType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.String:
		return true
	}
	return false
}

func (n *Node) IsSequence() bool {
	if n.nodeType == nil {
		return false
	}
	switch n.nodeType.Kind() {
	case reflect.Array, reflect.Slice:
		return true
	}
	return false
}

func (n *Node) IsMap() bool {
	if n.nodeType == nil {
		return false
	}
	switch n.nodeType.Kind() {
	case reflect.Map, reflect.Struct:
		return true
	}
	return false
}

func (n *Node) Set(value interface{}) {
	if value == nil {
		n.nodeType = nil
		n.Value = nil
		return
	}

	n.nodeType = reflect.TypeOf(value)
	
	if kind := n.nodeType.Kind(); kind == reflect.Ptr || kind == reflect.Interface {
		rval := reflect.ValueOf(value).Elem()
		n.nodeType = rval.Type()
		n.Value = rval.Interface()
	} else {
		n.Value = value
	}
}