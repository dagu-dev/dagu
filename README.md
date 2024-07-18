<p align="center">
  <img src="./assets/images/dagu-logo.webp" width="960" alt="dagu-logo">
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/dagu-dev/dagu">
    <img src="https://goreportcard.com/badge/github.com/dagu-dev/dagu" />
  </a>
  <a href="https://codecov.io/gh/dagu-dev/dagu">
    <img src="https://codecov.io/gh/dagu-dev/dagu/branch/main/graph/badge.svg?token=CODZQP61J2" />
  </a>
  <a href="https://github.com/dagu-dev/dagu/releases">
    <img src="https://img.shields.io/github/release/dagu-dev/dagu.svg" />
  </a>
  <a href="https://godoc.org/github.com/dagu-dev/dagu">
    <img src="https://godoc.org/github.com/dagu-dev/dagu?status.svg" />
  </a>
  <img src="https://github.com/dagu-dev/dagu/actions/workflows/test.yaml/badge.svg" />
</p>

<div align="center">

[Installation](https://dagu.readthedocs.io/en/latest/installation.html) | [Community](https://discord.gg/gpahPUjGRk) | [Quick Start](https://dagu.readthedocs.io/en/latest/quickstart.html)

</div>

<h1><b>Dagu</b></h1>

Dagu is a powerful Cron alternative that comes with a Web UI. It allows you to define dependencies between commands as a [Directed Acyclic Graph (DAG)](https://en.wikipedia.org/wiki/Directed_acyclic_graph) in a declarative [YAML format](https://dagu.readthedocs.io/en/latest/yaml_format.html). Dagu simplifies the management and execution of complex workflows. It natively supports running Docker containers, making HTTP requests, and executing commands over SSH.

-   [Documentation](https://dagu.readthedocs.io)
-   [Localized Documentation](#localized-documentation)
-   [Discord Community](https://discord.gg/gpahPUjGRk)

## **Highlights**

-   Single binary file installation
-   Declarative YAML format for defining DAGs
-   Web UI for visually managing, rerunning, and monitoring pipelines
-   Use existing programs without any modification
-   Self-contained, with no need for a DBMS

## **Table of Contents**

-   [**Highlights**](#highlights)
-   [**Table of Contents**](#table-of-contents)
-   [**Features**](#features)
-   [**Use Cases**](#use-cases)
-   [**Web UI**](#web-ui)
    -   [DAG Details](#dag-details)
    -   [DAGs](#dags)
    -   [Search](#search)
    -   [Execution History](#execution-history)
    -   [Log Viewer](#log-viewer)
-   [**Installation**](#installation)
    -   [Via Bash script](#via-bash-script)
    -   [Via GitHub Releases Page](#via-github-releases-page)
    -   [Via Homebrew (macOS)](#via-homebrew-macos)
    -   [Via Docker](#via-docker)
-   [**Quick Start Guide**](#quick-start-guide)
    -   [1. Launch the Web UI](#1-launch-the-web-ui)
    -   [2. Create a New DAG](#2-create-a-new-dag)
    -   [3. Edit the DAG](#3-edit-the-dag)
    -   [4. Execute the DAG](#4-execute-the-dag)
-   [**CLI**](#cli)
-   [**Localized Documentation**](#localized-documentation)
-   [**Documentation**](#documentation)
-   [**Running as a daemon**](#running-as-a-daemon)
-   [**Example DAG**](#example-dag)
-   [**Motivation**](#motivation)
-   [**Why Not Use an Existing DAG Scheduler Like Airflow?**](#why-not-use-an-existing-dag-scheduler-like-airflow)
-   [**How It Works**](#how-it-works)
-   [**License**](#license)
-   [**Support and Community**](#support-and-community)

## **Features**

-   Web User Interface
-   Command Line Interface (CLI) with several commands for running and managing DAGs
-   YAML format for defining DAGs, with support for various features including:
    -   Execution of custom code snippets
    -   Parameters
    -   Command substitution
    -   Conditional logic
    -   Redirection of stdout and stderr
    -   Lifecycle hooks
    -   Repeating task
    -   Automatic retry
-   Executors for running different types of tasks:
    -   Running arbitrary Docker containers
    -   Making HTTP requests
    -   Sending emails
    -   Running jq command
    -   Executing remote commands via SSH
-   Email notification
-   Scheduling with Cron expressions
-   REST API Interface
-   Basic Authentication over HTTPS

## **Use Cases**

-   **Data Pipeline Automation:** Schedule ETL tasks for data processing and centralization.
-   **Infrastructure Monitoring:** Periodically check infrastructure components with HTTP requests or SSH commands.
-   **Automated Reporting:** Generate and send periodic reports via email.
-   **Batch Processing:** Schedule batch jobs for tasks like data cleansing or model training.
-   **Task Dependency Management:** Manage complex workflows with interdependent tasks.
-   **Microservices Orchestration:** Define and manage dependencies between microservices.
-   **CI/CD Integration:** Automate code deployment, testing, and environment updates.
-   **Alerting System:** Create notifications based on specific triggers or conditions.
-   **Custom Task Automation:** Define and schedule custom tasks using code snippets.

## **Web UI**

### DAG Details

It shows the real-time status, logs, and DAG configurations. You can edit DAG configurations on a browser.

![example](assets/images/demo.gif?raw=true)

You can switch to the vertical graph with the button on the top right corner.

![Details-TD](assets/images/ui-details2.webp?raw=true)

### DAGs

It shows all DAGs and the real-time status.

![DAGs](assets/images/ui-dags.webp?raw=true)

### Search

It greps given text across all DAG definitions.
![History](assets/images/ui-search.webp?raw=true)

### Execution History

It shows past execution results and logs.

![History](assets/images/ui-history.webp?raw=true)

### Log Viewer

It shows the detail log and standard output of each execution and step.

![DAG Log](assets/images/ui-logoutput.webp?raw=true)

## **Installation**

You can install Dagu quickly using Homebrew or by downloading the latest binary from the Releases page on GitHub.

### Via Bash script

```sh
curl -L https://raw.githubusercontent.com/dagu-dev/dagu/main/scripts/installer.sh | bash
```

### Via GitHub Releases Page

Download the latest binary from the [Releases page](https://github.com/dagu-dev/dagu/releases) and place it in your `$PATH` (e.g. `/usr/local/bin`).

### Via Homebrew (macOS)

```sh
brew install dagu-dev/brew/dagu
```

Upgrade to the latest version:

```sh
brew upgrade dagu-dev/brew/dagu
```

### Via Docker

```sh
docker run \
--rm \
-p 8080:8080 \
-v $HOME/.config/dagu/dags:/home/dagu/.config/dagu/dags \
-v $HOME/.local/share/dagu:/home/dagu/.local/share/dagu \
ghcr.io/dagu-dev/dagu:latest dagu start-all
```

See [Environment variables](https://dagu.readthedocs.io/en/latest/config.html#environment-variables) to configure those default directories.

## **Quick Start Guide**

### 1. Launch the Web UI

Start the server and scheduler with the command `dagu start-all` and browse to `http://127.0.0.1:8080` to explore the Web UI.

### 2. Create a New DAG

Navigate to the DAG List page by clicking the menu in the left panel of the Web UI. Then create a DAG by clicking the `NEW` button at the top of the page. Enter `example` in the dialog.

_Note: DAG (YAML) files will be placed in `~/.config/dagu/dags` by default. See [Configuration Options](https://dagu.readthedocs.io/en/latest/config.html) for more details._

### 3. Edit the DAG

Go to the `SPEC` Tab and hit the `Edit` button. Copy & Paste the following example and click the `Save` button.

Example:

```yaml
schedule: "* * * * *" # Run the DAG every minute
steps:
    - name: s1
      command: echo Hello Dagu
    - name: s2
      command: echo done!
      depends:
          - s1
```

### 4. Execute the DAG

You can execute the example by pressing the `Start` button. You can see "Hello Dagu" in the log page in the Web UI.

## **CLI**

```sh
# Runs the DAG
dagu start [--params=<params>] <file>

# Displays the current status of the DAG
dagu status <file>

# Re-runs the specified DAG run
dagu retry --req=<request-id> <file>

# Stops the DAG execution
dagu stop <file>

# Restarts the current running DAG
dagu restart <file>

# Dry-runs the DAG
dagu dry [--params=<params>] <file>

# Launches both the web UI server and scheduler process
dagu start-all [--host=<host>] [--port=<port>] [--dags=<path to directory>]

# Launches the Dagu web UI server
dagu server [--host=<host>] [--port=<port>] [--dags=<path to directory>]

# Starts the scheduler process
dagu scheduler [--dags=<path to directory>]

# Shows the current binary version
dagu version
```

## **Localized Documentation**

-   [中文文档 (Chinese Documentation)](https://dagu.readthedocs.io/zh)
-   [日本語ドキュメント (Japanese Documentation)](https://dagu.readthedocs.io/ja)

## **Documentation**

-   [Installation Instructions](https://dagu.readthedocs.io/en/latest/installation.html)
-   ️[Quick Start Guide](https://dagu.readthedocs.io/en/latest/quickstart.html)
-   [Command Line Interface](https://dagu.readthedocs.io/en/latest/cli.html)
-   [Web User Interface](https://dagu.readthedocs.io/en/latest/web_interface.html)
-   Writing DAG
    -   [Minimal DAG Definition](https://dagu.readthedocs.io/en/latest/yaml_format.html#minimal-dag-definition)
    -   [Running Arbitrary Code Snippets](https://dagu.readthedocs.io/en/latest/yaml_format.html#running-arbitrary-code-snippets)
    -   [Environment Variables](https://dagu.readthedocs.io/en/latest/yaml_format.html#defining-environment-variables)
    -   [Parameters](https://dagu.readthedocs.io/en/latest/yaml_format.html#defining-and-using-parameters)
    -   [Command Substitution](https://dagu.readthedocs.io/en/latest/yaml_format.html#using-command-substitution)
    -   [Conditional Logic](https://dagu.readthedocs.io/en/latest/yaml_format.html#adding-conditional-logic)
    -   [Environment Variables with Standard Output](https://dagu.readthedocs.io/en/latest/yaml_format.html#setting-environment-variables-with-standard-output)
    -   [Redirecting Stdout and Stderr](https://dagu.readthedocs.io/en/latest/yaml_format.html#redirecting-stdout-and-stderr)
    -   [Lifecycle Hooks](https://dagu.readthedocs.io/en/latest/yaml_format.html#adding-lifecycle-hooks)
    -   [Repeating Task](https://dagu.readthedocs.io/en/latest/yaml_format.html#repeating-a-task-at-regular-intervals)
    -   [Minimal DAG Definition](https://dagu.readthedocs.io/en/latest/yaml_format.html#minimal-dag-definition)
    -   [Running Sub-DAG](https://dagu.readthedocs.io/en/latest/yaml_format.html#running-sub-dag)
    -   [All Available Fields for a DAG](https://dagu.readthedocs.io/en/latest/yaml_format.html#all-available-fields-for-dags)
    -   [All Available Fields for a Step](https://dagu.readthedocs.io/en/latest/yaml_format.html#all-available-fields-for-steps)
-   Example DAGs
    -   [Hello World](https://dagu.readthedocs.io/en/latest/examples.html#hello-world)
    -   [Conditional Steps](https://dagu.readthedocs.io/en/latest/examples.html#conditional-steps)
    -   [File Output](https://dagu.readthedocs.io/en/latest/examples.html#file-output)
    -   [Passing Output to Next Step](https://dagu.readthedocs.io/en/latest/examples.html#passing-output-to-next-step)
    -   [Running a Container Image](https://dagu.readthedocs.io/en/latest/examples.html#running-a-docker-container)
    -   [Making HTTP Requests](https://dagu.readthedocs.io/en/latest/examples.html#sending-http-requests)
    -   [JSON Processing](https://dagu.readthedocs.io/en/latest/examples.html#querying-json-data-with-jq)
    -   [Email](https://dagu.readthedocs.io/en/latest/examples.html#sending-email)
-   [Configurations](https://dagu.readthedocs.io/en/latest/config.html)
-   [Scheduler](https://dagu.readthedocs.io/en/latest/scheduler.html)
-   [Docker Compose](https://dagu.readthedocs.io/en/latest/docker-compose.html)
-   [REST API Documentation](https://app.swaggerhub.com/apis/YOHAMTA_1/dagu)

## **Running as a daemon**

The easiest way to make sure the process is always running on your system is to create the script below and execute it every minute using cron (you don't need `root` account in this way):

```bash
#!/bin/bash
process="dagu start-all"
command="/usr/bin/dagu start-all"

if ps ax | grep -v grep | grep "$process" > /dev/null
then
    exit
else
    $command &
fi

exit
```

## **Example DAG**

This example DAG showcases a data pipeline typically implemented in DevOps and Data Engineering scenarios. It demonstrates an end-to-end data processing cycle starting from data acquisition and cleansing to transformation, loading, analysis, reporting, and ultimately, cleanup.

![Details-TD](assets/images/example.webp?raw=true)

The YAML code below represents this DAG:

```yaml
# Environment variables used throughout the pipeline
env:
    - DATA_DIR: /data
    - SCRIPT_DIR: /scripts
    - LOG_DIR: /log
    # ... other variables can be added here

# Handlers to manage errors and cleanup after execution
handlerOn:
    failure:
        command: "echo error"
    exit:
        command: "echo clean up"

# The schedule for the DAG execution in cron format
# This schedule runs the DAG daily at 12:00 AM
schedule: "0 0 * * *"

steps:
    # Step 1: Pull the latest data from a data source
    - name: pull_data
      command: "sh"
      script: |
          echo `date '+%Y-%m-%d'`
      output: DATE

    # Step 2: Cleanse and prepare the data
    - name: cleanse_data
      command: echo cleansing ${DATA_DIR}/${DATE}.csv
      depends:
          - pull_data

    # Step 3: Transform the data
    - name: transform_data
      command: echo transforming ${DATA_DIR}/${DATE}_clean.csv
      depends:
          - cleanse_data

    # Parallel Step 1: Load the data into a database
    - name: load_data
      command: echo loading ${DATA_DIR}/${DATE}_transformed.csv
      depends:
          - transform_data

    # Parallel Step 2: Generate a statistical report
    - name: generate_report
      command: echo generating report ${DATA_DIR}/${DATE}_transformed.csv
      depends:
          - transform_data

    # Step 4: Run some analytics
    - name: run_analytics
      command: echo running analytics ${DATA_DIR}/${DATE}_transformed.csv
      depends:
          - load_data

    # Step 5: Send an email report
    - name: send_report
      command: echo sending email ${DATA_DIR}/${DATE}_analytics.csv
      depends:
          - run_analytics
          - generate_report

    # Step 6: Cleanup temporary files
    - name: cleanup
      command: echo removing ${DATE}*.csv
      depends:
          - send_report
```

## **Motivation**

Legacy systems often have complex and implicit dependencies between jobs. When there are hundreds of cron jobs on a server, it can be difficult to keep track of these dependencies and to determine which job to rerun if one fails. It can also be a hassle to SSH into a server to view logs and manually rerun shell scripts one by one. Dagu aims to solve these problems by allowing you to explicitly visualize and manage pipeline dependencies as a DAG, and by providing a web UI for checking dependencies, execution status, and logs and for rerunning or stopping jobs with a simple mouse click.

Dagu addresses these pain points by providing a user-friendly solution for explicitly defining and visualizing workflows. With its intuitive web UI, Dagu simplifies the management of workflows, enabling users to easily check dependencies, monitor execution status, view logs, and control job execution with just a few clicks.

## **Why Not Use an Existing DAG Scheduler Like Airflow?**

There are many existing tools such as Airflow, but many of these require you to write code in a programming language like Python to define your DAG. For systems that have been in operation for a long time, there may already be complex jobs with hundreds of thousands of lines of code written in languages like Perl or Shell Script. Adding another layer of complexity on top of these codes can reduce maintainability. Dagu was designed to be easy to use, self-contained, and require no coding, making it ideal for small projects.

## **How It Works**

Dagu is a single command line tool that uses the local file system to store data, so no database management system or cloud service is required. DAGs are defined in a declarative YAML format, and existing programs can be used without modification.

---

Feel free to contribute in any way you want! Share ideas, questions, submit issues, and create pull requests. Check out our [Contribution Guide](https://dagu.readthedocs.io/en/latest/contrib.html) for help getting started.

We welcome any and all contributions!

## **License**

This project is licensed under the GNU GPLv3.

## **Support and Community**

Join our [Discord community](https://discord.gg/gpahPUjGRk) to ask questions, request features, and share your ideas.
