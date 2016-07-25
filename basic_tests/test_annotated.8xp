2a2a 5449 3833 462a 1a0a 00 //signature
47 656e 6572 6174 6564 2062 7920 7468 6520 5449 2d42 4153 4943 2043 6f6d 7069 6c65 722e 0000 0000 0000 00 comment
23 00 //length of data section

0d 00 //variable entry beginning, 0x00 padded on right
12 00 //len of variable data
05 //variable type ID (05 means edit-unlocked, not sure what the locking flag byte is)
5445 5354 0000 0000 //variable name "TEST"
00 //version
00 //archive flag
1200 //len of variable data
1000 //number of tokens present in variable data, 0x00 padded on right; this is also part of the variable data section
de2a 4845 4c4c 4f29 574f 524c 442d 2a3f //variable data
4906 //checksum