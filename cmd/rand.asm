#include <midi>

start:
loop:
    jump loop

on_tick: 
    rand
    mov B, 0x3C
    and
    push
    push 0x1
    call midi_trig
    reti