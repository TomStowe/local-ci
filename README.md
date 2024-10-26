# Local CI
Local CI is a command-line tool for running and testing CI/CD pipelines locally. It supports executing GitLab and GitHub YAML files. The tool can selectively run specific pipeline stages and displays status in tables inspired by the GitLab GUI.

## Features
- **Run GitLab and GitHub CI pipelines locally**: Execute steps defined in a `.gitlab-ci.yml` file.
- **Exclude stages**: Ignore specific stages from execution, such as deploy scripts

## Usage

```bash
local-ci [options]

Options
--file (required): Path to the GitLab or GitHub CI YAML file.
--type (default: gitlab): The type of file to parse
--skip (optional): A comma-separated list of stages to skip
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request to help improve the project.

## License
This project is licensed under the MIT License. See the LICENSE file for more details.