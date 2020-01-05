bits 64
section .text
global _start
_start:
	mov rdx, len
	mov rsi, msg
	mov rdi, 1
	mov rax, 1
	syscall

	mov rdi, 0
	mov rax, 60
	syscall

section .rodata
msg: db "hello world", 10
len: equ $-msg
