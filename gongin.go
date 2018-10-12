package gongin

import (
	"github.com/gubsky90/gongin/render"
)

type Config struct {

}

type Gongin struct {
	handlers map[string][]func()
}

func New(conf Config) *Gongin {
	g := Gongin{}
	g.handlers = make(map[string][]func())
	return &g
}

func (g *Gongin) On(name string, cb func()) {
	if _, ok := g.handlers[name]; !ok {
		g.handlers[name] = make([]func(), 0)
	}
	g.handlers[name] = append(g.handlers[name], cb)
}

func (g *Gongin) fire(name string) {
	if cbs, ok := g.handlers[name]; ok {
		for _, cb := range cbs {
			cb()
		}
	}
}

var points = []float32 {
	-1.0,  1.0,  0.0,
	-1.0, -1.0,  0.0,
	 1.0,  1.0,  0.0,

	 1.0,  1.0,  0.0,
	-1.0, -1.0,  0.0,
	 1.0, -1.0,  0.0,
}

func (g *Gongin) Run() {
	r := render.New()
	defer r.Close()

	// run := true

	// r.On("close", func(){
	// 	run = false
	// })

	mesh := render.NewMesh(points)

	shader := render.NewShader(render.ShaderSource{
		Vertex: `
			#version 410
			in vec3 vp;
			void main() {
				gl_Position = vec4(vp, 1.0);
			}
		`,
		Fragment: `
			#version 410
			#define WAVES 8.0
			in vec4 gl_FragCoord;
			out vec4 fragColor;

			uniform vec3 iColor;
			uniform float iTime;

			vec2 iResolution = vec2(640.0, 480.0);

			float circle(vec2 uv, float r){
			    return smoothstep(r-0.02, r, length(vec2(0.0) - uv));
			}

			float satelit(vec2 uv, vec2 p, float r, vec2 d, float t){
			  	float c = cos(t + d.x);
			    float s = sin(t + d.y);
			    p *= mat2(c, -s, s, c);
			    return circle(uv + p, r);
			}

			float atom(vec2 uv){
			    float m = 1.0;
			    m *= circle(uv, 0.22);

			    float t = iTime * 0.5;
			    float c = 0.9;
			    float s = 0.1;

			    float t1 = iTime;

			    float c1 = abs(cos(t1)) * .5 + .5;
			    for(float i=0.0; i<20.0; i++){
			     	m *= satelit(uv, vec2(.1, .4), 0.003 * i * c1, vec2(c, c - 1.0), t1 + i);
			    }

			   	t1 += 0.1;
			    float c2 = abs(cos(t1)) * .5 + .5;
			    for(float i=0.0; i<20.0; i++){
			     	m *= satelit(uv, vec2(.4, .1), 0.003 * i * c2, vec2(s, s - 1.0), t1 + i);
			    }

			   	t1 += 0.5;
			    float c3 = abs(cos(t1)) * .5 + .5;
			    for(float i=0.0; i<20.0; i++){
			     	m *= satelit(uv, vec2(.3, -.3), 0.003 * i * c3, vec2(c, s), t1 + i);
			    }
			    return m;
			}


			void mainImage(out vec4 fragColor, in vec2 fragCoord){
			    vec2 uv = fragCoord.xy/iResolution.xy - 0.5;
			    uv.x *= iResolution.x/iResolution.y;

			    uv *= 1.0 + (sin(iTime * 0.1)*.5+.5) * 0.5;

			    float m = 1.0;

			    float t = iTime;
			    float c = cos(t + 1.0);
			    float s = sin(t + 0.5);

			    mat2 rot = mat2(c, -s, s, c);

				float z1 = 1.0 + (c*.2);
			    float z2 = 1.0 + (-c*.2);
			    m *= atom((uv*3.5*z1) - vec2(-0.7, 0.0) * rot);
			    m *= atom((uv*3.5*z2) - vec2(0.7, 0.0) * rot);
			    // m *= atom((uv*2.5) - vec2(0.0, -0.6) * rot);

			    vec3 col = mix(vec3(0.2, 0.2, 0.2), vec3(0.7, 0.8, 0.9), m);

			    fragColor = vec4(col,1.0);
			}

			void main() {
				mainImage(fragColor, gl_FragCoord.xy);
			}
		`,
	})

	r.UseShader(shader)


	g.fire("ready")

	for !r.ShouldClose() {
		r.Draw(mesh)
	}
}
