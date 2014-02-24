package main

import "os"

func main() {
    tmuxCapturePane()
    if len(os.Args) == 1 {
        config := NewConfig()
        tmuxOpenInnerWindow(config.TITLE, os.Args[0], config.WINDOWID)
    } else if os.Args[1] == "inner" {
        screen := NewScreen()
        screen.handleUserInput()
    }
}
