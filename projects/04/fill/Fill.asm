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

// Put your code here.
@f
M=0
(INPUT)
    @KBD
    D=M
    @CLEAR
    D;JEQ
    @f
    D=M
    @INPUT
    D;JNE
    @f
    M=-1
    @FILL
    0;JMP
(CLEAR)
    @f
    D=M
    @INPUT
    D;JEQ
    @f
    M=0
    @FILL
    0;JMP
(FILL)
    @SCREEN
    D=A
    @i
    M=D
(LOOP)
    @i
    D=M
    @KBD
    D=D-A
    @INPUT
    D;JEQ
    @f
    D=M
    @i
    A=M
    M=D
    @i
    M=M+1
    @LOOP
    0;JMP
