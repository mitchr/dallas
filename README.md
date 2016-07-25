# dallas

## About
dallas is a TI-BASIC compiler.

## Installation
`go get github.com/Mitchell-Riley/dallas`

## Usage
`dallas [flags] filename`

Flag|Type|Description
----|----|----
-a|bool|set the archive bit; if false, ram is used to store the program
-d|bool|set to true to disassemble .8xp files
-e|bool|set the edit-lock bit
-h|bool|display this help message
-o|string|set the name of the output file (default "PROG.8xp")
-p|string|set the program name (default "PROG")

Some tokens had to be changed from their unicode format to an ascii format:

Calculator|Dallas
---|---
r|radian
°|degree
ֿ¹|inverse
²|^2
T (transpose)|transpose
³|^3
→|->
θ|theta
≤|<=
≥|>=
≠|!=
π|pi
- (neg)|neg
√(|sqrt(
³√(|crt(
