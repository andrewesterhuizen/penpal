#include <midi>

__start:
    MOV A 0x80
    STORE A MIDI_ADDRESS_BPM
    MOV A 0x2
    STORE A MIDI_ADDRESS_PPQN
loop:
    JUMP loop

on_tick: 
    RAND
    MOV B 0x3C
    AND
    PUSH A
    PUSH 0x1
    CALL midi_trig
    RETI