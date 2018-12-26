# dallas
[![Build Status](https://travis-ci.org/mitchr/dallas.svg?branch=master)](https://travis-ci.org/mitchr/dallas)

## About
dallas is a TI-BASIC compiler. It supports the TI-83 and TI-83+/TI-84+.

## Installation
`go get github.com/Mitchell-Riley/dallas`

## Usage
`dallas [flags] filename`

Flag|Type|Description
----|----|----
-a|bool|set the archive bit; if false, ram is used to store the program
-d|bool|set to disassemble .8xp files
-e|bool|set the edit-lock bit
-h|bool|display this help message
-ti83|bool|compile for the TI-83
-p|string|set the program name (default "PROG")

Some tokens had to be changed from unicode to ascii, so that you can actually type them:

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
-(neg)|^-
√(|sqrt(
³√(|crt(


### Sources
[TI-83+/TI-84+ file format guide](http://merthsoft.com/linkguide/ti83+/fformat.html)
[Token List](http://tibasicdev.wikidot.com/one-byte-tokens)
