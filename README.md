# opengl-experimentation

Let's have fun with OpenGL in GO.

## Setup

GLFW Depencies: <https://github.com/go-gl/glfw#installation>

On windows, install a cgo compiler (Ex: TDM-GCC)

```bash
tdm64-gcc-9.2.0.exe from: https://jmeubank.github.io/tdm-gcc/
```

Install OpenGL and GLFW cgo binding libraries

```bash
go get -u github.com/go-gl/gl/v4.6-core/gl
go get -u github.com/go-gl/glfw/v3.3/glfw
```

## Build

Build the project

```bash
git clone git@github.com:devodev/opengl-experimentation.git
cd opengl-experimentation
go build .
```
