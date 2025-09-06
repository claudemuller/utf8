# Breaking Down UTF8

## Running

```bash
go run main.go
```

## Output

```bash
U+006f - o
[requires 1 bytes]
   byte 0 [0x6f - 01101111], masked [0x6f 01101111]

U+006b - k
[requires 1 bytes]
   byte 0 [0x6b - 01101011], masked [0x6b 01101011]

U+00a9 - ©
[requires 2 bytes]
   byte 0 [0xc2 - 11000010], masked [0x02 00000010]
   byte 1 [0xa9 - 10101001], masked [0x29 00101001] cont. byte

U+4e16 - 世
[requires 3 bytes]
   byte 0 [0xe4 - 11100100], masked [0x04 00000100]
   byte 1 [0xb8 - 10111000], masked [0x38 00111000] cont. byte
   byte 2 [0x96 - 10010110], masked [0x16 00010110] cont. byte

U+754c - 界
[requires 3 bytes]
   byte 0 [0xe7 - 11100111], masked [0x07 00000111]
   byte 1 [0x95 - 10010101], masked [0x15 00010101] cont. byte
   byte 2 [0x8c - 10001100], masked [0x0c 00001100] cont. byte

U+1f600 - 😀
[requires 4 bytes]
   byte 0 [0xf0 - 11110000], masked [0x00 00000000]
   byte 1 [0x9f - 10011111], masked [0x1f 00011111] cont. byte
   byte 2 [0x98 - 10011000], masked [0x18 00011000] cont. byte
   byte 3 [0x80 - 10000000], masked [0x00 00000000] cont. byte

```
