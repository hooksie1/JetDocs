# JetDocs

Uses NATS as a backend for storing markdown documents and viewing them as HTML pages.


## Usage

### Start the server with integrated NATS server
`jetdocs start`

### Star the server with existing NATS server
`jetdocs start --nats-urls=nats://10.0.0.5:4222`

### Sync a markdown file
`jetdocs sync file.md`

### Sync all files in current directory
`jetdocs sync -a`

### View markdown pages
http://localhost:8080/pages

# Disclaimer

This was written quickly and should NOT be used in production. There is probably a lot of things wrong in here.