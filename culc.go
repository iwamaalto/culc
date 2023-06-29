// calc0.go : 電卓プログラム
//
//	Copyright (C) 2014-2021 Makoto Hiroi
package culc

import (
	"fmt"
	"os"
	"strings"
	"text/scanner"
)

// 字句解析
type Lex struct {
	scanner.Scanner
	Token rune
}

// トークンを求める
func (lex *Lex) getToken() {
	lex.Token = lex.Scan()
}

// 因子の処理
func factor(lex *Lex) float64 {
	switch lex.Token {
	case '(':
		lex.getToken()
		val := expression(lex)
		if lex.Token != ')' {
			panic(fmt.Errorf("')' expected"))
		}
		lex.getToken()
		return val
	case '+':
		lex.getToken()
		return factor(lex)
	case '-':
		lex.getToken()
		return (-factor(lex))
	case scanner.Int, scanner.Float:
		var n float64
		fmt.Sscan(lex.TokenText(), &n)
		lex.getToken()
		return n
	case scanner.Ident:
		text := lex.TokenText()
		if text == "quit" {
			panic(text)
		}
		fallthrough
	default:
		panic(fmt.Errorf("unexpected token: %v", lex.TokenText()))
	}
}

// 項の処理
func term(lex *Lex) float64 {
	val := factor(lex)
	for {
		switch lex.Token {
		case '*':
			lex.getToken()
			val *= factor(lex)
		case '/':
			lex.getToken()
			val /= factor(lex)
		default:
			return val
		}
	}
}

// 式の処理
func expression(lex *Lex) float64 {
	val := term(lex)
	for {
		switch lex.Token {
		case '+':
			lex.getToken()
			val += term(lex)
		case '-':
			lex.getToken()
			val -= term(lex)
		default:
			return val
		}
	}
}

// 式の入力と評価
func toplevel(lex *Lex) (r bool) {
	r = false
	defer func() {
		err := recover()
		if err != nil {
			mes, ok := err.(string)
			if ok && mes == "quit" {
				r = true
			} else {
				fmt.Fprintln(os.Stderr, err)
				for lex.Token != ';' {
					lex.getToken()
				}
			}
		}
	}()
	for {
		fmt.Print("Calc> ")
		lex.getToken()
		val := expression(lex)
		if lex.Token != ';' {
			panic(fmt.Errorf("invalid expression"))
		} else {
			fmt.Println(val)
		}
	}
	return r
}

func Culc(expr string) string {
	var lex Lex
	lex.Init(strings.NewReader(expr))
	lex.getToken()                // 追加: 式の先頭のトークンを取得
	result := expression(&lex)    // Change this line
	if lex.Token != scanner.EOF { // 追加: 入力がすべて解析されたことをチェック
		return "ERROR: Invalid expression"
	}
	return fmt.Sprintf("%.2f", result)
}
