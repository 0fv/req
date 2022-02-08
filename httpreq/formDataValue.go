package httpreq

type FormDataType uint8

const (
	FormDataTypeStr FormDataType = iota + 1
	FormDataTypeFile
)

type FormDataValue struct {
	Value    string       `json:"value,omitempty"`
	Type     FormDataType `json:"type,omitempty"`
	FileName string       `json:"fileName,omitempty"`
}
