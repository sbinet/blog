package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot/vg"
)

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
}

func main() {
	costh := read("L3.dat")
	log.Printf("events: %d\n", len(costh))
	a, sig := asym(costh)
	log.Printf("A= %5.3f +/- %5.3f\n", a, sig)

	fcn := func(x []float64) float64 {
		lnL := 0.0
		A := x[0]
		const k = 3.0 / 8.0
		for _, v := range costh {
			lnL += math.Log(k*(1+v*v) + A*v)
		}
		return -lnL
	}

	grad := func(grad, x []float64) {
		fd.Gradient(grad, fcn, x, nil)
	}

	hess := func(h mat.MutableSymmetric, x []float64) {
		fd.Hessian(h.(*mat.SymDense), fcn, x, nil)
	}

	p := optimize.Problem{
		Func: fcn,
		Grad: grad,
		Hess: hess,
	}

	var meth = &optimize.Newton{}
	var p0 = []float64{0} // initial value for A

	res, err := optimize.Local(p, p0, nil, meth)
	if err != nil {
		log.Fatal(err)
	}
	display(res, p)
}

func read(fname string) []float64 {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var costh []float64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var pp, pm [3]float64
		_, err = fmt.Sscanf(
			sc.Text(), "%f %f %f %f %f %f",
			&pp[0], &pp[1], &pp[2],
			&pm[0], &pm[1], &pm[2],
		)
		if err != nil {
			log.Fatal(err)
		}
		v := pm[0]*pm[0] + pm[1]*pm[1] + pm[2]*pm[2]
		costh = append(
			costh,
			pm[2]/math.Sqrt(v),
		)
	}

	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}

	plot(costh)
	return costh
}

func asym(costh []float64) (float64, float64) {
	n := float64(len(costh))
	nf := 0
	for _, v := range costh {
		if v > 0 {
			nf++
		}
	}

	a := 2*float64(nf)/n - 1
	sig := math.Sqrt((1 - a*a) / n)
	return a, sig
}

func plot(vs []float64) {
	h := hbook.NewH1D(100, -1, 1)
	for _, v := range vs {
		h.Fill(v, 1)
	}

	p := hplot.New()
	p.Title.Text = "Cos Theta of mu-"
	hh := hplot.NewH1D(h)
	hh.FillColor = color.RGBA{0, 0, 255, 255}

	p.Add(hh)
	p.Add(hplot.NewGrid())

	err := p.Save(10*vg.Centimeter, -1, "plot.png")
	if err != nil {
		log.Fatal(err)
	}
}

func display(res *optimize.Result, p optimize.Problem) {
	fmt.Printf("results: %#v\n", res)

	errmat := mat.NewSymDense(len(res.X), nil)
	errmat.CopySym(res.Hessian)

	var ch mat.Cholesky
	if ok := ch.Factorize(errmat); !ok {
		log.Fatalf("could not factorize\n")
	}

	err := ch.InverseTo(errmat)
	if err != nil {
		log.Fatal(err)
	}
	for i, x := range errmat.RawSymmetric().Data {
		errmat.RawSymmetric().Data[i] = math.Sqrt(x)
	}

	fmt.Printf("minimal function value: %8.3f\n", res.F)
	// fmt.Printf("estimated diff to true minimum: %11.3e\n", edm(res, errmat))
	fmt.Printf("number of parameters: %d\n", len(res.X))

	fmt.Printf("grad=%v\n", res.Gradient)
	fmt.Printf("hess=%v\n", mat.Formatted(res.Hessian))

	fmt.Printf("errs= %v\n", mat.Formatted(errmat))

	for i, x := range res.X {
		xerr := errmat.At(i, i)
		fmt.Printf("par-%03d: %e +/- %e\n", i, x, xerr)
	}
}
