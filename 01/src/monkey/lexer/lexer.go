package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  //输入字符串当前位置
	readPosition int  //输入字符串读取位置（当前位置的下一个）
	ch           byte //当前正在查看的字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() //读取下一个字符，position=0，readPosition=1
	return l
}

func (l *Lexer) readChar() { //读取一个字符
	if l.readPosition >= len(l.input) {
		l.ch = 0 //是否达到input末尾
	} else {
		l.ch = l.input[l.readPosition] //记录下一个字符
	}
	l.position = l.readPosition //更新位置
	l.readPosition += 1         //更新位置+1
}

func (l *Lexer) NextToken() token.Token { //转换当前*Lexer的正在查看的字符ch，返回为对应Token结构包含类型和值
	var tok token.Token
	l.skipWhitespace() //跳过空格等无意义分隔符
	switch l.ch {      //匹配，得到语法单元<类型，值>
	case '=':
		if l.peekChar() == '=' { // '=='，peekChar()仅查看下一个字符
			ch := l.ch
			l.readChar() //读取下一个字符并移动
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' { // '!='，peekChar()仅查看下一个字符
			ch := l.ch
			l.readChar() //读取下一个字符并移动
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ('('):
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0: //空
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) { //字母或下划线
			tok.Literal = l.readIdentifiler()         //读出对应的字母下划线串
			tok.Type = token.LookUpIdent(tok.Literal) //区分关键字和用户定义标识符
			return tok                                //位置已改变，拿到tok，直接退出
		} else if isDigit(l.ch) { //判数字
			tok.Literal = l.readNumber() //读出对应的字母下划线串
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILIEGAL, l.ch) //其他的字符统一报错
		}
	}

	l.readChar() //读取下一个字符
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token { //传入对应类型的名称和byte类型的值，转换为token结构体里
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool { //判断是不是变量名
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readIdentifiler() string { //读出对应的字母下划线串
	position := l.position //第一个字母或下划线开始位置
	for isLetter(l.ch) {   //继续读直到非字母或下划线的位置
		l.readChar()
	}
	return l.input[position:l.position] //读出对应的字母下划线串
}

func (l *Lexer) skipWhitespace() { //跳过空格等无意义分隔符
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar() //直接下一个
	}
}

func isDigit(ch byte) bool { //判断是不是变量名
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string { //读出对应的字母下划线串
	position := l.position //第一个字母或下划线开始位置
	for isDigit(l.ch) {    //继续读直到非字母或下划线的位置
		l.readChar()
	}
	return l.input[position:l.position] //读出对应的字母下划线串
}

func (l *Lexer) peekChar() byte { //超前搜索
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition] //查看下一个单词
	}
}
