package evaluator

import (
	"context"
	"fmt"

	"bariq/ast"
	"bariq/object"
	"bariq/sched"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func Eval(node ast.Node, env *object.Env) object.Object {
	switch node := node.(type) {
	// Stmts
	case *ast.Program:
		return evalProgram(node.Stmts, env)
	// diff between expstmt and retstmt is HERE
	case *ast.ExprStmt:
		return Eval(node.Expr, env)
	case *ast.BlockStmt:
		return evalBlockStmt(node, env)
	case *ast.ReturnStmt:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStmt:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		// fmt.Println("Env of: ", node.Name.Value, env)
	// Exprs

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.IntLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.ArrayLiteral:
		elms := evalExprs(node.Elmnts, env)
		if len(elms) == 1 && isError(elms[0]) {
			return elms[0]
		}
		return &object.Array{Elements: elms}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.Boolean:
		// why to create an object every time
		// where you can just declare two values
		// and use them every time you need them
		return toBoolObj(node.Value)
	case *ast.PrefixExpr:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalInfixExpr(node.Operator, left, right)
	case *ast.YieldExpr:
		fmt.Printf("yield node: %v\n", node)
		val := Eval(node.Arg, env)
		if isError(val) {
			return val
		}
		return &object.YieldValue{Value: val}

	case *ast.AwaitExpr:
		return evalAwaitExpr(node, env)
		// return evalIfExpr(node, env)
	case *ast.IfExpr:
		return evalIfExpr(node, env)

	case *ast.Ident:
		return evalIdent(node, env)

	case *ast.FunctionLiteral:
		isAsync := node.Async
		params := node.Parameters
		body := node.Body
		isGen := node.Gen
		// fmt.Printf("created func with env addr: %p,and body: %+v\n", env, body)
		// INFO: when the function declared, the env is assigned
		return &object.Function{Parameters: params, Env: env, Body: body, IsAsync: isAsync, IsGen: isGen}

	case *ast.CallExpr:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExprs(node.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		// NOTE: to make this dynamic scope, instead of usnig function env
		//  to create a new env,pass the curent env and use it to create it.
		// return applyFunc(env, function, args)

		// TODO: check matching args
		// fmt.Printf("len(args): %v\n", len(node.Args))

		return applyFunc(function, args)

	case *ast.IndexExpr:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpr(left, index)
	}
	return nil
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Env) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for kn, vn := range node.Pairs {
		k := Eval(kn, env)
		if isError(k) {
			return k
		}
		hashKey, ok := k.(object.Hashable)
		if !ok {
			return newError("unusable as hashKey: %s", k.Type())
		}
		val := Eval(vn, env)
		if isError(val) {
			return val
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: k, Value: val}
	}
	return &object.Hash{Pairs: pairs}
}

func applyFunc(
	// env *object.Env,
	fn object.Object,
	args []object.Object,
) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		// fmt.Printf("fn being applied %+v, env_addr: %p\n", fn, fn.Env)
		// extendedEnv := extendedDynamicEnv(env, fn, args)
		extendedEnv := extendedStaticEnv(fn, args)
		if fn.IsAsync {
			fmt.Printf("fn.IsAsync: %v\n", fn.IsAsync)
			task := sched.Spawn(func(_ context.Context) (object.Object, error) {
				if fn.IsGen {
					fmt.Printf("fn.IsGen: %v\n", fn.IsGen)
					return &object.Generator{
						Fn:  fn,
						Env: extendedEnv,
					}, nil
				}
				evaluated := Eval(fn.Body, extendedEnv)
				return unwrapReturnValue(evaluated), nil
			})
			return &object.Task{
				Spawned: task,
			}
		}

		fmt.Printf("fn.IsGen: %v\n", fn.IsGen)
		if fn.IsGen {
			fmt.Printf("fn.IsGen: %v\n", fn.IsGen)
			return &object.Generator{
				Fn:  fn,
				Env: extendedEnv,
			}
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendedDynamicEnv(
	oldEnv *object.Env,
	fn *object.Function,
	args []object.Object,
) *object.Env {
	env := object.NewEnclosedEnv(oldEnv)
	for pIdx, p := range fn.Parameters {
		env.Set(p.Value, args[pIdx])
	}
	// fmt.Printf(
	// 	"created ext Env:%+v, addr: %p: \n%s\n",
	// 	env, env,
	// 	fn.Inspect(),
	// )
	return env
}

func extendedStaticEnv(
	fn *object.Function,
	args []object.Object,
) *object.Env {
	env := object.NewEnclosedEnv(fn.Env)
	for pIdx, p := range fn.Parameters {
		env.Set(p.Value, args[pIdx])
	}
	// fmt.Printf(
	// 	"created ext Env:%+v, addr: %p: \n%s\n",
	// 	env, env,
	// 	fn.Inspect(),
	// )
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}
	return obj
}

func evalExprs(exprs []ast.Expr, env *object.Env) []object.Object {
	var res []object.Object
	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		res = append(res, evaluated)
	}
	return res
}

func evalIndexExpr(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INT_OBJ:
		return evalArrayIndexExpr(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpr(left, index)
	default:
		return newError(`index operator not supported: %s `, left.Type())
	}
}

func evalHashIndexExpr(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalArrayIndexExpr(array, index object.Object) object.Object {
	arrObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrObj.Elements[idx]
}

func evalIdent(node *ast.Ident, env *object.Env) object.Object {
	if val, ok := env.Get(node.Value); ok {
		if val.Type() == object.GEN_OBJ {
			fmt.Printf("gen: %+v\n", val)
		}
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		// fmt.Println("found", node.Value)
		return builtin
	}
	return newError("ident not found: " + node.Value)
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalBlockStmt(block *ast.BlockStmt, env *object.Env) object.Object {
	var res object.Object
	// res will be the last evaluated stmt
	for _, stmt := range block.Stmts {

		res = Eval(stmt, env)
		// fmt.Println("block: ", block.String(), env)
		if res != nil {
			rt := res.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return res
			}
		}
	}
	return res
}

func evalAwaitExpr(a *ast.AwaitExpr, env *object.Env) object.Object {
	fmt.Println("Entered Await")
	evaluated := Eval(a.Arg, env)
	t, ok := evaluated.(*object.Task)
	if !ok {
		return evaluated
	}
	evalT, err := t.Spawned.Await()
	if err != nil {
		return newError("error has occured while awaiting - %+v", err)
	}
	return evalT
}

func evalIfExpr(ie *ast.IfExpr, env *object.Env) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpr(
	op string,
	left object.Object,
	right object.Object,
) object.Object {
	switch {
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpr(op, left, right)
	case left.Type() == object.INT_OBJ && right.Type() == object.INT_OBJ:
		return evalIntegerInfixExpr(op, left, right)
	case op == "==":
		// use pointer comparison to compare between bools
		return toBoolObj(left == right)
	case op == "!=":
		// use pointer comparison to compare between bools
		return toBoolObj(left != right)
	case left.Type() != right.Type():
		return newError(
			"type mismatch: %s %s %s",
			left.Type(),
			op,
			right.Type(),
		)
	default:
		return newError(
			"unkown operator: %s %s %s",
			left.Type(),
			op,
			right.Type(),
		)
		// return NULL
	}
}

func evalPrefixExpr(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpr(right)
	case "-":
		return evalMinusOperatorExpr(right)
		// op not supported
	default:
		return newError("unkown operator: %s%s", op, right.Type())
		// return NULL
	}
}

func evalIntegerInfixExpr(
	op string,
	l object.Object,
	r object.Object,
) object.Object {
	lVal := l.(*object.Integer).Value
	rVal := r.(*object.Integer).Value
	switch op {
	case "+":
		return &object.Integer{Value: lVal + rVal}
	case "-":
		return &object.Integer{Value: lVal - rVal}
	case "*":
		return &object.Integer{Value: lVal * rVal}
	case "/":
		return &object.Integer{Value: lVal / rVal}
	case "<":
		return toBoolObj(lVal < rVal)
	case ">":
		return toBoolObj(lVal > rVal)
	case "==":
		return toBoolObj(lVal == rVal)
	case "!=":
		return toBoolObj(lVal != rVal)
	default:
		return newError(
			"unkown operator: %s %s %s",
			l.Type(),
			op,
			r.Type(),
		)
		// return NULL
	}
}

func evalStringInfixExpr(
	op string,
	l object.Object,
	r object.Object,
) object.Object {
	lVal := l.(*object.String).Value
	rVal := r.(*object.String).Value
	switch op {
	case "+":
		return &object.String{Value: lVal + rVal}
	default:
		return newError(
			"unkown operator: %s %s %s",
			l.Type(),
			op,
			r.Type(),
		)
		// return NULL
	}
}

func evalMinusOperatorExpr(r object.Object) object.Object {
	if r.Type() != object.INT_OBJ {
		return newError("unkown operator: -%s", r.Type())
		// return NULL
	}
	value := r.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpr(r object.Object) object.Object {
	switch r {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
		// int, ... any truthy val but true,false,null..
	default:
		return FALSE
	}
}

func toBoolObj(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(stmts []ast.Stmt, env *object.Env) object.Object {
	var res object.Object
	// res will be the last evaluated stmt
	for _, stmt := range stmts {

		res = Eval(stmt, env)
		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}
	return res
}
