package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
	"sync"

	"bariq/ast"
	"bariq/sched"
)

const (
	HASH_OBJ         = "HASH"
	INT_OBJ          = "INTEGER"
	STRING_OBJ       = "STRING"
	BOOL_OBJ         = "BOOLEAN"
	NULL_OBJ         = "NULL"
	YIELD_VALUE_OBJ  = "YIELD_VALUE"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	TASK_OBJ         = "TASK_OBJ"
	GEN_OBJ          = "GEN_OBJ"
	ITER_OBJ         = "ITER_OBJ"
	ARRAY_OBJ        = "ARRAY"
	BUILTIN_OBJ      = "BUILTIN"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}
type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Hashable interface {
	HashKey() HashKey
}
type HashPair struct {
	Key   Object
	Value Object
}
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(
			pairs,
			fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()),
		)
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin Function" }

type Iteration struct {
	Done bool
	Val  Object
}

func (f *Iteration) Type() ObjectType { return GEN_OBJ }
func (f *Iteration) Inspect() string {
	var out bytes.Buffer
	return out.String()
}

type Generator struct {
	Env   *Env
	Fn    *Function
	Index int
	Done  bool
	Value Object
}

func (g *Generator) Type() ObjectType { return GEN_OBJ }
func (g *Generator) Inspect() string {
	var out bytes.Buffer
	return out.String()
}

type Task struct {
	Spawned *sched.Task[Object]
}

func (f *Task) Type() ObjectType { return TASK_OBJ }
func (f *Task) Inspect() string {
	var out bytes.Buffer
	return out.String()
}

type Function struct {
	Parameters []*ast.Ident
	Body       *ast.BlockStmt
	Env        *Env
	IsAsync    bool
	IsGen      bool
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type Env struct {
	mu    sync.RWMutex
	store map[string]Object
	outer *Env
}

func NewEnclosedEnv(outer *Env) *Env {
	// INFO: take care that it creates a new env
	env := NewEnv()
	env.outer = outer
	return env
}

func NewEnv() *Env {
	s := make(map[string]Object)
	// WARN: &env is the pointer to the pointer var containing pioner value
	// and not the actual pointer value
	env := &Env{store: s, outer: nil}
	// fmt.Printf("created env: %+v, with addr %p\n", env, env)
	return env
}

func (e *Env) Get(name string) (Object, bool) {
	e.mu.RLock()
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	e.mu.RUnlock()
	return obj, ok
}

func (e *Env) Set(name string, val Object) Object {
	e.mu.Lock()
	e.store[name] = val
	e.mu.Unlock()
	return val
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type YieldValue struct {
	Value Object
}

func (yv *YieldValue) Inspect() string  { return yv.Value.Inspect() }
func (yv *YieldValue) Type() ObjectType { return YIELD_VALUE_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Array struct {
	Elements []Object
}

func (arr *Array) Inspect() string {
	var out bytes.Buffer
	elmnts := []string{}
	for _, e := range arr.Elements {
		elmnts = append(elmnts, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elmnts, ", "))
	out.WriteString("]")
	return out.String()
}
func (arr *Array) Type() ObjectType { return ARRAY_OBJ }

type String struct {
	Value string
}

func (i *String) Inspect() string  { return i.Value }
func (i *String) Type() ObjectType { return STRING_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INT_OBJ }

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string  { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) Type() ObjectType { return BOOL_OBJ }

type Null struct{}

func (i *Null) Inspect() string  { return "null" }
func (i *Null) Type() ObjectType { return NULL_OBJ }
