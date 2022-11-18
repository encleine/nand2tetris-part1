// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.
@R5 
D = M
@color
M = D
(reset)
	@SCREEN
	D = A
	@address
	M = D
//? checking if a key is pressed
(keyboard)
	@KBD
	D = M
	@color
	M = D
	@loop
	D;JEQ
	@color
	M = -1
//* changes the screen color
(loop)
	@color
	D = M
	@address
	A = M
	M = D

	@24576 
	D = A
	@address
	D = D - M 
	M = M + 1
		@keyboard
		D;JGT 
	@reset
	0;JMP