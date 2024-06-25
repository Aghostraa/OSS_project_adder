# OSS Project yaml Creator

A Go application to help users create YAML files for web3 projects and automatically commit them to a OSS project db.

## Features

- Step-by-step prompts to gather project details
- Generates a YAML file based on user input
- Automatically commits and pushes the YAML file to a specified GitHub repository

## Prerequisites

- [Go](https://golang.org/dl/) installed
- [Git](https://git-scm.com/downloads) installed
- A GitHub account with a forked repository

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/<your-github-username>/yaml-project-creator.git
    cd yaml-project-creator
    ```

2. Build the application:

    ```sh
    go build -o yaml_project_creator
    ```

## Usage

1. Run the application:

    ```sh
    ./yaml_project_creator
    ```

2. Follow the prompts to enter your project details.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License.
