# Distributed Personal Knowledge Graph (DPKG)

A distributed system for building, querying, and sharing personal knowledge graphs. Designed with real-time collaboration, offline-first functionality, and data privacy in mind.

## Features
- Core graph engine with nodes and edges
- REST API for graph operations
- Extensible architecture for real-time updates and storage persistence

## Getting Started

### Prerequisites
- Go 1.20+
- Git

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/dpkgraph.git
   cd dpkgraph

2. Build and run the project:
    ```bash
    go run cmd/main.go

| Method | Endpoint                                   | Description           |
|--------|--------------------------------------------|-----------------------|
| POST   | `/nodes`                                   | Add a new node        |
| GET    | `/nodes/:id`                               | Get details of a node |
| DELETE | `/nodes/:id`                               | Delete a node         |
| POST   | `/edges`                                   | Add a new edge        |
| GET    | `/edges?from=<from>&to=<to>&label=<label>` | Get edges             |


### **Add TODOs as GitHub Issues**
Track your next steps using GitHub issues.

#### **Steps:**
1. Open your GitHub repository.
2. Go to the "Issues" tab.
3. Create issues for each TODO item. Example titles:
   - Implement persistent storage for graph data.
   - Add querying features for nodes and edges.
   - Enable real-time collaboration with WebSockets.
   - Add user authentication with JWT and roles.
   - Build a React-based frontend for graph visualization.

For each issue, provide a description, acceptance criteria, and labels (e.g., `enhancement`, `backend`, `frontend`).
