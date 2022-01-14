# go-scaffold

A simple golang scaffolding to help me to create new api projects or workers with golang on k8s.

## Patterns

* Struct folder Layout: (Modern Go Application)[https://github.com/sagikazarmark/modern-go-application], (Golang Standard Layout)[https://github.com/golang-standards/project-layout];

## Code flow

![alt text](./docs/assets/architecture.drawio.svg)

## Development

How to develop with this project.

### VS Code and Remote-Control Plugin

1. Install Remote-Control plugin on VS Code.
2. copy file ./env/application.env.sample to ./env/application.env
3. Reopen in Container mode, like [this](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
4. Run the command `make hot`, for start with hot reload or on main.go file opened debug with pressing "f5".

### Lint Troublehooting

Take a look into [this](https://github.com/facily-tech/go-core/blob/main/LINT.md).
