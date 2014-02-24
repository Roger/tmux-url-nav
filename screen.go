package main

import (
    "os"
    "fmt"
    "strconv"
    "os/exec"
    "strings"
    "regexp"
)

var url_pattern = `(?:https?://|ftp://|news://|git://|mailto:|file://|www\.)` +
    `[\w\-\@;\/?:&=%\$_.+!*\x27(),~#]+[\w\-\@;\/?&=%\$_+!*\x27()~]`

var url_regex = regexp.MustCompile(url_pattern)


type Screen struct {
    buffer string
    config *Config
}

func NewScreen() *Screen {
    screen := &Screen{buffer:getBuffer(), config: NewConfig()}
    screen.initScreen()
    return screen
}

func (screen *Screen) initScreen() {
    // Disable input buffering and not display entered characters
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1", "-echo").Run()

    // Hide cursor
    fmt.Print("\033[?25l")

    // Display buffer
    fmt.Print(screen.buffer)
}

func (screen *Screen) findUrls(selected int, lines []string) (string, int) {
    current := 0
    url := ""
    dis := 0
    num := ""
    for i, line := range lines {
        for _, match := range url_regex.FindAllStringIndex(line, -1) {
            color := screen.config.COLORS.DEFAULT_URL
            if selected == current {
                color = screen.config.COLORS.ACTIVE_URL
                url = line[match[0]:match[1]]
            }

            if screen.config.SHOW_POSITION {
                // Displacement of the identifier
                dis = len(fmt.Sprintf("%d", current+1))
                // Identifier of the match
                num = colored(strconv.Itoa(current+1),
                    &screen.config.COLORS.POSITION)
            }
            // Position in screen of the match
            pos := moveCursor(i, match[0])
            // Text of the match
            str := colored(line[match[0]+dis:match[1]], &color)

            // Replace match in the screen
            fmt.Print(pos + num + str)
            current += 1
        }
    }
    return url, current
}

func (screen *Screen) handleUserInput() {
    // Get ascii lines of the buffer
    lines := getLines(screen.buffer)

    selected := 0
    b := make([]byte, 1)

MainLoop:
    for {
        url, current := screen.findUrls(selected, lines)

        if current == 0 {
            tmuxDisplayMsg("No URLs found")
            break
        }

        os.Stdin.Read(b)
        key := string(b)

        switch {
        case key == "j":
            selected++
        case key == "k":
            selected--
        case key == "o":
            exec.Command(screen.config.OPENER, url).Run()
            break MainLoop
        case key == "O": // Don't quit, just open
            exec.Command(screen.config.OPENER, url).Run()
        case key == "y":
            cmd := exec.Command("xsel", "-bi")
            cmd.Stdin = strings.NewReader(url)
            cmd.Run()
            break MainLoop
        case key == "q":
            break MainLoop
        case b[0] >= '0' && b[0] <= '9':
            selected, _ = strconv.Atoi(key)
            selected -= 1
        }

        if selected > current-1 {
            selected = 0
        } else if selected < 0 {
            selected = current - 1
        }
    }
}

