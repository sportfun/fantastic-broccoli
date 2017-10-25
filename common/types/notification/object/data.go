package object

type DataObject struct {
	module string
	value  interface{}
}

func NewDataObject(module string, value interface{}) *DataObject {
	return &DataObject{module: module, value: value}
}

func (dataObject *DataObject) Module() string {
	return dataObject.module
}

func (dataObject *DataObject) Value() interface{} {
	return dataObject.value
}
