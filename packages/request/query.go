package request

type RowForm struct {
	Columns string `json:"columns"`
}

type PaginatorForm struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListForm struct {
	Name string `json:"name"` //table name
	PaginatorForm
	RowForm
}

type GetList struct {
	ListForm
	Order any `json:"order"`
	Where any `json:"where"`
}

type FileType struct {
	Name  string `json:"name"`  //file name
	Type  string `json:"type"`  //data type if exist
	Value string `json:"value"` //if file name is null, save result to value
}

func NewFileType(fileName string) *FileType {
	var str FileType
	str.Name = fileName
	return &str
}

// Omitempty if input type zero value then return nil
func Omitempty(v any) any {
	if v == nil {
		return nil
	}
	var r any
	r = v
	switch v.(type) {
	case int:
		if v.(int) == 0 {
			return nil
		}
	case float64:
		if v.(float64) == 0 {
			return nil
		}
	case int64:
		if v.(int64) == 0 {
			return nil
		}
	case string:
		if v.(string) == "" {
			return nil
		}
	case []byte:
		if v.([]byte) == nil {
			return nil
		}
	default:
		return nil
	}
	return r
}

// BlockIdOrHash
// block id or block hash must be only choose one
type BlockIdOrHash struct {
	Id   int64  `json:"id,omitempty"`
	Hash string `json:"hash,omitempty"`
}
