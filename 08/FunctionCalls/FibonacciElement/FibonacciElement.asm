@256
D=A
@SP
M=D
// function Main.fibonacci 0
(FibonacciElement$Main.fibonacci)
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
// push constant 2
@2
D=A
@SP
A=M
M=D
@SP
M=M+1
// lt
@SP
AM=M-1
D=M
@SP
AM=M-1
D=M-D
M=-1
@skip$0
D;JLT
@SP
A=M
M=0
(skip$0)
@SP
M=M+1
// if-goto IF_TRUE
@SP
AM=M-1
D=M
@FibonacciElement$Main.fibonacci$IF_TRUE
D;JNE
// goto IF_FALSE
@FibonacciElement$Main.fibonacci$IF_FALSE
0;JMP
// label IF_TRUE
(FibonacciElement$Main.fibonacci$IF_TRUE)
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
// return
@return-prep
0;JMP
// label IF_FALSE
(FibonacciElement$IF_FALSE)
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
// push constant 2
@2
D=A
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
// call Main.fibonacci 1
@address$Main.fibonacci$0
D=A
@R13
M=D
@FibonacciElement$Main.fibonacci
D=A
@R15
M=D
@1
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Main.fibonacci$0)
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
// push constant 1
@1
D=A
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
// call Main.fibonacci 1
@address$Main.fibonacci$1
D=A
@R13
M=D
@FibonacciElement$Main.fibonacci
D=A
@R15
M=D
@1
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Main.fibonacci$1)
// add
@SP
AM=M-1
D=M
@SP
AM=M-1
M=M+D
@SP
M=M+1
// return
@return-prep
0;JMP
// function Sys.init 0
(FibonacciElement$Sys.init)
// push constant 4
@4
D=A
@SP
A=M
M=D
@SP
M=M+1
// call Main.fibonacci 1
@address$Main.fibonacci$0
D=A
@R13
M=D
@FibonacciElement$Main.fibonacci
D=A
@R15
M=D
@1
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Main.fibonacci$0)
// label WHILE
(FibonacciElement$Sys.init$WHILE)
// goto WHILE
@FibonacciElement$Sys.init$WHILE
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
