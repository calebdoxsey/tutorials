#!/bin/bash
go run -compiler gccgo -gccgoflags '-lgsl -lm -lgslcblas' main.go
