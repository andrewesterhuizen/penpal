#include <midi>

__start:
    PUSH 0x30
    PUSH 0x1
    CALL trig
    HALT
