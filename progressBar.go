package main

import (
	"fmt"
	lolgopher "github.com/kris-nova/lolgopher"
	terminal "golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

var spinnerStrings = [][2]string{
	[2]string{"⡇   ", "    "},
	[2]string{"⠏   ", "    "},
	[2]string{"⠋⠁  ", "    "},
	[2]string{"⠉⠉  ", "    "},
	[2]string{"⠈⠉⠁ ", "    "},
	[2]string{" ⠉⠉ ", "    "},
	[2]string{" ⠈⠉⠁", "    "},
	[2]string{"  ⠉⠉", "    "},
	[2]string{"  ⠈⠙", "    "},
	[2]string{"   ⠹", "    "},
	[2]string{"   ⢸", "    "},
	[2]string{"   ⢰", "   ⠈"},
	[2]string{"   ⢠", "   ⠘"},
	[2]string{"   ⢀", "   ⠸"},
	[2]string{"    ", "   ⢸"},
	[2]string{"    ", "   ⣰"},
	[2]string{"    ", "  ⢀⣠"},
	[2]string{"    ", "  ⣀⣀"},
	[2]string{"    ", " ⢀⣀⡀"},
	[2]string{"    ", " ⣀⣀ "},
	[2]string{"    ", "⢀⣀⡀ "},
	[2]string{"    ", "⣀⣀  "},
	[2]string{"    ", "⣄⡀  "},
	[2]string{"    ", "⣆   "},
	[2]string{"    ", "⡇   "},
	[2]string{"⡀   ", "⠇   "},
	[2]string{"⡄   ", "⠃   "},
	[2]string{"⡆   ", "⠁   "},
	[2]string{"⡆   ", "⠁   "},
}

var nyanStrings = [2]string{
	"\u001b[95m╔═══╗\u001b[0m⠡⠡",
	"\u001b[95m╚═\u001b[90m(OﻌO)\u001b[0m",
}

type DummyWriter struct {
	buffer *([]byte)
}

func (dw DummyWriter) Write(p []byte) (n int, err error) {
	*dw.buffer = append(*dw.buffer, p...)
	return len(p), nil
}

func (dw *DummyWriter) ReadBuffer() (buffer []byte) {
	buffer = *dw.buffer
	*dw.buffer = (*dw.buffer)[:0]
	return
}

type ProgressBarWriter struct {
	mux         sync.Mutex
	wnl         *WebNodeList
	barShown    bool
	lastBarLen  int
	bar1        string
	bar2        string
	currentSpin int
	lol         bool
	lolWriter   io.Writer
	lolDw       DummyWriter
}

func CreateProgressBarWriter(wnl *WebNodeList, showProgress, lol bool) (pbwRef *ProgressBarWriter, pbwoRef *PbwStdout, pbweRef *PbwStderr) {
	var pbw ProgressBarWriter
	var pbwo PbwStdout
	var pbwe PbwStderr
	pbwo.pbw = &pbw
	pbwe.pbw = &pbw
	pbw.wnl = wnl
	pbw.lol = lol
	pbw.barShown = showProgress
	if lol {
		lolbuf := make([]byte, 0)
		pbw.lolDw.buffer = &lolbuf
		pbwo.outputWriter = &lolgopher.Writer{
			Output:    os.Stdout,
			ColorMode: lolgopher.ColorModeTrueColor,
		}
		pbw.lolWriter = &lolgopher.Writer{
			Output:    &pbw.lolDw,
			ColorMode: lolgopher.ColorModeTrueColor,
		}
	} else {
		pbwo.outputWriter = os.Stdout
	}
	ticker := time.NewTicker(time.Millisecond * 40)
	go func() {
		for {
			<-ticker.C
			pbw.currentSpin++
			if pbw.currentSpin >= len(spinnerStrings) {
				pbw.currentSpin = 0
			}
			pbw.updateBar()
		}
	}()
	return &pbw, &pbwo, &pbwe
}

func (pbw *ProgressBarWriter) HideBar() {
	if pbw.barShown {
		esc := "\033[A\r"
		os.Stderr.WriteString(esc + strings.Repeat(" ", pbw.lastBarLen) + "\n" + strings.Repeat(" ", pbw.lastBarLen) + esc)
		pbw.barShown = false
	}
}

func (pbw *ProgressBarWriter) ShowBar() {
	if !pbw.barShown {
		os.Stderr.WriteString(pbw.bar1 + "\n" + pbw.bar2)
		pbw.barShown = true
	}
}

func (pbw *ProgressBarWriter) getRainbow(input []byte) (rainbow []byte) {
	input = append(input, '\n')
	pbw.lolWriter.Write(input)
	rainbow = pbw.lolDw.ReadBuffer()
	rainbow = rainbow[:len(rainbow)-1]
	return
}

func (pbw *ProgressBarWriter) updateBar() {
	pbw.mux.Lock()
	defer pbw.mux.Unlock()
	if !pbw.barShown {
		return
	}
	pbw.HideBar()
	defer pbw.ShowBar()
	termWidth, _, _ := terminal.GetSize(2)
	done, fail, total := pbw.wnl.GetStats()
	barWidth := int(termWidth - 55)
	doneWidth := barWidth * done / total
	percent := 100 * done / total
	padWidth := barWidth - doneWidth
	var loadBar1 string
	var loadBar2 string
	if pbw.lol {
		if barWidth > 8 {
			doneWidth -= 8
			if doneWidth < 0 {
				padWidth += doneWidth
				doneWidth = 0
			}
			loadBar1 = string(pbw.getRainbow([]byte(strings.Repeat("⡪", doneWidth)))) + nyanStrings[0] + strings.Repeat("⠡", padWidth) + " ┋"
			loadBar2 = string(pbw.getRainbow([]byte(strings.Repeat("⡪", doneWidth)))) + nyanStrings[1] + strings.Repeat("⠡", padWidth) + " ┋"
		}
	} else {
		if barWidth > 0 {
			loadBar1 = strings.Repeat("█", doneWidth) + strings.Repeat("░", padWidth) + " ┋"
			loadBar2 = loadBar1
		}
	}
	pbw.bar1 = fmt.Sprintf("%s ┋ Failed:%12d ┋   Total   ┋ %s Progress ┋", spinnerStrings[pbw.currentSpin][0], fail, loadBar1)
	pbw.bar2 = fmt.Sprintf("%s ┋ Completed:%9d ┋%10d ┋ %s   %3d%%   ┋", spinnerStrings[pbw.currentSpin][1], done, total, loadBar2, percent)
	pbw.lastBarLen = int(termWidth)
}

type PbwStderr struct {
	pbw *ProgressBarWriter
}

func (pbwe *PbwStderr) Write(p []byte) (n int, err error) {
	pbwe.pbw.mux.Lock()
	defer pbwe.pbw.mux.Unlock()
	if pbwe.pbw.barShown {
		pbwe.pbw.HideBar()
		defer pbwe.pbw.ShowBar()
	}
	return os.Stderr.Write(p)
}

// This is buffered to ensure the progress bar flickers as little as possible. It won't flush automatically. Flush must be called once this is done.
type PbwStdout struct {
	pbw          *ProgressBarWriter
	outputWriter io.Writer
	buffer       []byte
}

func (pbwo *PbwStdout) Write(p []byte) (n int, err error) {
	pbwo.buffer = append(pbwo.buffer, p...)
	return len(p), nil
}

func (pbwo *PbwStdout) Flush() {
	if len(pbwo.buffer) > 0 {
		pbwo.pbw.mux.Lock()
		defer pbwo.pbw.mux.Unlock()
		if pbwo.pbw.barShown {
			pbwo.pbw.HideBar()
			defer pbwo.pbw.ShowBar()
		}
		pbwo.outputWriter.Write(pbwo.buffer)
		pbwo.buffer = pbwo.buffer[:0]
	}
}
