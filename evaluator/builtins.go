package evaluator

import (
	"fmt"
	"time"

	"bariq/object"
)

var builtins map[string]*object.Builtin

func init() {
	builtins = map[string]*object.Builtin{
		"sleep": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError(
						"wrong number of args, got %d, want 1",
						len(args),
					)
				}
				val, ok := args[0].(*object.Integer)
				if !ok {
					return newError(
						"argument to `len` not supported, got %s",
						args[1].Type(),
					)
				}
				time.Sleep(time.Second * time.Duration(val.Value))
				return NULL
			},
		},
		"len": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError(
						"wrong number of args, got %d, want 1",
						len(args),
					)
				}
				switch arg := args[0].(type) {
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
		"puts": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL
			},
		},
		"first": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError(
						"wrong number of args, got %d, want 1",
						len(args),
					)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError(
						"argument to `first` not supported, got %s",
						args[0].Type(),
					)
				}
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}
				return NULL
			},
		},
		"last": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError(
						"wrong number of args, got %d, want 1",
						len(args),
					)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError(
						"argument to `last` not supported, got %s",
						args[0].Type(),
					)
				}
				arr := args[0].(*object.Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[len(arr.Elements)-1]
				}
				return NULL
			},
		},
		"tail": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError(
						"wrong number of args, got %d, want 1",
						len(args),
					)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError(
						"argument to `last` not supported, got %s",
						args[0].Type(),
					)
				}
				arr := args[0].(*object.Array)
				ln := len(arr.Elements)
				if ln > 0 {
					newElmnts := make([]object.Object, ln-1, ln-1)
					copy(newElmnts, arr.Elements[1:ln])
					return &object.Array{Elements: newElmnts}
				}
				return NULL
			},
		},
		// immutalbe
		"push": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError(
						"wrong number of args, got %d, want 2",
						len(args),
					)
				}

				if args[0].Type() != object.ARRAY_OBJ {
					return newError(
						"argument to `last` not supported, got %s",
						args[0].Type(),
					)
				}
				arr := args[0].(*object.Array)
				ln := len(arr.Elements)
				newElmnts := make([]object.Object, ln+1, ln+1)
				copy(newElmnts, arr.Elements)
				newElmnts[ln] = args[1]
				return &object.Array{Elements: newElmnts}
			},
		},

		"next": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError(
						"wrong number of args for next, got %d, want 1",
						len(args),
					)
				}

				switch arg := args[0].(type) {
				case *object.Generator:
					var res object.Object
					fmt.Printf("%+v is arg\n", arg)
					fmt.Printf("i is %d , arg.stmts: %v\n", arg.Index, arg.Fn.Body.Stmts[arg.Index:])
					for i, stmt := range arg.Fn.Body.Stmts[arg.Index:] {
						fmt.Printf("stmt: %v\n", stmt)
						res = Eval(stmt, arg.Env)
						// fmt.Println("block: ", block.String(), env)
						if res != nil {
							// fmt.Printf("res: %v\n", res.Type())
							rt := res.Type()
							if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
								break
							}
							if rt == object.YIELD_VALUE_OBJ {
								fmt.Printf("i is %d, idx is %d\n", i, arg.Index)
								arg.Index += i + 1
								arg.Value = res

								if v, ok := arg.Value.(*object.YieldValue); ok {

									fmt.Printf("v.Value: %v\n", v.Value)
									return &object.Iteration{Done: false, Val: v.Value}
								}
								fmt.Println("not yield")
							}
						}
						fmt.Printf("evaling %v is nil\n", stmt)

					}

					if v, ok := arg.Value.(*object.YieldValue); ok {
						return &object.Iteration{Done: true, Val: v.Value}
					}
					return newError("not an yield value")
				default:
					return newError("argument to `next` not supported, got %s with %s", args[0].Type(), args[0].Inspect())
				}
			},
		},
	}
}
