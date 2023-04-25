package request

import "fmt"

type MapParams map[string]any

func (cp *MapParams) Get(key string) string {
	if _, ok := (*cp)[key]; !ok {
		return ""
	}
	return fmt.Sprintf("%v", (*cp)[key])
}

func (cp *MapParams) GetRaw(key string) any {
	return (*cp)[key]
}
