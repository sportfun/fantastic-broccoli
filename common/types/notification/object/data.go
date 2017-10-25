package object

type DataObject struct {
	module string
	value  interface{}
}

func NewDataObject(module string, value interface{}) *DataObject {
	return &DataObject{module: module, value: value}
}

func (d *DataObject) Module() string {
	return d.module
}

func (d *DataObject) Value() interface{} {
	return d.value
}
