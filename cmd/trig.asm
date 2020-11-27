#include <midi>

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