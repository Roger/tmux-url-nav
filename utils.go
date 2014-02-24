package main

import (
    "fmt"
    "regexp"
    "strconv"
    "os/exec"
    "strings"
)

func colored(text string, color *Color) string {
    colorstr := fmt.Sprintf("\033[38;5;%vm", color.FOREGROUND)
    colorstr += fmt.Sprintf("\033[48;5;%vm", color.BACKGROUND)

    attrs := ""
    if color.BOLD {
        attrs += "1;"
    }
    if color.DIM {
        attrs += "3;"
    }
    if color.UNDERLINED {
        attrs += "4;"
    }
    if color.BLINK {
        attrs += "5;"
    }
    if color.REVERSE {
        attrs += "7;"
    }
    if color.HIDDEN {
        attrs += "8;"
    }
    if attrs != "" {
        attrs = fmt.Sprintf("\033[%sm", strings.TrimSuffix(attrs, ";"))
    }

    return fmt.Sprintf("%s%s%s\033[0m", colorstr, attrs, text)
}

func moveCursor(x int, y int) string {
    return fmt.Sprintf("\033[%d;%dH", x+1, y+1)
}

func tmuxDisplayMsg(msg string) {
    exec.Command("tmux", "display-message", "[tmux-url-nav] "+msg).Run()
}

func tmuxCapturePane() {
    exec.Command("tmux", "capture-pane", "-eJ").Run()
}

func tmuxOpenInnerWindow(title string, command string, windowid int) {
    exec.Command("tmux", "new-window", "-dn", title,
                 "-t", strconv.Itoa(windowid),
                 command + " inner").Run()
}

func tmuxSelectWindow(windowid int) {
    exec.Command("tmux", "select-window", "-t", strconv.Itoa(windowid)).Run()
}

func getBuffer() string {
    out, err := exec.Command("tmux", "show-buffer").Output()

    if err != nil {
        panic(err)
    }

    buffer := strings.TrimSuffix(string(out), "\n")
    return buffer
}

func getLines(buffer string) []string {
    rectrl := regexp.MustCompile(`\x1b[^m]*m`)
    reunicode := regexp.MustCompile(`[^\x00-\x7F]`)

    return strings.Split(
        // Replace Unicode Characters
        reunicode.ReplaceAllString(
            // Remove Escape Sequences
            rectrl.ReplaceAllString(buffer, ""),
            " "),
        "\n")
}
