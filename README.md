# !dices
## Overview
A simple Dice bot for discord written in go made by me :)  

## Usage
In the discord chat we have to call the command !dices for interacting with the bot. The idea here is to use the patterns already known for commands in discord, basically adding the "!" after the command itself.  

Rolling one d20:
```
!dices d20
```
or:
```
!dices 1d20
```
output:
```
d20: (10) sum: 10 med: 10
```
Running multiple dices of the same type (d10):  
```
!dices 3d10
```
output:
```
d10: (7) (8) (9) sum: 24 med: 8
```
