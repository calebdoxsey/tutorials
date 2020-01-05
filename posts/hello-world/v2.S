bits 64
section .text
global _start
_start:
	mov rdi, [rsp + 16]
	call atoi

	mov r12, rax

.helloloop:
	mov rdx, len
	mov rsi, msg
	mov rdi, 1
	mov rax, 1
	syscall
	dec r12
	cmp r12, 0
	jg .helloloop

	mov rdi, 0
	mov rax, 60
	syscall


atoi:
    mov rax, 0              ; Set initial total to 0
.convert:
    movzx rsi, byte [rdi]   ; Get the current character
    test rsi, rsi           ; Check for \0
    je .done

    cmp rsi, 48             ; Anything less than 0 is invalid
    jl .error

    cmp rsi, 57             ; Anything greater than 9 is invalid
    jg .error

    sub rsi, 48             ; Convert from ASCII to decimal
    imul rax, 10            ; Multiply total by 10
    add rax, rsi            ; Add current digit to total

    inc rdi                 ; Get the address of the next character
    jmp .convert
.error:
    mov rax, -1             ; Return -1 on error
.done:
    ret                     ; Return total or error code

section .rodata
msg: db "hello world", 10
len: equ $-msg


