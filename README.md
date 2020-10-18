# aloneMP

A minimal music player terminal based.

aloneMP is composed of two elements, the daemon (aloneMPd) and the clients, for now is avaiable only the terminal client with a TUI

<img src="demo.gif" alt="demo">

# Installation

With the installers aloneMP client and the daemon are installed on the same machine, if you want configure all the things manually there are the zip with the binary.

## Debian and friends

You need to have systemd.

Download the .deb package <a href="https://github.com/marcktomack/aloneMP/releases">here</a>

## Windows

Download the installer <a href="https://github.com/marcktomack/aloneMP/releases">here</a>



# Usage

```
Usage of aloneMP:
  -addr string
        aloneMP daemon address (default "hostname:3777")
  -dir string
        Directory with audio files (default "~/Music")
  -tui
        run tui client (default true)
  -version
        show version

```

You can pass a directory with audio files:

`aloneMP -dir /your/directory`

Or by default it will load audio files from `~/Music`

You need to provide the daemon server address, the default one is hostname:3777, same on the daemon

Supported formats are <b>MP3</b>, <b>FLAC</b>, <b>OGG</b> and <b>WAV</b>

# TODO

- Web interface

Since I'm not a UI guy, if anyone want make a web interface for aloneMP are welcome!<br>

Drop me an email before start the work.


