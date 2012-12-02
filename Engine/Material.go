package Engine

import (
	"fmt"
	"github.com/vova616/gl"
)

type Material interface {
	Load() error
	Begin(gobj *GameObject)
	End(gobj *GameObject)
}

type BasicMaterial struct {
	Program        gl.Program
	vertexShader   string
	fragmentShader string

	ViewMatrix, ProjMatrix, ModelMatrix, BorderColor, AddColor, Texture gl.UniformLocation
	Verts, UV                                                           gl.AttribLocation
}

func NewBasicMaterial(vertexShader, fragmentShader string) *BasicMaterial {
	return &BasicMaterial{Program: gl.CreateProgram(), vertexShader: vertexShader, fragmentShader: fragmentShader}
}

func (b *BasicMaterial) Load() error {
	program := b.Program
	vrt := gl.CreateShader(gl.VERTEX_SHADER)
	frg := gl.CreateShader(gl.FRAGMENT_SHADER)

	vrt.Source(vertexShader)
	frg.Source(fragmentShader)

	vrt.Compile()
	if vrt.Get(gl.COMPILE_STATUS) != 1 {
		return fmt.Errorf("Error in Compiling Vertex Shader:%s\n", vrt.GetInfoLog())
	}
	frg.Compile()
	if frg.Get(gl.COMPILE_STATUS) != 1 {
		return fmt.Errorf("Error in Compiling Fragment Shader:%s\n", frg.GetInfoLog())
	}

	program.AttachShader(vrt)
	program.AttachShader(frg)

	program.BindAttribLocation(0, "vertexPos")
	program.BindAttribLocation(1, "vertexUV")

	program.Link()

	b.Verts = program.GetAttribLocation("vertexPos")
	b.UV = program.GetAttribLocation("vertexUV")
	b.ViewMatrix = program.GetUniformLocation("MView")
	b.ProjMatrix = program.GetUniformLocation("MProj")
	b.ModelMatrix = program.GetUniformLocation("MModel")
	b.BorderColor = program.GetUniformLocation("bcolor")
	b.Texture = program.GetUniformLocation("mytexture")
	b.AddColor = program.GetUniformLocation("addcolor")
	return nil
}

func (b *BasicMaterial) Begin(gobj *GameObject) {
	b.Program.Use()
}

func (b *BasicMaterial) End(gobj *GameObject) {

}

var TextureShader gl.Program
var TextureMaterial *BasicMaterial

const vertexShader = `
#version 110

uniform mat4 MProj;
uniform mat4 MView;
uniform mat4 MModel;

attribute  vec3 vertexPos;
attribute  vec2 vertexUV; 
varying vec2 UV;


 
void main(void)
{
	gl_Position = MProj * MView * MModel * vec4(vertexPos, 1.0);
	UV = vertexUV;
}
`

const fragmentShader = `
#version 110

varying vec2 UV; 
uniform sampler2D mytexture;
uniform vec4 bcolor;
uniform vec4 addcolor;

void main(void)
{ 
  	vec4 tcolor = texture2D(mytexture, UV);
	if (tcolor.a > 0.0) {
		tcolor += bcolor;
	}
	tcolor = tcolor*addcolor;

	//nice alpha detection
	//vec4 t = addcolor;
	//t.a = 0;
	//tcolor = mix(tcolor, t, tcolor.a);

	gl_FragColor = tcolor;
}
`
