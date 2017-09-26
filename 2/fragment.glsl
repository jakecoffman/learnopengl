#version 330 core
out vec4 FragColor;
in vec3 ourColor; // the input variable is from the vertex shader (same name and type)

void main()
{
   FragColor = vec4(ourColor, 1.0);
}
