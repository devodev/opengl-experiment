# opengl-experimentation

Let's have fun with OpenGL in GO.

## Setup

Install OpenGL and GLFW

```bash
# GLFW dependencies (Ubuntu on WSL2)
# See here fore details: https://github.com/go-gl/glfw#installation
sudo apt-get update
sudo apt-get install libgl1-mesa-dev xorg-dev

go get -u github.com/go-gl/gl/v4.6-core/gl
go get -u github.com/go-gl/glfw/v3.3/glfw
```