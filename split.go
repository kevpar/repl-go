package repl

import "errors"

func split(s string) ([]string, error) {
	errBadString := errors.New("invalid string")
	const (
		stateNormal = iota
		stateSlash
		stateHexStart
		stateHexMid
	)
	var (
		elems   []string
		prev    int
		state   int
		partial string
		hex     rune
	)
	for i, c := range s {
		switch state {
		case stateHexStart:
			switch c {
			case '{':
				state = stateHexMid
			default:
				return nil, errBadString
			}
		case stateHexMid:
			if c >= '0' && c <= '9' {
				if hex > 0xfffffff {
					return nil, errBadString
				} else if hex != 0 {
					hex <<= 4
				}
				hex |= c - '0'
			} else if c >= 'A' && c <= 'F' {
				if hex > 0xfffffff {
					return nil, errBadString
				} else if hex != 0 {
					hex <<= 4
				}
				hex |= c - 'A' + 10
			} else if c >= 'a' && c <= 'f' {
				if hex > 0xfffffff {
					return nil, errBadString
				} else if hex != 0 {
					hex <<= 4
				}
				hex |= c - 'a' + 10
			} else if c == '}' {
				partial += string(hex)
				hex = 0
				prev = i + 1
				state = stateNormal
			} else {
				return nil, errBadString
			}
		case stateSlash:
			switch c {
			case ' ':
				partial += " "
				prev = i + 1
				state = stateNormal
			case '\\':
				partial += `\`
				prev = i + 1
				state = stateNormal
			case 'x':
				state = stateHexStart
			default:
				return nil, errBadString
			}
		case stateNormal:
			switch c {
			case '\\':
				partial += s[prev:i]
				state = stateSlash
			case ' ':
				partial += s[prev:i]
				if partial != "" {
					elems = append(elems, partial)
				}
				partial = ""
				prev = i + 1
			}
		}
	}
	if state != stateNormal {
		return nil, errBadString
	}
	partial += s[prev:]
	if partial != "" {
		elems = append(elems, partial)
	}
	return elems, nil
}
