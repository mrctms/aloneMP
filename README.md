# aloneMP

A minimal music player terminal based.
Insipred by <a href="https://github.com/Depado/parakeet">parakeet</a>

<img src="demo.gif" alt="demo">

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

