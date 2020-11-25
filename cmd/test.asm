send_midi:
    // status
    STORE +5(fp) 0x0 
    // data1
    STORE +6(fp) 0x1 
    // data2
    STORE +7(fp) 0x2

    SEND
    RET

trig:
    // note on
    PUSH 0x63
    PUSH +5(fp)
    PUSH 0x90
    PUSH 0x3
    CALL send_midi

    // note off
    PUSH 0x63
    PUSH +5(fp)
    PUSH 0x80
    PUSH 0x3
    CALL send_midi

    RET

__start:
    PUSH 0x30
    PUSH 0x1
    CALL trig
    HALT
