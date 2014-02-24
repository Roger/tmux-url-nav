package main

import (
    "os/user"
    "io/ioutil"
    // external deps
    "launchpad.net/goyaml"
)

type Color struct {
    BACKGROUND int
    FOREGROUND int
    BOLD       bool
    DIM        bool
    BLINK      bool
    REVERSE    bool
    HIDDEN     bool
    UNDERLINED bool
}

type Colors struct {
    DEFAULT_URL Color
    ACTIVE_URL  Color
    POSITION    Color
}

type Config struct {
    TITLE         string
    WINDOWID      int
    COLORS        Colors
    OPENER        string
    SHOW_POSITION bool
}

func NewConfig() *Config {
    config := &Config{OPENER: "firefox", TITLE: "Select Url", WINDOWID: 999,
                      SHOW_POSITION: true,
        COLORS: Colors{
            DEFAULT_URL: Color{FOREGROUND: 6},
            ACTIVE_URL:  Color{BACKGROUND: 17, FOREGROUND: 7,
                               UNDERLINED: true},
            POSITION:    Color{FOREGROUND: 226},
        }}
    config.parseConfig()
    return config
}


func (config *Config) parseConfig() {
    usr, err := user.Current()
    if err != nil {
        panic(err)
    }

    data, err := ioutil.ReadFile(usr.HomeDir + "/.config/tmux-url-config.yml")

    // if can't open the config, just ignore it
    if err != nil {
        return
    }

    err = goyaml.Unmarshal(data, &config)
    if err != nil {
        panic(err)
    }
}
