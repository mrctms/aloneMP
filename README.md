# aloneMP

A minimal music player terminal based.
Insipred by <a href="https://github.com/Depado/parakeet">parakeet</a>

# Download and installation on Linux

Download <a href="https://github.com/marcktomack/aloneMP/releases">here</a>

Extract the zip

```
cd aloneMP
chmod +x install.sh
sudo ./install.sh
```

# Usage

```
Usage of aloneMP:
  -addr string
        http server address (default "127.0.0.1:3777")
  -dir string
        Directory with audio files (default "~/Music")
  -s    Run the http server

```

You can pass a directory with audio files:

`aloneMP -dir /your/directory`

Or by default it will load audio files from `~/Music`

You can also run the http server with -s and pass the address with -addr (default is 127.0.0.1:3777)

Docs for http server and http server it self is under construction

Supported formats are <b>MP3</b>, <b>FLAC</b>, <b>OGG</b> and <b>WAV</b>

# Build from source

You need a <a href="https://golang.org/dl/">Go compiler</a>

```
git clone https://github.com/MarckTomack/aloneMP
cd aloneMP
go build
```

# Depndencies

aloneMP uses <a href="https://github.com/faiface/beep">beep</a> that uses <a href="https://github.com/hajimehoshi/oto">oto</a> see the dependencies <a href="https://github.com/hajimehoshi/oto#prerequisite">here</a>

