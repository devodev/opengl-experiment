import sys


if __name__ == '__main__':
    d = dict()
    for line in sys.stdin:
        idx, text = line.split(':')
        if 'TIMEOUT_IGNORED' in text:
            continue
        d.setdefault(int(idx, 16), []).append(text.strip())

    sortedD = list(sorted(d.items()))
    key_padding_size = len(str(sortedD[len(sortedD)-1][0])) + 1

    output = list()
    output.append('package opengl')
    output.append('')
    output.append('// GlEnum string lookup')
    output.append('var (')
    output.append('    GlEnums = map[uint32][]string{')
    for key, value in sortedD:
        map_key = '{}:'.format(key)
        map_value = ', '.join(map(lambda x: '"{}"'.format(x), value))
        output.append('        {:<{}} {{{}}},'.format(map_key, key_padding_size, map_value))
    output.append('    }')
    output.append(')')
    output.append('')

    print('\n'.join(output), end='')
