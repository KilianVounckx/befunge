#BEFUNGE

A [befunge](https://esolangs.org/wiki/Befunge) interpreter written in [go](www.golang.org).

The interpreter can run befunge files.

###How To Use?
Create a file with a befunge program in it (e.g. [factorial.bf](factorial.bf), this program asks
the user for a number and calculates the factorial). Then run:
```shell
>> befunge factorial.bf
```

###Debug Mode
There is a basic debug mode. Instead of going through the whole program at once, the interpreter
pauses at every step. It then prints the position of the PC, the character under the PC and the
current stack. Enable debug mode by adding the '-debug' flag before the file name.
