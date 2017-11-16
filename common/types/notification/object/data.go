package object

type DataObject struct {
	Module string      `json:"module" mapstructure:"module"`
	Value  interface{} `json:"value" mapstructure:"value"`
}

func NewDataObject(module string, value interface{}) *DataObject {
	return &DataObject{Module: module, Value: value}
}
