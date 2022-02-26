package basic

func IfElse(ok bool, v1 interface{}, v2 interface{}) interface{} {
	if ok {
		return v1
	}
	return v2
}
