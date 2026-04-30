# Swarm Mainframe

Swarm is a polyglot task-orchestration system designed for automation, security workflows, and local agent management. It allows you to coordinate a swarm of independent bot binaries, controlled by a central Go-Mainframe and driven by Lua-based mission files.

## The Concept
Swarm separates the Logic (Lua), the Orchestration (Go Mainframe), and the Execution (Polyglot Bots). 

- Mainframe: The brain. A persistent Go daemon that manages the database, schedules tasks, and monitors bot lifecycles.
- Bots: The soldiers. Any executable binary (Go, C, Rust, Python-wrappers) that follows the Swarm-Header protocol.
- Tasks: The missions. Lua scripts that define complex workflows, dependencies, and automation rules.

## Tech Stack
- Core: Go (Golang)
- Database: SQLite (CGO-free)
- Scripting: Lua (Gopher-Lua)
- Communication: JSON over Stdin/Stdout/IPC

## Project Structure
- /cmd/swarm: Main entry point.
- /internal/db: Database schemas and persistence logic.
- /internal/models: Central data structures.
- /internal/bot: Bot registration and management logic.
- /bots: Directory for installed bot binaries.
- /tasks: Lua mission files.
- /data: Persistent storage (SQLite).

## Features (Work in Progress)
- [x] Persistent SQLite Backend
- [ ] Bot-Handshake & Registration Protocol
- [ ] Async Task Scheduler
- [ ] Lua-DSL Integration
- [ ] Real-time Bot Monitoring (swarm listen)

## Installation & Usage
1. Clone the repository:
   git clone https://github.com/whazzabii7/swarm.git

2. Build the mainframe:
   go build -o swarm ./cmd/swarm/main.go

3. Run the mainframe:
   ./swarm

## License
MIT - Feel free to use and expand.
