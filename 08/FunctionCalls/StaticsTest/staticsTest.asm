@256
D=A
@SP
M=D
// function Class1.set 0
(staticsTest$Class1.set)
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
// pop static 0
@SP
AM=M-1
D=M
@staticsTest$0
M=D
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
// pop static 1
@SP
AM=M-1
D=M
@staticsTest$1
M=D
// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
// return
@return-prep
0;JMP
// function Class1.get 0
(staticsTest$Class1.get)
// push static 0
@staticsTest$0
D=M
@SP
A=M
M=D
@SP
M=M+1
// push static 1
@staticsTest$1
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
// function Class2.set 0
(staticsTest$Class2.set)
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
// pop static 0
@SP
AM=M-1
D=M
@staticsTest$0
M=D
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
// pop static 1
@SP
AM=M-1
D=M
@staticsTest$1
M=D
// push constant 0
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
// return
@return-prep
0;JMP
// function Class2.get 0
(staticsTest$Class2.get)
// push static 0
@staticsTest$0
D=M
@SP
A=M
M=D
@SP
M=M+1
// push static 1
@staticsTest$1
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
// function Sys.init 0
(staticsTest$Sys.init)
// push constant 6
@6
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 8
@8
D=A
@SP
A=M
M=D
@SP
M=M+1
// call Class1.set 2
@address$Class1.set$0
D=A
@R13
M=D
@staticsTest$Class1.set
D=A
@R15
M=D
@2
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Class1.set$0)
// pop temp 0
@SP
AM=M-1
D=M
@5
M=D
// push constant 23
@23
D=A
@SP
A=M
M=D
@SP
M=M+1
// push constant 15
@15
D=A
@SP
A=M
M=D
@SP
M=M+1
// call Class2.set 2
@address$Class2.set$1
D=A
@R13
M=D
@staticsTest$Class2.set
D=A
@R15
M=D
@2
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Class2.set$1)
// pop temp 0
@SP
AM=M-1
D=M
@5
M=D
// call Class1.get 0
@address$Class1.get$2
D=A
@R13
M=D
@staticsTest$Class1.get
D=A
@R15
M=D
@0
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Class1.get$2)
// call Class2.get 0
@address$Class2.get$3
D=A
@R13
M=D
@staticsTest$Class2.get
D=A
@R15
M=D
@0
D=A
@R14
M=D
@functionCallPrep
0;JMP
(address$Class2.get$3)
// label WHILE
(staticsTest$Sys.init$WHILE)
// goto WHILE
@staticsTest$Sys.init$WHILE
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
