+++
title = "Go-HEP Manifesto"
date = 2018-07-31T10:22:04+02:00
Tags = []
Categories = ["go", "go-hep", "HEP"]
+++

Hello again.

I am starting today an article for arXiv about [Go](https://golang.org) and [Go-HEP](https://go-hep.org).
I thought structuring my thoughts a bit (in the form of a blog post) would help fluidify the process.

## (HEP) Software is painful

In my introduction talk(s) about [Go](https://golang.org) and [Go-HEP](https://go-hep.org), such as [here](https://talks.godoc.org/github.com/sbinet/talks/2018/2018-04-03-go-hep-jlab/talk.slide#3), I usually talk about software being painful.
HENP software is no exception.
It *is* painful.

As a C++/Python developer and former software architect of one of the four LHC experiments, I can tell you from vivid experience that software is painful to develop.
One has to tame deep and complex software stacks with huge dependency lists.
Each dependency comes with its own way to be configured, built and installed.
Each dependency comes with its own dependencies.
When you start working with one of these software stacks, installing them on your own machine is no walk in the park, even for experienced developers.
These software stacks are real snowflakes: they need their unique cocktail of dependencies, with the right version, compiler toolchain and OS, tightly integrated on usually a single development platform.

Granted, the _de facto_ standardization on `CMake` and `docker` did help with some of these aspects, allowing projects to cleanly encapsulate the list of dependencies in a reproducible way, in a container.
Alas, this renders code easier to deploy but less portable: everything is `linux/amd64` plus some arbitrary Linux distribution.

In HENP, with `C++` being now the _lingua franca_ for everything that is related with framework or infrastructure, we get unwiedly compilation times and thus a very unpleasant edit-compile-run development cycle.
Because `C++` is a very complex language to learn, read and write - each new revision more complex than the previous one - it is becoming harder to bring new people on board with existing `C++` projects that have accumulated a lot of technical debt over the years: there are many layers of accumulated cruft, different styles, different ways to do things, etc...

Also, HENP projects heavily rely on shared libraries: not because of security, not because they are faster at runtime (they are not), but because as `C++` is so slow to compile, it is more convenient to not recompile everything into a static binary.
And thus, we have to devise sophisticated deployment scenarii to deal with all these shared libraries, properly configuring `$LD_LIBRARY_PATH`, `$DYLD_LIBRARY_PATH` or `-rpath`, adding yet another moving piece in the machinery.
We did not have to do that in the `FORTRAN` days: we were building static binaries.

From a user perspective, HENP software is also - even more so - painful.
One needs to deal with:

- overly complicated Object Oriented systems,
- overly complicated inheritance hierarchies,
- overly complicated meta-template programming,

and, of course, dependencies.
It's 2018 and there are still no simple way to handle dependencies, nor a standard one that would work across operating systems, experiments or analysis groups, when one lives in a `C++` world.
Finally, there is no standard way to retrieve documentation - and here we are just talking about APIs - nor a system that works across projects and across dependencies.

All of these issues might explain why many physicists are migrating to `Python`.
The ecosystem is much more integrated and standardized with regard to installation procedures, serving, fetching and describing dependencies and documentation tools.
`Python` is also simpler to learn, teach, write and read than `C++`.
But it is also slower.

Most physicists and analysts are willing to pay that price, trading reduced runtime efficiency for a wealth of scientific, turn-key pure-Python tools and libraries.
Other physicists strike a different compromise and are willing to trade the relatively seamless installation procedures of pure-Python software with some runtime efficiency by wrapping `C/C++` libraries.

To summarize, `Python` and `C++` are no _panacea_ when you take into account the vast diversity of programming skills in HENP, the distributed nature of scientific code development in HENP, the many different teams' sizes and the constraints coming from the development of scientific analyses (agility, fast edit-compile-run cycles, reproducibility, deployment, portability, ...)
To add insult to injury, these languages are rather ill equiped to cope with distributed programming and parallel programming: either because of a technical limitation (`CPython`'s Global Interpreter Lock) or because the current toolbox is too low-level or error-prone.

Are we really left with either:

- a language that is relatively fast to develop with, but slow at runtime, or
- a language that is painful to develop with but fast at runtime ?

![nogo](/code/2018-07-31/funfast-nogo.svg)

## Mending software with Go

Of course, I think [Go](https://golang.org) can greatly help with the general situation of software in HENP.
It is not a magic wand, you still have to think and apply work.
But it is a definitive, positive improvement.

![go-logo](/code/2018-07-31/golang-logo.png?bounds=10x20)

Go was created to tackle all the challenges that `C++` and `Python` couldn't overcome.
Go was designed for ["programming in the large"](https://talks.golang.org/2012/splash.article).
Go was designed to strive at scales: software development at Google-like scale but also at 2-3 people scale.

But, most importantly, Go wasn't designed to be a good programming language, it was designed for [*software engineering*](https://research.swtch.com/vgo-eng):

```
  Software engineering is what happens to programming 
  when you add time and other programmers.
```

Go is a simple language - not a simplistic language - so one can easily learn most of it in a couple of days and be proficient with it in a few weeks.

Go has builtin tools for concurrency (the famed `goroutines` and `channels`) and that is what made me try it initially.
But I stayed with Go for everything else, _ie_ the tooling that enables:

- code refactoring with `gorename` and `eg`,
- code maintenance with `goimports`, `gofmt` and `go fix`,
- code discoverability and completion with `gocode`,
- local documentation (`go doc`) and across projects ([godoc.org](https://godoc.org)),
- integrated, *simple*, build system (`go build`) that handles dependencies (`go get`), without messing around with `CMakeList.txt`, `Makefile`, `setup.py` nor `pom.xml` build files: all the needed information is in the source files,
- easiest cross-compiling toolchain to date.

And all these tools are usable from every single editor or IDE.

Go compiles optimized code really quickly.
So much so that the `go run foo.go` command, that compiles a complete program and executes it on the fly, feels like running `python foo.py` - but with builtin concurrency and better runtime performances (CPU and memory.)
Go produces static binaries that usually do not even require `libc`.
One can take a binary compiled for `linux/amd64`, copy it on a Centos-7 machine or on a Debian-8 one, and it will happily perform the requested task.

As a _Gedankexperiment_, take a standard `centos7` `docker` image from docker-hub and imagine having to build your entire experiment software stack, from the exact gcc version down to the last wagon of your train analysis.

- How much time would it take?
- How much effort of tracking dependencies and ensuring internal consistency would it take?
- How much effort would it be to deploy the binary results on another machine? on another non-Linux machine?

Now consider this script:

```sh
#!/bin/bash

yum install -y git mercurial curl

mkdir /build
cd /build

## install the Go toolchain
curl -O -L https://golang.org/dl/go1.10.3.linux-amd64.tar.gz
tar zxf go1.10.3.linux-amd64.tar.gz
export GOROOT=`pwd`/go
export GOPATH=/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

## install Go-HEP and its dependencies
go get -v go-hep.org/x/hep/...
```

Running this script inside said container yields:

```sh
$> time ./install.sh
[...]
go-hep.org/x/hep/xrootd/cmd/xrd-ls
go-hep.org/x/hep/xrootd/server
go-hep.org/x/hep/xrootd/cmd/xrd-srv

real  2m30.389s
user  1m09.034s
sys   0m14.015s
```

In less than 3 minutes, we have built a container with (almost) all the tools to perform a HENP analysis.
The bulk of these 3 minutes is spent cloning repositories.

Building [root-dump](https://godoc.org/go-hep.org/x/hep/rootio/cmd/root-dump), a program to display the contents of a [ROOT](https://root.cern) file for, say, Windows, can easily performed in one single command:

```sh
$> GOOS=windows \
   go build go-hep.org/x/hep/rootio/cmd/root-dump
$> file root-dump.exe 
root-dump.exe: PE32+ executable (console) x86-64 (stripped to external PDB), for MS Windows

## now, for windows-32b
$> GOARCH=386 GOOS=windows \
   go build go-hep.org/x/hep/rootio/cmd/root-dump
$> file root-dump.exe 
root-dump.exe: PE32 executable (console) Intel 80386 (stripped to external PDB), for MS Windows
```

Fun fact: [Go-HEP](https://go-hep.org) was supporting Windows users wanting to read ROOT-6 files *before* ROOT itself (ROOT-6 support for Windows landed with `6.14/00`.)

## Go & Science

Most of the needed scientific tools are available in Go at [gonum.org](https://gonum.org):

- plots,
- network graphs,
- integration,
- statistical analysis,
- linear algebra,
- optimization,
- numerical differentiation,
- probability functions (univariate and multivariate),
- discrete Fourier transforms

Gonum is almost at feature parity with the `numpy/scipy` stack.
Gonum is still missing some tools, like [ODE](https://en.wikipedia.org/wiki/Ordinary_differential_equation) or more interpolation tools, but the chasm is closing.

Right now, in a HENP context, it is not possible to perform an analysis in Go and insert it in an already existing C++/Python pipeline.
At least not easily: while reading is possible, [Go-HEP](https://go-hep.org) is still missing the ability to write ROOT files.
This restriction should be lifted before the end of 2018.

That said, Go can already be quite useful and usable, now, in science and HENP, for data acquisition, monitoring, cloud computing, control frameworks and some physics analyses.
Indeed, [Go-HEP](https://go-hep.org) provides HEP-oriented tools such as histograms and n-tuples, Lorentz vectors, fitting, interoperability with [HepMC](https://gitlab.cern.ch/hepmc/HepMC) and other Monte-Carlo programs (HepPDT, LHEF, SLHA), a toolkit for a fast detector simulation à la [Delphes](https://cp3.irmp.ucl.ac.be/projects/delphes) and libraries to interact with [ROOT](https://root.cern) and [XRootD](http://xrootd.org).

I think building the missing scientific libraries in Go is a better investment than trying to fix the `C++/Python` languages and ecosystems.

Go is a better trade-off for software engineering and for science:

![with-go](/code/2018-07-31/funfast.svg)

---

**PS:** There's a nice discussion about this post on the [Go-HEP forum](https://groups.google.com/forum/#!topic/go-hep/H-_Mj1JKeT4).
