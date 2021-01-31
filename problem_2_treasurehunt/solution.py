

def findStart(fields):
    for x in range(len(fields)):
        for y in range(len(fields[x])):
            if fields[x][y] == 'X':
                return x,y


def drawFields(field, a , b, c):
    result = []
    xs, ys = findStart(field)
    # print (xs, ys)
    x,y = xs, ys

    # move up/north a steps
    for i in range(a):
        x -= 1
        # print (field[x][y])
        if(field[x][y] == '#'):
            break
        if(field[x][y] == '.'):
            result += [(x, y)]
            field[x][y] = '$'

    # move right/east b steps
    for i in range(b):
        y += 1
        if(field[x][y] == '#'):
            break
        if(field[x][y] == '.'):
            result += [(x, y)]
            field[x][y] = '$'

    # move down/south c steps
    for i in range(c):
        x += 1
        if(field[x][y] == '#'):
            break
        if(field[x][y] == '.'):
            result += [(x, y)]
            field[x][y] = '$'

    for l in field:
        print("".join(l))

    print()

    for t in result:
        print(t)
    

def main():
    
    fields = [
        "########",
        "#......#",
        "#.###..#",
        "#...#.##",
        "#X#....#",
        "########"
    ]

    params = input("")
    a, b, c = params.split(" ")
    arrField = [list(a) for a in fields]

    drawFields(arrField, int(a), int(b), int(c))    

if __name__ == "__main__":
    main()