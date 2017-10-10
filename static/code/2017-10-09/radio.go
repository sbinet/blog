package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/vg"
)

func main() {
	rs := []float64{1, 1, 5, 4, 2, 0, 3, 2, 4, 1, 2, 1, 1, 0, 1, 1, 2, 1}
	rfac := []float64{1, 1, 120, 24, 2, 1, 6, 2, 24, 1, 2, 1, 1, 1, 1, 1, 2, 1}

	mean := stat.Mean(rs, nil)
	merr := math.Sqrt(mean / float64(len(rs)))

	fmt.Printf("mean=%v\n", mean)
	fmt.Printf("merr=%v\n", merr)

	fcn := func(x []float64) float64 {
		mu := x[0]
		lnl := 0.0
		for i := range rs {
			lnl += rs[i]*math.Log(mu) - mu - math.Log(rfac[i])
		}
		return -lnl
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
	var p0 = []float64{1} // initial value for mu

	res, err := optimize.Local(p, p0, nil, method)
	if err != nil {
		log.Fatal(err)
	}

	display(res, p)
	plot(rs)
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
		xerr := errmat.At(i, i) * upValue
		fmt.Printf("par-%03d: %e +/- %e\n", i, x, xerr)
	}
}

const upValue = 2 // 2 for Log-likelihood fits, 1 for Chi-2

func edm(res *optimize.Result, errs mat.Symmetric) float64 {
	edm := 0.0
	for i := range res.X {
		grad := res.Gradient[i]
		// off-diagonal elements
		for j := 0; j < i; j++ {
			z := 2 * errs.At(i, j)
			edm += grad * z * res.Gradient[j]
		}
		// diagonal elements
		edm += errs.At(i, i) * grad * grad
	}
	return edm
}

func plot(rs []float64) {
	// hbook is go-hep.org/x/hep/hbook.
	// here we create a 1-dim histogram with 10 bins,
	// from 0 to 10.
	h := hbook.NewH1D(10, 0, 10)
	for _, x := range rs {
		h.Fill(x, 1)
	}

	// hplot is a convenience package built on top
	// of gonum.org/v1/plot.
	// hplot is go-hep.org/x/hep/hplot.
	p := hplot.New()
	p.X.Label.Text = "r"
	hh := hplot.NewH1D(h)
	hh.FillColor = color.RGBA{255, 0, 0, 255}
	p.Add(hh)
	p.Add(hplot.NewGrid())

	err := p.Save(10*vg.Centimeter, -1, "plot.png")
	if err != nil {
		log.Fatal(err)
	}
}
