package enum

// Enum 用map[int]string实现枚举
type Enum map[int]string

func (e Enum) Code(value string) int {
	for k, v := range e {
		if v == value {
			return k
		}
	}
	return -1
}

func (e Enum) Value(code int) string {
	value, _ := e[code]
	return value
}
