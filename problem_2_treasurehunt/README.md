# Task 2: Treasure Hunt

For the solution, I am assuming the coordinate ```(0,0)``` is on the top left corner and ```(W-1, H-1)``` is on bottom right corner where ```W``` and ```H``` is the width and height of the layout measure in number of characters.

To execute the program:
```
python3 solution.py
```

Input the layout matrix as it is and followed by the values of A, B and C consecutively separated by single space in a new line. See example bellow

Example:
```
$ python3 solution.py
########
#......#
#.###..#
#...#.##
#X#....#
########
3 4 3
```

The program will output like this:

```
########
#$$$$$.#
#$###$.#
#$..#$##
#X#..$.#
########

(3, 1)
(2, 1)
(1, 1)
(1, 2)
(1, 3)
(1, 4)
(1, 5)
(2, 5)
(3, 5)
(4, 5)

```