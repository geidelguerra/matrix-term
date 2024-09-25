package main

import "fmt"
import "bytes"
import "os"
import "golang.org/x/term"
import "time"
import "math/rand"

func ClearScreen(b *bytes.Buffer) {
	b.WriteString("\033[2J")
}

func SetCursor(buffer *bytes.Buffer, x, y int) {
	buffer.WriteString(fmt.Sprintf("\033[%d;%dH", y, x))
}

func SetForegroundColor(buffer *bytes.Buffer, r, g, b int) {
	buffer.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b))
}

func ResetForegroundColor(buffer *bytes.Buffer) {
	buffer.WriteString("\033[39m")
}

func SetBackgroundColor(buffer *bytes.Buffer, r, g, b int) {
	buffer.WriteString(fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b))
}

func ResetBackroundColor(buffer *bytes.Buffer) {
	buffer.WriteString("\033[49m")
}

func main() {
	var buffer bytes.Buffer
	
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	
	defer term.Restore(int(os.Stdin.Fd()), state)
	
	quit := false

	go func () {
		for {
			var buf [64]byte
			n, err := os.Stdin.Read(buf[:])
			if err != nil {
				quit = true
				panic(err)
			}

			if n == 1 && buf[0] == 0x03 {
				quit = true
				break
			}
		}
	}()
	
	os.Stdout.WriteString("\x1b[?25l")
	
	screenW, screenH, _ := term.GetSize(int(os.Stdout.Fd()))

	charOffset := 0
	yOffsets := make([]int, screenW)
	tailLengths := make([]int, screenW)
	
	for i := 0; i < screenW; i ++ {
		yOffsets[i] = rand.Intn(10)
		tailLengths[i] = rand.Intn(screenH - 2)
	}

	for !quit {	
		buffer.Reset()
		ClearScreen(&buffer)
		
		x := 1
		for x < screenW {
			y := 1
			var r rune
			yOffset := yOffsets[x - 1]
			tailLength := tailLengths[x - 1]

			for y < screenH {
				SetCursor(&buffer, x, y)
				SetBackgroundColor(&buffer, 0, 0, 0)

				if y <= yOffset + 1 && y >= yOffset - tailLength + 2  {
					SetForegroundColor(&buffer, 3, 160, 98)
				} else {
					SetForegroundColor(&buffer, 20, 20, 20)
				}

				r = rune(0x30a1 + charOffset)
				buffer.WriteRune(r)

				charOffset = (charOffset + 1) % 95
				y += 1
			}

			yOffsets[x - 1] = (yOffset + 1) % screenH
			x += 2
		}

		buffer.WriteTo(os.Stdout)
		time.Sleep(100 * time.Millisecond)
	}

	os.Stdout.WriteString("\x1b[?25h")
}
