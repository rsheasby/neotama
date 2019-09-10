package main

import (
	lolgopher "github.com/kris-nova/lolgopher"
	"io"
	"net/url"
	"os"
	"strings"
)

type WnlPrinter struct {
	wnl          *WebNodeList
	outputFormat OutputFormat
	outputWriter io.Writer
	colorOption  ColorOption
	colorCodes   map[string]string
	treeLines    []bool
	lastDepth    int
	printIndex   int
}

func CreateWnlPrinter(wnl *WebNodeList, of OutputFormat, colorOption ColorOption) (wp WnlPrinter) {
	wp.wnl = wnl
	wp.outputFormat = of
	wp.colorOption = colorOption
	if wp.colorOption == lol {
		wp.outputWriter = lolgopher.NewTruecolorLolWriter()
	} else {
		wp.outputWriter = os.Stdout
	}
	cc, ccExists := os.LookupEnv("LS_COLORS")
	if !ccExists {
		// I know this is kinda bad and it would be better for this to be a compile time constant map, but this only takes up one line and is less work.
		cc = "rs=0:su=37;41:di=01;34:*.tar=01;31:*.tgz=01;31:*.arc=01;31:*.arj=01;31:*.taz=01;31:*.lha=01;31:*.lz4=01;31:*.lzh=01;31:*.lzma=01;31:*.tlz=01;31:*.txz=01;31:*.tzo=01;31:*.t7z=01;31:*.zip=01;31:*.z=01;31:*.dz=01;31:*.gz=01;31:*.lrz=01;31:*.lz=01;31:*.lzo=01;31:*.xz=01;31:*.zst=01;31:*.tzst=01;31:*.bz2=01;31:*.bz=01;31:*.tbz=01;31:*.tbz2=01;31:*.tz=01;31:*.deb=01;31:*.rpm=01;31:*.jar=01;31:*.war=01;31:*.ear=01;31:*.sar=01;31:*.rar=01;31:*.alz=01;31:*.ace=01;31:*.zoo=01;31:*.cpio=01;31:*.7z=01;31:*.rz=01;31:*.cab=01;31:*.wim=01;31:*.swm=01;31:*.dwm=01;31:*.esd=01;31:*.jpg=01;35:*.jpeg=01;35:*.mjpg=01;35:*.mjpeg=01;35:*.gif=01;35:*.bmp=01;35:*.pbm=01;35:*.pgm=01;35:*.ppm=01;35:*.tga=01;35:*.xbm=01;35:*.xpm=01;35:*.tif=01;35:*.tiff=01;35:*.png=01;35:*.svg=01;35:*.svgz=01;35:*.mng=01;35:*.pcx=01;35:*.mov=01;35:*.mpg=01;35:*.mpeg=01;35:*.m2v=01;35:*.mkv=01;35:*.webm=01;35:*.ogm=01;35:*.mp4=01;35:*.m4v=01;35:*.mp4v=01;35:*.vob=01;35:*.qt=01;35:*.nuv=01;35:*.wmv=01;35:*.asf=01;35:*.rm=01;35:*.rmvb=01;35:*.flc=01;35:*.avi=01;35:*.fli=01;35:*.flv=01;35:*.gl=01;35:*.dl=01;35:*.xcf=01;35:*.xwd=01;35:*.yuv=01;35:*.cgm=01;35:*.emf=01;35:*.ogv=01;35:*.ogx=01;35:*.aac=00;36:*.au=00;36:*.flac=00;36:*.m4a=00;36:*.mid=00;36:*.midi=00;36:*.mka=00;36:*.mp3=00;36:*.mpc=00;36:*.ogg=00;36:*.ra=00;36:*.wav=00;36:*.oga=00;36:*.opus=00;36:*.spx=00;36:*.xspf=00;36"
	}
	ccl := strings.Split(cc, ":")
	wp.colorCodes = make(map[string]string, len(ccl))
	for i := range ccl {
		splitcc := strings.Split(ccl[i], "=")
		if len(splitcc) != 2 {
			continue
		}
		wp.colorCodes[splitcc[0]] = splitcc[1]
	}
	return
}

func (wp *WnlPrinter) treePrintNode(index int) {
	if index < 0 || index >= len(wp.wnl.list) {
		return
		// TODO: Proper error return
	}
	node := wp.wnl.list[index]
	if node.nodeDepth >= len(wp.treeLines) {
		wp.treeLines = append(wp.treeLines, false)
	}
	for i := 0; i < node.nodeDepth-1; i++ {
		if i < len(wp.treeLines) && wp.treeLines[i] {
			wp.outputWriter.Write([]byte("│   "))
		} else {
			wp.outputWriter.Write([]byte("    "))
		}
	}
	if node.nodeDepth > 0 {
		if node.nodeLastSibling {
			wp.outputWriter.Write([]byte("└── "))
			wp.treeLines[node.nodeDepth-1] = false
		} else {
			wp.outputWriter.Write([]byte("├── "))
			wp.treeLines[node.nodeDepth-1] = true
		}
	}
	if wp.colorOption == on {
		setSgr(wp.getSgrCode(&wp.wnl.list[wp.printIndex]))
	}
	wp.outputWriter.Write([]byte(node.name + "\n"))
}

func setSgr(sgrCode string) {
	os.Stdout.WriteString("\u001b[" + sgrCode + "m")
}

func (wp *WnlPrinter) resetSgr() {
	setSgr(wp.colorCodes["rs"])
}

func (wp *WnlPrinter) getSgrCode(node *WebNode) (sgrCode string) {
	if node.nodeFail {
		return wp.colorCodes["su"]
	}
	if node.nodeType == directory {
		return wp.colorCodes["di"]
	} else {
		split := strings.Split(node.path, ".")
		ext := split[len(split)-1]
		sgr, valid := wp.colorCodes["*."+ext]
		if !valid {
			sgr = wp.colorCodes["rs"]
		}
		return sgr
	}
}

func (wp *WnlPrinter) PrintDone() {
	l := wp.wnl
	l.mux.Lock()
	defer l.mux.Unlock()
	for ; wp.printIndex < len(l.list); wp.printIndex++ {
		if l.list[wp.printIndex].nodeStatus == done {
			switch wp.outputFormat {
			case tree:
				wp.treePrintNode(wp.printIndex)
			case urlencoded:
				if wp.colorOption == on {
					setSgr(wp.getSgrCode(&l.list[wp.printIndex]))
				}
				wp.outputWriter.Write([]byte(wp.wnl.list[wp.printIndex].path + "\n"))
			case list:
				if wp.colorOption == on {
					setSgr(wp.getSgrCode(&l.list[wp.printIndex]))
				}
				path, _ := url.PathUnescape(wp.wnl.list[wp.printIndex].path)
				wp.outputWriter.Write([]byte(path + "\n"))
			}
			wp.resetSgr()
		} else {
			return
		}
	}
}
