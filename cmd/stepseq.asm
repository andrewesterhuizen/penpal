#include <midi>

i: db 0
length: db 16

notes:
    db 48
    db 49
    db 53
    db 56
    db 58
    db 60
    db 55
    db 65
    db 48
    db 51
    db 53
    db 56
    db 58
    db 51
    db 55
    db 53

steps:
    db 1
    db 1
    db 0
    db 1
    db 0
    db 1
    db 1
    db 0
    db 0
    db 1
    db 0
    db 1
    db 1
    db 1
    db 0
    db 1

start:
    mov A, 130
    store A, midi_bpm

    mov A, 4
    store A, midi_ppqn
loop:
    jump loop

inc_step:
    // increment index
    load i, A
    mov B, 1
    add
    store A, i

    // check if index == length
    load length, B
    gte
    jumpz inc_step_end

    // reset
    mov A, 0
    store A, i

    inc_step_end:
    ret

on_tick: 
    push 1
    call inc_step
    
    // check if step is active
    load i, A
    load (steps[A]), A
    jumpz skip

    // load note
    load i, A
    load (notes[A]), A
    push
    push 1
    call midi_trig

    skip:
    reti