# OSS Project YAML Creator

A Go application and Chrome extension to help users create YAML files for web3 projects and automatically commit them to an OSS project database.

## Features

- Step-by-step prompts to gather project details
- Generates a YAML file based on user input
- Automatically commits and pushes the YAML file to a specified GitHub repository
- Chrome extension for an easy-to-use interface

## Prerequisites

- [Go](https://golang.org/dl/) installed
- [Git](https://git-scm.com/downloads) installed
- A GitHub account with a forked repository
- [Node.js](https://nodejs.org/en/download/) and [npm](https://www.npmjs.com/get-npm) installed (for building the Chrome extension)

## Installation

### Go Application

1. Clone the repository:

    ```sh
    git clone https://github.com/aghostraa/yaml-project-creator.git
    cd yaml-project-creator
    ```

2. Build the application:

    ```sh
    go build -o yaml_project_creator
    ```

### Chrome Extension

1. Navigate to the extension directory:

    ```sh
    cd oss_project_adder_ext
    ```

2. Install the dependencies:

    ```sh
    npm install
    ```

3. Build the extension:

    ```sh
    npm run build
    ```

## Usage

### Go Application

1. Run the application:

    ```sh
    ./yaml_project_creator
    ```

2. Follow the prompts to enter your project details.

### Chrome Extension

1. Open Chrome and go to `chrome://extensions/`.
2. Enable "Developer mode" using the toggle in the top right.
3. Click on "Load unpacked" and select the `oss_project_adder_ext` directory.
4. The extension should now be loaded and its icon should appear in the toolbar.
5. Click the extension icon and follow the prompts to enter your project details.
6. The extension will communicate with the local Go server to create and manage YAML files.

### Running the Go Server

To use the Chrome extension, ensure that the Go server is running:

1. Run the server:

    ```sh
    go run main.go
    ```

2. The server will start on `http://localhost:8080`.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License.
