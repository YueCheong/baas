package contract

//链码所能使用的语言类型
type ChaincodeLanguageType int

func (l ChaincodeLanguageType) Int() int {
	return int(l)
}

func (l ChaincodeLanguageType) String() string {
	if result, ok := chaincodeLanguageName[l]; ok {
		return result
	} else {
		return "UNKNOW"
	}
}

func (l ChaincodeLanguageType) Valid() bool {
	if _, ok := chaincodeLanguageName[l]; ok {
		return true
	} else {
		return false
	}
}

var chaincodeLanguageName = map[ChaincodeLanguageType]string{
	Golang: "Golang",
	Java:   "Java",
	Node:   "Node",
}

const (
	Golang ChaincodeLanguageType = iota
	Java
	Node
)
