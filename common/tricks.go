package common

func If(b bool, trueValue interface{}, falseValue interface{}) interface{} {
	if b {
		return trueValue
	}
	return falseValue
}