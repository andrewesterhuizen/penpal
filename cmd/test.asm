#include <midi>

__start:
    mov A 0x80
    store A MIDI_ADDRESS_BPM
    mov A 0x2
    store A MIDI_ADDRESS_PPQN
loop:
    jump loop

on_tick: 
    rand
    mov B 0x3C
    and
    push A
    push 0x1
    call midi_trig
    reti