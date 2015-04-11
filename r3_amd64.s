
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

TEXT ·Add3(SB),7,$0-24
	MOVD	v1+0(FP),BP
	MOVUPS	(BP), X0
	MOVD	v2+8(FP), BP
	ADDPS	(BP), X0
	MOVD	v3+16(FP), BP
	MOVUPS	X0, (BP)
	RET

TEXT ·OrthoBoxAdd(SB),7,$0-16
	MOVD	bb2+8(FP), 	BP
	MOVUPS  (BP),      	X0
	MOVUPS  16(BP),   	X1
	MOVD	bb1+0(FP), 	BP
	MINPS	(BP),       X0
	MAXPS	16(BP),    	X1
	RET

TEXT ·doNop(SB),0,$0-0
	RET

