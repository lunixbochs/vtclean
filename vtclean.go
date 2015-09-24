package vtclean

import (
	"regexp"
	"strconv"
)

// see regex.txt for a slightly separated version of this regex
var vt100re = regexp.MustCompile(`^\033(([A-KZ=>12<]|Y\d{2})|\[\d+[A-D]|\[\d+;\d+[Hf]|#[1-68]|\[(\d+|;)*[qm]|\[[KJg]|\[[0-2]K|\[[02]J|\([ABCEHKQRYZ0-7=]|[\[K]\d+;\d+r|\[[03]g|\[\?[1-9][lh]|\[20[lh]|\[[56]n|\[0?c|\[2;[1248]y|\[!p|\[([01457]|254)}|\[\?(12;)?(25|50)[lh]|[78DEHM]|\[[ABCDHJKLMP]|\[4[hl]|\[\?1[46][hl]|\[\*[LMP]|\[[12][JK]|\]\d*;\d*[^\x07]+\x07|\[\d*[@ABCDEFGIJKLMPSTXZ1abcdeghilmnp])`)
var vt100color = regexp.MustCompile(`^\033\[(\d+|;)*[m]`)
var lineEdit = regexp.MustCompile(`^\033\[(\d*)([CDPK])`)

func vt100scan(line string) int {
	return len(vt100re.FindString(line))
}

func isColor(line string) bool {
	return len(vt100color.FindString(line)) > 0
}

func Clean(line string, color bool) string {
	out := make([]rune, len(line))
	liner := []rune(line)
	hadColor := false
	pos, max := 0, 0
	for i := 0; i < len(liner); {
		c := liner[i]
		str := string(liner[i:])
		switch c {
		case '\b':
			pos -= 1
		case '\x7f':
			copy(out[pos:max], out[pos+1:max])
			max -= 1
		case '\033':
			if m := lineEdit.FindStringSubmatch(str); m != nil {
				i += len(lineEdit.FindString(str))
				n, err := strconv.Atoi(m[1])
				if err != nil || n > 10000 {
					n = 1
				}
				switch m[2] {
				case "C":
					pos += n
				case "D":
					pos -= n
				case "P":
					most := max - pos
					if n > most {
						n = most
					}
					copy(out[pos:], out[pos+n:])
					max -= n
				case "K":
					max = pos
				}
				if pos < 0 {
					pos = 0
				}
				if pos > max {
					pos = max
				}
				continue
			}
			if !(color && isColor(str)) {
				skip := vt100scan(str)
				if skip > 0 {
					i += skip
					continue
				}
			} else {
				hadColor = true
				out[pos] = c
				pos++
			}
		default:
			if c == '\n' || c >= ' ' {
				out[pos] = c
				pos++
			}
		}
		if pos > max {
			max = pos
		}
		i += 1
	}
	out = out[:max]
	if hadColor {
		out = append(out, []rune("\033[0m")...)
	}
	return string(out)
}
