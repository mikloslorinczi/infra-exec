# infra-exec

Execute shell commands remotely

This small architecture consist of three main part:

### Infra CLI
The CLI is a command liene tool to add / list / inspect tasks on the Infra Server\
Downloading logfiles of executed tasks from the server is also possible with the CLI.

### Infra Server
The Infra Server is responsible of storing the tasks in its small JSON database, and distribute them to Infra Client(s) to be executed.\
The server also stores the logfiles sent from client(s), which can be downloaded with the CLI.

### Infra Client
The Infra Client regulary polls the Infra Server for new tasks to execute. When it finds a match (based on Tags) it tries to execute the task and channel its output (both stdin & stderr) to a logfile and send it back to the Infra Server. 

## Test
