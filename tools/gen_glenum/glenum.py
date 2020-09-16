import sys


if __name__ == '__main__':
    d = dict()
    for line in sys.stdin:
        idx, text = line.split(':')
        d.setdefault(int(idx, 16), []).append(text.strip())

    sortedD = list(sorted(d.items()))
    key_padding_size = len(str(sortedD[len(sortedD)-1][0]))

    comment = ''
    print('package opengl')
    print()
    print('// GlEnum string lookup')
    print('var (')
    print('    GlEnums = map[int][]string{')
    for key, value in sortedD:
        print('        ', end='')
        if 'TIMEOUT_IGNORED' in value:
            print('// ', end='')
        print('{:<{}} []string{{{}}},'.format('{}:'.format(key), key_padding_size+1, ', '.join(map(lambda x: '"{}"'.format(x), value))))
    print('    }')
    print(')')
    print()
