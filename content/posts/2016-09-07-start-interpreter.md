+++
Categories = ["go", "interpreter"]
Description = "How to write a Go interpreter"
Tags = ["go", "interpreter", "gonum", "gopherds"]
date = "2016-09-07T10:37:07+02:00"
title = "Starting a Go interpreter"

+++

In this series of posts, I'll try to explain how one can write an interpreter
in `Go` and for `Go`.
If, like me, you lack a bit in terms of interpreters know-how, you should be
in for a treat.

## Introduction

`Go` is starting to get traction in the science and data science communities.
And, why not?
`Go` is fast to compile and run, is statically typed and thus presents a nice
"edit/compile/run" development cycle.
Moreover, a program written in `Go` is easily deployable and cross-compilable
on a variety of machines and operating systems.

`Go` is also starting to have the foundation libraries for scientific work:

- [gonum/blas](https://github.com/gonum/blas)
- [gonum/lapack](https://github.com/gonum/lapack)
- [gonum/integrate](https://github.com/gonum/integrate)
- [gonum/matrix](https://github.com/gonum/matrix)
- [gonum/optimize](https://github.com/gonum/optimize)
- [gonum/plot](https://github.com/gonum/plot)
- [gonum/stat](https://github.com/gonum/stat)

And the data science community is bootstrapping itself around the [gopherds](https://github.com/gopherds)
community (slack channel: [#data-science](https://gophers.slack.com/messages/data-science)).

For data science, a central tool and workflow is the [Jupyter](https://jupyter.org) and its
notebook.
The Jupyter notebook provides a nice "REPL"-based workflow and the ability
to share algorithms, plots and results.
The REPL (Read-Eval-Print-Loop) allows people to engage fast exploratory
work of someone's data, quickly iterating over various algorithms or
different ways to interpret data.
For this kind of work, an interactive interpreter is paramount.

But `Go` is compiled and even if the compilation is lightning fast, a true
interpreter is needed to integrate well with a REPL-based workflow.

The [go-interpreter](https://github.com/go-interpreter) project (also available
on Slack: [#go-interpreter](https://gophers.slack.com/messages/go-interpreter))
is starting to work on that: implement a `Go` interpreter, in `Go` and for `Go`.
The first step is to design a bit this beast: [here](https://github.com/go-interpreter/proposal/issues/1).

Before going there, let's do a little detour: writing a (toy) interpreter
in `Go` for `Python`.
Why? you ask...
Well, there is a very nice article in the AOSA series:
[A Python interpreter written in Python](http://www.aosabook.org/en/500L/a-python-interpreter-written-in-python.html).
I will use it as a guide to gain a bit of knowledge in writing interpreters.

## PyGo: A (toy) Python interpreter

In the following, I'll show how one can write a toy `Python` interpreter in `Go`.
But first, let us define exactly what `pygo` will do.
`pygo` won't lex, parse and compile `Python` code.

No.
`pygo` will take directly the already compiled bytecode, produced with a
`python3` program, and then interpret the bytecode instructions:

```sh
shell> python3 -m compileall -l my-file.py
shell> pygo ./__pycache__/my-file.cpython-35.pyc
```

`pygo` will be a simple _bytecode interpreter_, with a main loop fetching
bytecode instructions and then executing them.
In pseudo `Go` code:

```go
func run(instructions []instruction) {
	for _, instruction := range instructions {
		switch inst := instruction.(type) {
			case opADD:
				// perform a+b
			case opPRINT:
				// print values
			// ...
		}
	}
}
```

`pygo` will export a few types to implement such an interpreter:

- a virtual machine `pygo.VM` that will hold the call stack of frames
  and manage the execution of instructions inside the context of these frames,
- a `pygo.Frame` type to hold informations about the stack (globals, locals, 
  functions' code, ...),
- a `pygo.Block` type to handle the control flow (`if`, `else`, `return`, 
  `continue`, _etc_...),
- a `pygo.Instruction` type to model opcodes (`ADD`, `LOAD_FAST`, `PRINT`, ...)
  and their arguments (if any).

Ok.
That's enough for today.
Stay tuned...

In the meantime, I recommend reading the [AOSA article](http://www.aosabook.org/en/500L/a-python-interpreter-written-in-python.html).
