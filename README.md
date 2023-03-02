# GoVelocity - Simple Go CI/CD Tool

A simple CI/CD tool written in Go that reads commands to be executed from a YAML configuration file based on a specific GitHub webhook event.

## Usage

To use the tool, you must first create a YAML configuration file in your repository specifying the commands to be executed for each GitHub webhook event. An example configuration file is provided in `config.yml.example`. You should copy this file to `config.yml` and modify it according to your needs.

The `name` field specifies the name of the command, and the `command` field specifies the actual command to be executed. You can include any shell command, and the tool will execute it.
By default, the tool will read the `config.yml` file and execute the commands defined in the `push` event. You can specify a different event by passing it as an argument to the command:

### CLI

The tool also comes with a command-line interface (CLI) that allows you to easily configure and run the tool. The CLI provides the following commands:

