# CSC445 Programming Project 

Repo for development of a programming project for CSC445: OPERATIONS RESEARCH: LINEAR PROGRAMMING in Summer 2021.

## Usage
`go build`  
`./solver < [lp]`

## Implementation
This program is an implementation of the Revised Simplex Method using Bland's Rule as defined in the lecture slides.

I have three packages utils, lp and simplex.

lp contains the struct for the LP. So a constructor and the various methods for accessing subsets of the matrices and vectors.

utils contains functions for various things, mostly indexing matrices and vectors with an array of ints.

simplex contains both the primal and dual procedures for the solver.

## Tests
`sh test.sh`

At time of submission:
- All the vanderbei tests are passing except vanderbei_example14.1.txt. 
- Three netlib tests are passing (netlib_bgprtr.txt, netlib_itest2.txt and netlib_itest6.txt), the others cycle or don't return a correct result. 
- All the volume 2 tests are passing except optimal_3x3_6.txt due to some floating point error.
