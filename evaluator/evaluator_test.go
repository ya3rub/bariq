package evaluator

import (
	"fmt"
	"testing"

	"bariq/lexer"
	"bariq/object"
	"bariq/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnv()
	return Eval(program, env)
}

func TestEvalIntegerExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-4", -4},
		{"(5 + 5 + 5 + 5 - 10 * 2) + 5 / 5", 1},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got %T (%+v)", obj, obj)
		return false

	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got %s want %s",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not integer. got %T (%+v)", obj, obj)
		return false

	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got %d want %d",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func TestEvalBooleanExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 1", false},
		{"1 != 2", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not integer. got %T (%+v)", obj, obj)
		return false

	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got %t want %t",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!false", true},
		{"!!false", false},
		{"!10", false},
		{"!!10", true},
		{"true == true", true},
		{"false == true", false},
		{"true != true", false},
		{"false != true", true},
		{"(1 < 2) != true", false},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 < 2) != false", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExprs(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (true) { 10 } else { 20 }", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL  got %T (%+v)", obj, obj)
		return false

	}
	return true
}

func TestReturnStmts(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9;return 2 * 5; 9;", 10},
		{`
			if(10 > 1) {
			if(10>1){
			return 10;
			}
			return 1;
			}
		`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5; false + true;5", "unkown operator: BOOLEAN + BOOLEAN"},
		{"false + true;", "unkown operator: BOOLEAN + BOOLEAN"},
		{"-true;", "unkown operator: -BOOLEAN"},
		{"if (10 > 1) {true + false;}", "unkown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "ident not found: foobar"},
		{`"hello" - "world";`, "unkown operator: STRING - STRING"},
		{`{"name":"y"}[fn(x){x}];`, "unusable as hash key: FUNCTION"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf(
				"no error object returned, got %T(%+v)",
				evaluated,
				evaluated,
			)
			continue
		}
		if errObj.Message != tt.expected {
			t.Errorf(
				"wrong error msg, expected %q, got %q",
				tt.expected,
				errObj.Message,
			)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a;let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evalutated := testEval(input)
	fn, ok := evalutated.(*object.Function)
	if !ok {
		t.Fatalf("object is not func, got %T (%+v) ", evalutated, evalutated)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong paramters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters[0].String() != "x" {
		t.Fatalf("paramete is not 'x', got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q got %q", expectedBody, fn.Body.String())
	}
}

func TestAsyncFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`let s = async fn(x) { sleep(5);5 };
			let task = s(5);
			puts("4");
			await(task)`,
			5,
		},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let i = fn(x) { x; }; i(5);", 5},
		{"let i = fn(x) { return x; }; i(5);", 5},
		{"let double = fn(x) { x * x; }; double(5);", 25},
		{"let add = fn(x, y) { x + x; }; add(5,5);", 10},
		{"let add = fn(x, y) { x + x; }; add(5 + 5, add(5,5));", 20},
		{"fn(x) { x; }(5);", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosure(t *testing.T) {
	input := `
	let newAddr = fn(x) {
	  fn(y) { x + y };
	};
	let addTwo = newAddr(2);
	addTwo(2);`
	testIntegerObject(t, testEval(input), 4)
}

func TestEvalStringExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello world"`, "hello world"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"hello" + " " + "world";`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("obj is not a string, got %T (%+v)", evaluated, evaluated)
	}
	if str.Value != "hello world" {
		t.Errorf(
			"string is wrong, expected %s, got %s",
			"hello world",
			str.Value,
		)
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("");`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(0)`, "argument to `len` not supported, got INTEGER"},
		{`len("one","two")`, "wrong number of args, got 2, want 1"},
	}
	for _, tt := range tests {
		fmt.Println("entered", tt.expected)
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error, got %T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("error message not matched, got %s, want %s", errObj.Message, expected)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`
	evaluated := testEval(input)
	res, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not an Array. got %T (%+v)", evaluated, evaluated)
	}
	if len(res.Elements) != 3 {
		t.Fatalf(
			"array has wrong number of elements. got %d ",
			len(res.Elements),
		)
	}
	testIntegerObject(t, res.Elements[0], 1)
	testIntegerObject(t, res.Elements[1], 4)
	testIntegerObject(t, res.Elements[2], 6)
}

func TestHashIndexExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`{"foo":5}["foo"]`, 5},
		{`{"foo":5}["bar"]`, nil},
		{`let key = "foo";{"foo":5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5:5}[5]`, 5},
		{`{true:5}[true]`, 5},
		{`{false:5}[false]`, 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayIndexExpr(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`[1,2,3][0]`, 1},
		{`[1,2,3][1]`, 2},
		{`[1,2,3][2]`, 3},
		{`[1,2,3][1 + 1]`, 3},
		{`let i = 0; [1][i];`, 1},
		{`let s = [1,2,3];s[0];`, 1},
		{`let s = [1,2,3];let i = s[0]; s[i];`, 2},
		{`let s = [1,2,3];s[0] + s[1] + s[2];`, 6},
		{`[1,2,3][3]`, nil},
		{`[1,2,3][-1]`, nil},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
	"one":10 - 9,
	two:1 +1,
	"thr"+"ee":6/2,
	4:4,
	true:5,
	false:6
	}
	`
	evaluated := testEval(input)
	res, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("eval didn't return a hash,got %T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(res.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got %d.", len(res.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := res.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestGenFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected *object.Iteration
	}{
		{
			`let s =  fn gen () { yield 2;yield 0;yield 6;yield 1; };
			let genr = s();
			next(genr);
			next(genr);
			next(genr);
			`,
			&object.Iteration{
				Val:  &object.Integer{Value: 6},
				Done: false,
			},
		},
		{
			`let s =  fn gen () { yield 2;yield 0;yield 6;yield 1; };
			let genr = s();
			next(genr);
			next(genr);
			next(genr);
			next(genr);
			`,
			&object.Iteration{
				Val:  &object.Integer{Value: 1},
				Done: false,
			},
		},
		{
			`let s =  fn gen () { yield 2;yield 0;yield 6;yield 1; };
			let genr = s();
			next(genr);
			next(genr);
			next(genr);
			next(genr);
			next(genr);
			`,
			&object.Iteration{
				Val:  &object.Integer{Value: 1},
				Done: true,
			},
		},
		{
			`
			let w = fn (){1}
			let s =  fn gen () {
					let q = w()
					yield q;
				};
			let genr = s();
			next(genr);
			`,
			&object.Iteration{
				Val:  &object.Integer{Value: 1},
				Done: false,
			},
		},
	}
	for _, tt := range tests {
		testIterationObject(t, testEval(tt.input), tt.expected)
	}
}

func testIterationObject(
	t *testing.T,
	obj object.Object,
	expected *object.Iteration,
) bool {
	result, ok := obj.(*object.Iteration)
	if !ok {
		t.Errorf("object is not integer. got %T (%+v)", obj, obj)
		return false

	}
	if result.Done != expected.Done {
		t.Errorf(
			"object has wrong done value. got %T want %T",
			result.Done,
			expected.Done,
		)
		return false
	}
	if result.Val.Inspect() != expected.Val.Inspect() {
		t.Errorf(
			"object has wrong done value. got %s want %s",
			result.Val.Inspect(),
			expected.Val.Inspect(),
		)
		return false
	}
	return true
}
