here:
    db 0xaa
    db 0xcc

__start:
    mov A, 2
    load (here[A]), B
    halt  

