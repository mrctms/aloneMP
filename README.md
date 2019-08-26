# aloneMP

A minimal music player terminal based.

# Download and installation on Linux

Download <a href="https://github.com/marcktomack/aloneMP/releases">here</a>

Extract the zip

```
cd aloneMP
chmod +x install.sh
sudo ./install.sh
```

# Usage

You can pass a directory with audio files:

`aloneMP -dir /your/directory`

Or by default it will load audio files from `~/Music`

Supported formats are <b>MP3</b>, <b>FLAC</b> and <b>WAV</b>

# Build from source

You need a <a href="https://golang.org/dl/">Go compiler</a>

```
git clone https://github.com/MarckTomack/aloneMP
cd aloneMP
go build
```
