Input "A? ",A
Input "ROOT? ", R
If A^(1/R)=int(A^(1/R))
Then
Disp A^(1/R)
End
1→Z
Goto B
Lbl B
If A>0
Then
2→H
1→L
End
If A<0
Then
-2→H
-1→L
End
For(X,H,int(A/2),L)
X^(1/R)→Q
If (A/X)=int(A/X) and Q=int(Q)
Then
Q*Z→Z
(A/X)→A
Goto B
End
End
End
Disp "A*B^(1/R)"
Disp "A:",Z
Disp "B:",A