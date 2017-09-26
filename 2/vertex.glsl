#version 330 core
layout (location = 0) in vec3 aPos; // the position variable has attribute position 0
layout (location = 1) in vec3 aColor; // the color variable has attribute position 1

out vec3 ourColor; // specify a color output to the fragment shader

void main()
{
   gl_Position = vec4(aPos, 1.0); // example of giving a vec3 in a vec4 constructor
   ourColor = aColor; // set ourColor to the input color we got from the vertex data
}
