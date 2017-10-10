.PHONY: all build go-build cxx-build

ROOT_CXX_FLAGS := `root-config --cflags --libs` -lMinuit

build: go-build cxx-build

go-build: fba.go
	go build -o fba-gonum fba.go

cxx-build: fba.cxx
	$(CXX) -o fba-MINUIT fba.cxx $(ROOT_CXX_FLAGS)

