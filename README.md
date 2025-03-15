# Go Generator

A powerful code generator tool for Go projects that helps streamline development by automating common code patterns and boilerplate.

## Features

- Automated code generation for common patterns
- Built-in templates for various use cases
- Customizable code templates
- Easy to use CLI interface

## Installation the package

```bash
go install github.com/indalyadav56/go-generator@latest
```

## Examples

```bash
# Initialize a backend project with authentication, user, and todo apps
go-generator new backend --app=auth --app=user --app=todo #add more if needed

# Generate a new app
go-generator app another-app

<!-- with htmx -->

go-generator new project --app=auth --app=user --app=todo --frontend=htmx

<!-- with react -->

go-generator new project --app=auth --app=user --app=todo --frontend=react

```

## Screenshots

![Go Generator Screenshot](/screenshots/image.png)

## Source Code Clone

Clone the repository with submodules:

```bash
git clone --recurse-submodules https://github.com/indalyadav56/go-generator.git
```

## Prerequisites

- Go 1.16 or higher
- Git

## Usage

1. Navigate to the project directory:

```bash
cd go-generator
```

2. Build the project:

```bash
make build
```

3. Run the generator:

```bash
./bin/go-generator [command] [flags]

Example:
    ./bin/go-generator init backend --app=auth --app=user --app=todo

    or

    go run main.go init backend --app=auth --app=user --app=todo
```

## Available Commands

- `init`: Initialize a new project
- `generate`: Generate code from templates
- More commands coming soon...

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository.
