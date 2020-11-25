send_midi:
    // status
    MOV A +5(fp) 
    STORE 0x0 

    // data1
    MOV A +6(fp) 
    STORE 0x1 

    // data2
    MOV A +7(fp) 
    STORE 0x2

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
    PUSH 0x60
    PUSH 0x1
    CALL trig
    HALT
