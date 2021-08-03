# CSC445 Programming Project 

This application is a simple linear program solver built for CSC445: OPERATIONS RESEARCH: LINEAR PROGRAMMING in Summer 2021.

## Usage
`rm ._main.go && go build`  
`./solver < [lp file]`

## Example input
```
13       12        9
0.5      0.4       0.4     10
0.3      0         0       5
0.1      0.2       0.4     10
0        0.3       0       1
0        0.1       0.2     2
```
gives: 
```
optimal
264.1667
16.66667 3.333333 0.8333333
```

## Implementation
This program is an implementation of the Revised Simplex Method using Bland's Rule for pivoting.

I have three packages `utils`, `lp` and `simplex`.

`lp` contains the struct for the LP. So a constructor and the various methods for accessing subsets of the matrices and vectors.

`utils` contains functions for various things, mostly indexing matrices and vectors with an array of ints.

`simplex` contains both the primal and dual procedures for the solver.

## Tests
`sh test.sh`

## To Do
- Solve linear equations instead of doing inverse. DONE
    - except the dzn calculation https://github.com/iamkahvi/CSC445_project/blob/c4dfc95df9d7a72c5a4c3e4d850b67c2a085b016/lp/lp.go#L206
- Solve test cases `optimal_3x3_6.txt`, `netlib_share2b.txt`, `netlib_share1b.txt`, `netlib_adlittle.txt` and `netlib_afiro.txt`.
- Correct unwanted cycling on `netlib_klein2.txt`.
    - Perturbation method https://github.com/iamkahvi/CSC445_project/blob/eabc4d09c3d6513eedbfa5ff9f6f21ecdbc0ce46/main.go#L32
