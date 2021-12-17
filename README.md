# opengl-experiment

![Current Project State](docs/images/currentProjectState2.gif?raw=true "Current Project State")

## Setup

Install a cgo compiler(Windows): <https://jmeubank.github.io/tdm-gcc/>

Install GLFW Depencies: <https://github.com/go-gl/glfw#installation>

Install OpenGL and GLFW cgo binding libraries

```bash
go get -u github.com/go-gl/gl/v4.6-core/gl
go get -u github.com/go-gl/glfw/v3.3/glfw
```

## Example

Run the example project

```bash
git clone https://github.com/devodev/opengl-experiment.git
cd opengl-experiment/examples/textured_quad
go run .
```
