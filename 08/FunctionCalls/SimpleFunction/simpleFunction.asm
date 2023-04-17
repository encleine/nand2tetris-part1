@256
D=A
@SP
M=D
// function SimpleFunction.test 2
(simpleFunction$SimpleFunction.test)
@2
D=A
@SP
M=D+M
// push local 0
@0
D=A
@LCL
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// push local 1
@1
D=A
@LCL
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M+D
@SP
M=M+1
// not
@SP
A=M-1
M=!M
// push argument 0
@0
D=A
@ARG 
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M+D
@SP
M=M+1
// push argument 1
@1
D=A
@ARG 
A=M+D
D=M
@SP
A=M
M=D
@SP
M=M+1
// sub
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M-D
@SP
M=M+1
// return
@return-prep
0;JMP
(return-prep)
// pop last in stack to arg
@SP
AM=M-1
D=M
@ARG
A=M
M=D
D=A
@R13
M=D
// pop that
@LCL
D=M
@SP
AM=D
D=M
@THAT
M=D
// pop this
@SP
AM=M-1
D=M
@THIS
M=D
// pop to arg
@SP
AM=M-1
D=M
@ARG
M=D
// pop to lcl
@SP
AM=M-1
D=M
@LCL
M=D
// return address
@SP
AM=M-1
D=M
@R15
M=D
@R13
D=M
@SP
M=D
@R15
A=M
0;JMP
(functionCallPrep)
// push return address 
@R13
D=M
@SP
A=M
M=D
@SP
M=M+1
// push lcl
@1
D=M
@SP
A=M
M=D
@SP
M=M+1
// push args
@2
D=M
@SP
A=M
M=D
@SP
M=M+1
// push this
@3
D=M
@SP
A=M
M=D
@SP
M=M+1
// push that
@4
D=M
@SP
A=M
M=D
@SP
M=M+1
//+ARG = (*sp) - 5 - arg num
@5
D=A
@R14
D=M+D
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
// R15 stores the call address 
@R15
A=M
0;JMP
