# opengl-experimentation

![Current Project State](docs/images/currentProjectState.gif?raw=true "Current Project State")

## Setup

Install GLFW Depencies: <https://github.com/go-gl/glfw#installation>

Install a cgo compiler(Windows): <https://jmeubank.github.io/tdm-gcc/>

Install OpenGL and GLFW cgo binding libraries

```bash
go get -u github.com/go-gl/gl/v4.6-core/gl
go get -u github.com/go-gl/glfw/v3.3/glfw
```

## Build

Build the project

```bash
git clone https://github.com/devodev/opengl-experimentation.git
cd opengl-experimentation
go build .
```
