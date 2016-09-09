+++
Categories = ["go", "interpreter"]
Description = "A tiny interpreter"
Tags = ["go", "interpreter", "pygo"]
date = "2016-09-09T10:37:07+02:00"
title = "A tiny python-like interpreter"

+++

Last episode saw me slowly building up towards setting the case for
a `pygo` interpreter: a `python` interpreter in `Go`.

Still following the [Python interpreter written in Python](http://www.aosabook.org/en/500L/a-python-interpreter-written-in-python.html)
blueprints, let me first do (yet another!) little detour:
let me build a tiny (`python`-like) interpreter.

## A Tiny Interpreter

This tiny interpreter will understand three instructions:

- `LOAD_VALUE`
- `ADD_TWO_VALUES`
- `PRINT_ANSWER`

As stated before, my interpreter doesn't care about lexing, parsing nor compiling.
It has just, somehow, got the instructions from somewhere.

So, let's say I want to interpret:

```python
7 + 5
```

The instruction set to interpret would look like:

```go
code := Code{
	Prog: []Instruction{
		OpLoadValue, 0, // load first number
		OpLoadValue, 1, // load second number
		OpAdd,
		OpPrint,
	},
	Numbers: []int{7, 5},
}

var interp Interpreter
interp.Run(code)
```

----
The astute reader will probably notice I have slightly departed from
AOSA's `python` code.
In the book, each instruction is actually a 2-tuple `(Opcode, Value)`.
Here, an instruction is just a stream of "integers", being (implicitly) either
an `Opcode` or an operand.

----

The `CPython` interpreter is a _stack machine_.
Its instruction set reflects that implementation detail and thus,
our tiny interpreter implementation will have to cater for this aspect too:

```go
type Interpreter struct {
	stack stack
}

type stack struct {
	stk []int
}
```

Now, the interpreter has to actually run the code, iterating over each
instructions, pushing/popping values to/from the stack, according to
the current instruction.
That's done in the `Run(code Code)` method:

```go
func (interp *Interpreter) Run(code Code) {
	prog := code.Prog
	for pc := 0; pc < len(prog); pc++ {
		op := prog[pc].(Opcode)
		switch op {
		case OpLoadValue:
			pc++
			val := code.Numbers[prog[pc].(int)]
			interp.stack.push(val)
		case OpAdd:
			lhs := interp.stack.pop()
			rhs := interp.stack.pop()
			sum := lhs + rhs
			interp.stack.push(sum)
		case OpPrint:
			val := interp.stack.pop()
			fmt.Println(val)
		}
	}
}
```

And, yes, sure enough, running this:

```go
func main() {
	code := Code{
		Prog: []Instruction{
			OpLoadValue, 0, // load first number
			OpLoadValue, 1, // load second number
			OpAdd,
			OpPrint,
		},
		Numbers: []int{7, 5},
	}

	var interp Interpreter
	interp.Run(code)
}
```

outputs:

```sh
$> go run ./cmd/tiny-interpreter/main.go
12
```

The full code is here: [github.com/sbinet/pygo/cmd/tiny-interp](https://github.com/sbinet/pygo/blob/4938a159499724011a7175a4f344560372ccd468/cmd/tiny-interp/main.go).

## Variables

The `AOSA` article sharply notices that, even though this `tiny-interp` interpreter
is quite limited, its overall architecture and *modus operandi* are quite comparable
to how the real `python` interpreter works.

Save for variables.
`tiny-interp` doesn't do variables.
Let's fix that.

Consider this code fragment:
```python
a = 1
b = 2
print(a+b)
```

`tiny-interp` needs to be modified so that:

- values can be associated to names (variables), and
- new `Opcodes` need to be added to describe these associations.

Under these new considerations, the above code fragment would be compiled
down to the following program:

```go
func main() {
	code := Code{
		Prog: []Instruction{
			OpLoadValue, 0,
			OpStoreName, 0,
			OpLoadValue, 1,
			OpStoreName, 1,
			OpLoadName, 0,
			OpLoadName, 1,
			OpAdd,
			OpPrint,
		},
		Numbers: []int{1, 2},
		Names:   []string{"a", "b"},
	}

	interp := New()
	interp.Run(code)
}
```

The new opcodes `OpStoreName` and `OpLoadName` respectively store the current
value on the stack with some variable name (the index into the `Names` slice) and
load the value (push it on the stack) associated with the current variable.

The `Interpreter` now looks like:

```go
type Interpreter struct {
	stack stack
	env   map[string]int
}
```

where `env` is the association of variable names with their current value.

The `Run` method is then modified to handle `OpLoadName` and `OpStoreName`:
```diff
 func (interp *Interpreter) Run(code Code) {
@@ -63,6 +79,16 @@ func (interp *Interpreter) Run(code Code) {
                case OpPrint:
                        val := interp.stack.pop()
                        fmt.Println(val)
+               case OpLoadName:
+                       pc++
+                       name := code.Names[prog[pc].(int)]
+                       val := interp.env[name]
+                       interp.stack.push(val)
+               case OpStoreName:
+                       pc++
+                       name := code.Names[prog[pc].(int)]
+                       val := interp.stack.pop()
+                       interp.env[name] = val
                }
        }
 }
```

At this point, `tiny-interp` correctly handles variables:

```sh
$> tiny-interp
3
```

which is indeed the expected result.

The complete code is here: [github.com/sbinet/pygo/cmd/tiny-interp](https://github.com/sbinet/pygo/blob/79e9815cafa9c32e898141858502931acb3daf05/cmd/tiny-interp/main.go)

## Control flow

`tiny-interp` is already quite great.
I think.
But there is at least one glaring defect.
Consider:

```python
def cond():
	x = 3
	if x < 5:
		return "yes"
	else:
		return "no"
```

`tiny-interp` doesn't handle conditionals.
It's also completely ignorant about loops.
In a nutshell, there is **no control flow** in `tiny-interp`.
Yet.

To properly implement control flow, though, `tiny-interp` will need
to grow a new concept: activation records, also known as `Frames`.

Stay tuned...
