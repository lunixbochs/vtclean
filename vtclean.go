package vtclean

import (
	"bytes"
	"regexp"
	"strconv"
)

// see regex.txt for a slightly separated version of this regex
var vt100re = regexp.MustCompile(`^\033(([A-KZ=>12<]|Y\d{2})|\[\d+[A-D]|\[\d+;\d+[Hf]|#[1-68]|\[(\d+|;)*[qm]|\[[KJg]|\[[0-2]K|\[[02]J|\([ABCEHKQRYZ0-7=]|[\[K]\d+;\d+r|\[[03]g|\[\?[1-9][hl]|\[\?1[4689][hl]|\[1[289][hl]|\[20[hl]|\[[024][hl]|\[[56]n|\[0?c|\[2;[1248]y|\[!p|\[([01457]|254)}|\[\?(12;)?(25|50)[lh]|[78DEHM]|\[[ABCDHJKLMP]|\[\*[LMP]|\[[12][JK]|\]\d*;\d*[^\x07]+\x07|\[\d*[@ABCDEFGIJKLMPSTXZ1abcdeghilmnp])`)
var vt100color = regexp.MustCompile(`^\033\[(\d+|;)*[m]`)
var lineEditRe = regexp.MustCompile(`^\033\[(\d*)([@CDPK])`)

func vt100scan(line []byte) int {
	return len(vt100re.Find(line))
}

func isColor(line []byte) bool {
	return len(vt100color.Find(line)) > 0
}

func Clean(line string, color bool) string {
	var edit = lineEdit{buf: make([]byte, len(line))}
	lineb := []byte(line)

	hadColor := false
	for i := 0; i < len(lineb); {
		c := lineb[i]
		switch c {
		case '\b':
			edit.Move(-1)
		case '\033':
			// set terminal title
			if bytes.HasPrefix(lineb[i:], []byte("\x1b]0;")) {
				pos := bytes.Index(lineb[i:], []byte("\a"))
				if pos != -1 {
					i += pos + 1
					continue
				}
			}
			if m := lineEditRe.FindSubmatch(lineb[i:]); m != nil {
				i += len(m[0])
				n, err := strconv.Atoi(string(m[1]))
				if err != nil || n > 10000 {
					n = 1
				}
				switch m[2][0] {
				case '@':
					edit.Insert(bytes.Repeat([]byte{' '}, n))
				case 'C':
					edit.Move(n)
				case 'D':
					edit.Move(-n)
				case 'P':
					edit.Delete(n)
				case 'K':
					switch string(m[1]) {
					case "", "0":
						edit.ClearRight()
					case "1":
						edit.ClearLeft()
					case "2":
						edit.Clear()
					}
				}
				continue
			}
			if !(color && isColor(lineb[i:])) {
				skip := vt100scan(lineb[i:])
				if skip > 0 {
					i += skip
					continue
				}
			} else {
				hadColor = true
				edit.Write([]byte{c})
			}
		default:
			if c == '\n' || c >= ' ' {
				edit.Write([]byte{c})
			}
		}
		i += 1
	}
	out := edit.Bytes()
	if hadColor {
		out = append(out, []byte("\033[0m")...)
	}
	return string(out)
}
