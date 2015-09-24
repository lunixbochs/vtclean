package vtclean

import (
	"regexp"
)

// see regex.txt for a slightly separated version of this regex
var vt100re = regexp.MustCompile(`^\033(([A-KZ=>12<]|Y\d{2})|\[\d+[A-D]|\[\d+;\d+[Hf]|#[1-68]|\[(\d+|;)*[qm]|\[[KJg]|\[[0-2]K|\[[02]J|\([ABCEHKQRYZ0-7=]|[\[K]\d+;\d+r|\[[03]g|\[\?[1-9][lh]|\[20[lh]|\[[56]n|\[0?c|\[2;[1248]y|\[!p|\[([01457]|254)}|\[\?(12;)?(25|50)[lh]|[78DEHM]|\[[ABCDHJKLMP]|\[4[hl]|\[\?1[46][hl]|\[\*[LMP]|\[[12][JK]|\]\d*;\d*)`)
var vt100color = regexp.MustCompile(`^\033\[(\d+|;)*[m]`)

func vt100scan(line string) int {
	return len(vt100re.FindString(line))
}

func isColor(line string) bool {
	return len(vt100color.FindString(line)) > 0
}

func Clean(line string, color bool) string {
	out := make([]rune, 0, len(line))
	liner := []rune(line)
	hadColor := false
	for i := 0; i < len(liner); {
		c := liner[i]
		switch c {
		case '\b':
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
		case '\033':
			if !(color && isColor(string(liner[i:]))) {
				skip := vt100scan(string(liner[i:]))
				if skip > 0 {
					i += skip
					continue
				}
			} else {
				hadColor = true
			}
			out = append(out, c)
		default:
			if c == '\n' || c >= ' ' {
				out = append(out, c)
			}
		}
		i += 1
	}
	if hadColor {
		out = append(out, []rune("\033[0m")...)
	}
	return string(out)
}
