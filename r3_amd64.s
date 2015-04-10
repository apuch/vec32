
TEXT ·LengthR3(SB),7,$0-12
	MOVD        v+0(FP), BP
    MOVUPS      (BP), X0
    MULPS       X0, X0
    MOVLHPS     X0, X2
    MOVHLPS     X0, X2
    ADDPS       X2, X0
    MOVUPS      X0, X1
    SHUFPS      $0x03, X1, X1
    ADDSS       X1, X0
	SQRTSS      X0, X0
	MOVSS       X0, ret+8(FP)
    RET 
