package parser

import (
	"fmt"
	"testing"

	"bariq/ast"
	"bariq/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}
	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatal("ParseProgram() returned nil")
		}
		if len(program.Stmts) != 1 {
			t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
		}
		stmt := program.Stmts[0]
		if !testLetStmt(t, stmt, tt.expectedIdent) {
			return
		}
		letStmt, ok := stmt.(*ast.LetStmt)
		if !ok {
			t.Errorf("s is not *ast.LetStmt. got %T", letStmt)
		}
		val := letStmt.Value
		fmt.Println("val is :", val, i, letStmt)
		if !testLiteralExpr(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStmt(t *testing.T, s ast.Stmt, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got %s ", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStmt)
	if !ok {
		t.Errorf("s is not *ast.LetStmt. got %T", letStmt)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf(
			"letStmt.Name.Vale not '%s'. got '%s'",
			name,
			letStmt.Name.Value,
		)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf(
			"letStmt.Name.TokenLiteral() not '%s'. got %s",
			name,
			letStmt.Name.TokenLiteral(),
		)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return add(15);
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Stmts) != 3 {
		t.Fatalf("expected 3 stmts but got %d", len(program.Stmts))
	}
	for _, stmt := range program.Stmts {
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStmt, got %T ", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf(
				"returnStmt.TokenLiteral not 'return',got %q ",
				returnStmt.TokenLiteral(),
			)
		}

	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser Error %q", msg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	ident, ok := stmt.Expr.(*ast.Ident)
	if !ok {
		t.Fatalf("exp is not as Ident ,got %T", stmt.Expr)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident value is not %s, got %s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf(
			"ident TokenLiteral is not %s, got %s",
			"foobar",
			ident.TokenLiteral(),
		)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	literal, ok := stmt.Expr.(*ast.IntLiteral)
	if !ok {
		t.Fatalf("exp is not as Ident ,got %T", stmt.Expr)
	}
	if literal.Value != 5 {
		t.Errorf("ident value is not %s, got %d", "5", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf(
			"ident TokenLiteral is not %s, got %s",
			"5",
			literal.TokenLiteral(),
		)
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatal("ParseProgram() returned nil")
		}
		if len(program.Stmts) != 1 {
			t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
		}
		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
		}
		exp, ok := stmt.Expr.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("exp is not as prefixExpr ,got %T", stmt.Expr)
		}
		if exp.Operator != tt.operator {
			t.Errorf("ident value is not %s, got %s", tt.operator, exp.Operator)
		}
		if !testLiteralExpr(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegeralLiteral(t *testing.T, il ast.Expr, value int64) bool {
	i, ok := il.(*ast.IntLiteral)
	if !ok {
		t.Errorf("il is not integerLiteral, got %T ", il)
		return false
	}
	if i.Value != value {
		t.Errorf("integer.value not %d, got %d", value, i.Value)
		return false
	}
	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.Literal not %d, got %se ", value, i.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 == 5;", 5, "==", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"true == false;", true, "==", false},
	}
	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatal("ParseProgram() returned nil")
		}
		if len(program.Stmts) != 1 {
			t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
		}
		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
		}
		if !testInfixExpr(
			t,
			stmt.Expr,
			tt.leftValue,
			tt.operator,
			tt.rightValue,
		) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},

		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},

		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"1 / (2 + 3)",
			"(1 / (2 + 3))",
		},
		{
			"-(2 + 3)",
			"(-(2 + 3))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(1 * 2) + d",
			"((a + add((1 * 2))) + d)",
		},
		{
			"a * [1,2,3,4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2],b[1],2 * [1,2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected %q but got %q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expr, value string) bool {
	ident, ok := exp.(*ast.Ident)
	if !ok {
		t.Errorf("exp not an ident, got %T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got %s", value, ident.Value)
		return false
	}
	if ident.TokenLiteral() != value {
		t.Errorf(
			"ident.TokenLiteral not %s, got %s ",
			value,
			ident.TokenLiteral(),
		)
		return false
	}
	return true
}

func testLiteralExpr(t *testing.T, exp ast.Expr, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegeralLiteral(t, exp, int64(v))
	case int64:
		return testIntegeralLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled, got %T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expr, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got %T ", exp)
		return false
	}
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got %t", value, bo.Value)
		return false
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not  %t got %s ", value, bo.TokenLiteral())
		return false
	}
	return true
}

// use this to test infix above
func testInfixExpr(
	t *testing.T,
	exp ast.Expr,
	left any,
	operator string,
	right any,
) bool {
	opExp, ok := exp.(*ast.InfixExpr)
	if !ok {
		t.Errorf("type of exp not infix expr, got %T", exp)
	}
	if !testLiteralExpr(t, opExp.Left, left) {
		return false
	}
	if opExp.Operator != operator {
		t.Errorf("operator is not %s, got %s ", opExp.Operator, operator)
	}
	if !testLiteralExpr(t, opExp.Right, right) {
		return false
	}
	return true
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Stmts) != 1 {
			t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
		}
		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
		}
		boolean, ok := stmt.Expr.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp is not as boolean ,got %T", stmt.Expr)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf(
				"ident value is not %t, got %t",
				tt.expectedBoolean,
				boolean.Value,
			)
		}
	}
}

func TestAwaitExpr(t *testing.T) {
	input := `await (x < y)`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}

	exp, ok := stmt.Expr.(*ast.AwaitExpr)
	if !ok {
		t.Fatalf("stmt expr is not await expr, got %T ", stmt.Expr)
	}

	if !testInfixExpr(t, exp.Arg, "x", "<", "y") {
		return
	}
}

func TestIfExpr(t *testing.T) {
	input := `if (x < y) { x }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	exp, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("stmt expr is not if expr, got %T ", stmt.Expr)
	}
	if !testInfixExpr(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Stmts) != 1 {
		t.Errorf(
			"expected 1 stmt for Consequence,got %d ",
			len(exp.Consequence.Stmts),
		)
	}
	consequence, ok := exp.Consequence.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("stmt[0] is not exprStmt,got %T ", exp.Consequence.Stmts[0])
	}
	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}
	if exp.Alternative != nil {
		t.Errorf(
			"exp.Alternative statemnts was not nil. got %+v ",
			exp.Alternative,
		)
	}
}

func TestIfElseExpr(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	exp, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("stmt expr is not if expr, got %T ", stmt.Expr)
	}
	if !testInfixExpr(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Stmts) != 1 {
		t.Errorf(
			"expected 1 stmt for Consequence,got %d ",
			len(exp.Consequence.Stmts),
		)
	}
	consequence, ok := exp.Consequence.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("stmt[0] is not exprStmt,got %T ", exp.Consequence.Stmts[0])
	}
	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}
	if len(exp.Alternative.Stmts) != 1 {
		t.Errorf(
			"exp Alternative stmts != 1, got %d ",
			len(exp.Alternative.Stmts),
		)
	}

	alter, ok := exp.Alternative.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("stmt[0] is not exprStmt,got %T ", exp.Alternative.Stmts[0])
	}
	if !testIdentifier(t, alter.Expr, "y") {
		return
	}
}

func TestGenFunctionLiteralParsing(t *testing.T) {
	input := ` fn gen (x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	fn, ok := stmt.Expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt expr is not FunctionLiteral expr, got %T ", stmt.Expr)
	}
	fmt.Printf("fn.Async: %v\n", fn.Async)
	fmt.Printf("fn.Gen: %v\n", fn.Gen)
	if !fn.Gen {
		t.Fatalf(
			"stmt expr is not gen FunctionLiteral expr, got %T",
			fn,
		)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Parameters))
	}
	testLiteralExpr(t, fn.Parameters[0], "x")
	testLiteralExpr(t, fn.Parameters[1], "y")
	if len(fn.Body.Stmts) != 1 {
		t.Fatalf("expected 1 stmt for func body, got %d", len(fn.Body.Stmts))
	}
	bodyStmt, ok := fn.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf(
			"func body stmt is not an expr stmt. got %T",
			fn.Body.Stmts[0],
		)
	}
	testInfixExpr(t, bodyStmt.Expr, "x", "+", "y")
}

func TestAsyncFunctionLiteralParsing(t *testing.T) {
	input := `async fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	fn, ok := stmt.Expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt expr is not FunctionLiteral expr, got %T ", stmt.Expr)
	}
	fmt.Println(fn.Async)
	if !fn.Async {
		t.Fatalf(
			"stmt expr is not async FunctionLiteral expr, got %T",
			fn,
		)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Parameters))
	}
	testLiteralExpr(t, fn.Parameters[0], "x")
	testLiteralExpr(t, fn.Parameters[1], "y")
	if len(fn.Body.Stmts) != 1 {
		t.Fatalf("expected 1 stmt for func body, got %d", len(fn.Body.Stmts))
	}
	bodyStmt, ok := fn.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf(
			"func body stmt is not an expr stmt. got %T",
			fn.Body.Stmts[0],
		)
	}
	testInfixExpr(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	exp, ok := stmt.Expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt expr is not FunctionLiteral expr, got %T ", stmt.Expr)
	}

	if len(exp.Parameters) != 2 {
		t.Fatalf("expected 2 params, got %d", len(exp.Parameters))
	}
	testLiteralExpr(t, exp.Parameters[0], "x")
	testLiteralExpr(t, exp.Parameters[1], "y")
	if len(exp.Body.Stmts) != 1 {
		t.Fatalf("expected 1 stmt for func body, got %d", len(exp.Body.Stmts))
	}
	bodyStmt, ok := exp.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf(
			"func body stmt is not an expr stmt. got %T",
			exp.Body.Stmts[0],
		)
	}
	testInfixExpr(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn () {};", expectedParams: []string{}},
		{input: "fn (x) {};", expectedParams: []string{"x"}},
		{input: "fn (x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Stmts[0].(*ast.ExprStmt)
		function := stmt.Expr.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf(
				"length doesn't match, epxected %d, got %d\n",
				len(function.Parameters),
				len(tt.expectedParams),
			)
		}
		for i, ident := range tt.expectedParams {
			testLiteralExpr(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExprParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	exp, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("stmt expr is not a call expr. got %T", stmt.Expr)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Args) != 3 {
		t.Fatalf("wrong len of args, got %d", len(exp.Args))
	}
	testLiteralExpr(t, exp.Args[0], 1)
	testInfixExpr(t, exp.Args[1], 2, "*", 3)
	testInfixExpr(t, exp.Args[2], 4, "+", 5)
}

func TestStringLiteralExpr(t *testing.T) {
	input := `"foobar";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	str, ok := stmt.Expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("s is not *ast.StringLiteral. got %T", stmt)
	}
	if str.Value != "foobar" {
		t.Errorf("literal value is not foobar,got %s", str.Value)
	}
}

func TestArrayLiteralExpr(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	array, ok := stmt.Expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("s is not *ast.ArrayLiteral. got %T", stmt.Expr)
	}
	if len(array.Elmnts) != 3 {
		t.Errorf("array elements are not 3, got %d", len(array.Elmnts))
	}
	testIntegeralLiteral(t, array.Elmnts[0], 1)
	testInfixExpr(t, array.Elmnts[1], 2, "*", 2)
	testInfixExpr(t, array.Elmnts[2], 3, "+", 3)
}

func TestParsingIndexExpr(t *testing.T) {
	input := `myArray[1 + 1];`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	// if len(program.Stmts) != 1 {
	// 	for _, s := range program.Stmts {
	// 		fmt.Println(s.String())
	// 	}
	// 	t.Fatalf(
	// 		"expected 1 stmts but got %d",
	// 		len(program.Stmts),
	// 	)
	// }
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	ie, ok := stmt.Expr.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("s is not *ast.IndexExpr. got %T", stmt.Expr)
	}
	if !testIdentifier(t, ie.Left, "myArray") {
		return
	}
	if !testInfixExpr(t, ie.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralBooleanKeys(t *testing.T) {
	input := `{true:1,false:2}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	hl, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("s is not *ast.HashLiter. got %T", stmt.Expr)
	}
	if len(hl.Pairs) != 2 {
		t.Errorf("array elements are not 2, got %d", len(hl.Pairs))
	}
	expected := map[bool]int64{
		true:  1,
		false: 2,
	}
	for k, v := range hl.Pairs {
		lit, ok := k.(*ast.Boolean)
		if !ok {
			t.Errorf("key is not ast.Boolean Literal, got %T", k)
		}
		expVal := expected[lit.Value]
		testIntegeralLiteral(t, v, expVal)
	}
}

func TestParsingHashLiteralStringKeys(t *testing.T) {
	input := `{"one":1,"two":2,"three":3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	hl, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("s is not *ast.HashLiter. got %T", stmt.Expr)
	}
	if len(hl.Pairs) != 3 {
		t.Errorf("array elements are not 3, got %d", len(hl.Pairs))
	}
	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	for k, v := range hl.Pairs {
		lit, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.String Literal, got %T", k)
		}
		expVal := expected[lit.String()]
		testIntegeralLiteral(t, v, expVal)
	}
}

func TestParsingHashLiteralIntegerKeys(t *testing.T) {
	input := `{1:1,2:2,3:3}`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	hl, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("s is not *ast.HashLiter. got %T", stmt.Expr)
	}
	if len(hl.Pairs) != 3 {
		t.Errorf("array elements are not 3, got %d", len(hl.Pairs))
	}
	expected := map[int64]int64{
		1: 1,
		2: 2,
		3: 3,
	}
	for k, v := range hl.Pairs {
		lit, ok := k.(*ast.IntLiteral)
		if !ok {
			t.Errorf("key is not ast.Boolean Literal, got %T", k)
		}
		expVal := expected[lit.Value]
		testIntegeralLiteral(t, v, expVal)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{};"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	hl, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("s is not *ast.HashLiter. got %T", stmt.Expr)
	}

	if len(hl.Pairs) != 0 {
		t.Errorf("array elements are not 0, got %d", len(hl.Pairs))
	}
}

func TestParsingHashLiteralWithExprs(t *testing.T) {
	input := `{"one":0 +1,"two": 10 - 8,"three":15/5};`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("expected 1 stmts but got %d", len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("s is not *ast.exprStmt. got %T", stmt)
	}
	hl, ok := stmt.Expr.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("s is not *ast.HashLiter. got %T", stmt.Expr)
	}

	if len(hl.Pairs) != 3 {
		t.Errorf("array elements are not 3, got %d", len(hl.Pairs))
	}

	tests := map[string]func(ast.Expr){
		"one": func(e ast.Expr) {
			testInfixExpr(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expr) {
			testInfixExpr(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expr) {
			testInfixExpr(t, e, 15, "/", 5)
		},
	}
	for k, v := range hl.Pairs {
		lit, ok := k.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not string literal, got %T", k)
		}
		testFunc, ok := tests[lit.String()]
		if !ok {
			t.Errorf("No test func is associated with key %q", lit.String())
			continue
		}
		testFunc(v)
	}
}
