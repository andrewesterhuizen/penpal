#include <midi>

i: db 0

notes:
    db 48
    db 51
    db 53
    db 56
    db 58
    db 60
    db 65

__start:
    mov A, 4
    store A, midi_ppqn
loop:
    jump loop


next_note:
    // increment index
    load i, A
    mov B, 1
    add
    store A, i

    // check if index == 4
    mov B, 6
    eq
    jumpz after_reset

    // if 4, reset i to 0
    mov A, 0
    store A, i

    after_reset:
    load i, A
    load (notes[A]), A

    ret

on_tick: 
    call next_note
    push
    push 0x1
    call midi_trig

    reti