package token

type TokenType string

type Token struct { //文本中读取出的单个字符串的类型和值
	Type    TokenType
	Literal string
}

const (
	ILIEGAL = "ILIEGAL"
	EOF     = "EOF"

	//标识符+字面量
	IDENT = "IDENT" //字母或下划线组成的用户定义标识符
	INT   = "INT"
	//运算符
	ASSIGN = "="
	PLUS   = "+"
	//分隔符
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	//关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{ //关键字
	"fn":  FUNCTION,
	"let": LET,
}

func LookUpIdent(ident string) TokenType { //区分关键字和用户定义标识符
	if tok, ok := keywords[ident]; ok {
		return tok //关键字类型
	}
	return IDENT //用户定义标识符类型名
}
