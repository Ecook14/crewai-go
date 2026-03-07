# Crew-GO Core Architecture 🏗️

Understanding the inner workings of Crew-GO allows developers to optimize performance and debug complex autonomous behavior.

## The Agent Reasoning Protocol (ReAct)

Crew-GO Agents do not just predictably "call an LLM". They utilize the **ReAct (Reason + Act)** pattern to create infinite, autonomous reasoning loops.

### The Inner Loop Segment (`pkg/agents/agent.go`)
When a Task is given to an Agent:
1.  **Context Construction**: The Engine compiles the Agent's Role, Goal, the current Task, AND the historical context of previous tasks into an immutable System Prompt.
2.  **Memory Recall**: The Agent calculates the vector embedding of the Task description and queries its `MemoryStore` (e.g., ChromaDB) for relevant past experiences, injecting them into context.
3.  **The LLM Call**: The LLM evaluates the prompt.
4.  **Action Evaluation**:
    *   If the LLM responds with a final answer, the loop BREAKS and the text is returned.
    *   If the LLM responds with a `ToolCall` (e.g., `SearchWeb(query="...")`), the loop pausses.
5.  **Execution & Observation**: The Go engine natively executes the requested Tool payload. The output (or error) is injected back into the LLM context as a new "Observation" message.
6.  **Self-Healing**: If a Go panic occurred or an API timed out inside the Tool, the LLM literally "reads" the Go error output and attempts to adjust parameters for the next iteration.
7.  **Loop Returns to Step 3**. This continues until the agent believes it has satisfied the Task.

---

## The Go-Native Orchestration Engine (`pkg/crew/crew.go`)

Unlike Python frameworks which rely on simulated asynchronous loops, Crew-GO utilizes true native hardware threads (`goroutines`).

### Hierarchical Processing Deep-Dive
When `Process: crew.Hierarchical` is invoked:
1.  The `Crew` generates an autonomous `ManagerAgent`.
2.  A Go `sync.WaitGroup` is initialized.
3.  A loop over all Tasks triggers a massive, instantaneous Fan-Out. Every single Task drops into its own `goroutine`.
4.  Inside each goroutine, the `ManagerAgent` is pinged with the Task payload to determine WHICH worker Agent is best suited.
5.  The worker Agents run their ReAct loops simultaneously across all cores.
6.  As tasks finish, an `errCh := make(chan error, len(tasks))` captures success limits.
7.  The `WaitGroup.Wait()` blocks until all parallel streams conclude.
8.  The `ManagerAgent` receives a fan-in of all results and synthesizes the finalized buffer.

---

## Global Telemetry & Observability Bus (`pkg/telemetry`)

AI operations are historically "black boxes". Crew-GO fixes this entirely with the Global Event Bus.

### Event Propagation
The `telemetry.EventBus` is an internal Pub/Sub broker wrapped in an `RWMutex`.

As an Agent performs operations deep within the call stack, it calls:
```go
telemetry.GlobalBus.Publish(telemetry.Event{
    Type:      telemetry.EventToolStarted,
    AgentRole: "Researcher",
    Payload:   map[string]interface{}{"tool": "SearchWeb"},
})
```

Because it uses isolated channels, this publishing has negligible impact on execution latency natively.

### The Dashboard Bridge (`internal/server/ws.go`)
When `--ui` is utilized, the `StartDashboardServer` function initializes:
1.  An HTTP handler for standard `html/css/js` delivery.
2.  A WebSocket `/ws` upgrade handler.
3.  It calls `telemetry.GlobalBus.Subscribe()`, grabbing the live firehose of ReAct events and marshaling them into JSON chunks over the TCP socket, giving the browser real-time frame rates of the AI reasoning cycle.
