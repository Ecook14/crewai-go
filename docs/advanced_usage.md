# Crew-GO Advanced Orchestration 🔥

This guide covers the most technically complex, industrial-grade features of the Crew-GO engine: Parallelism, Cyclic Graphs, Structured Outputs, and Flow State Machines.

---

## 1. Top-Level Process Orchestration Modes

When you create a `Crew`, the `Process` enum dictates exactly how the graph of tasks is resolved.

### A. Sequential (`crew.Sequential`)
Tasks execute strictly in the order they define in the slice.
- **Use Case**: Simple data pipelines. `Scrape -> Summarize -> Save`.

### B. Hierarchical (`crew.Hierarchical`)
Task definitions execute **in parallel**, but they are not executed blindly.
1. The Crew spins up an invisible autonomous `ManagerAgent`.
2. The Manager looks at the parallel tasks and dynamically assigns them to the worker Agents in your slice.
3. Once all parallel go-routines finish, the Manager synthesizes all results into a single final master output.
- **Use Case**: Running massive, concurrent research sweeps (e.g. Scrape 10 websites at once) and merging the results.

### C. Consensual (`crew.Consensual`)
Takes a single task and launches it across **every Agent** concurrently. The ManagerAgent then forces consensus.
- **Use Case**: Evaluating high-risk decisions where you want 5 different LLM personas to vote on the best answer.

### D. Reflective (`crew.Reflective`)
Sequential execution, but after every single task, the ManagerAgent executes a rigorous `ReviewPrompt`. If it rejects the output, the worker agent must retry from scratch.

---

## 2. Cyclic Graphs & State Machines (Elite)

Standard frameworks use direct A -> B -> C flows. Crew-GO supports **DAGs with Cycles**. This means agents can get stuck in autonomous feedback loops until an exact condition is met.

### Enabling Graph Mode
```go
myCrew := crew.Crew{
    Agents:  agents,
    Tasks:   tasks,
    Process: crew.Graph, // Enable DAG Mode
    MaxCycles: 50,       // Global infinite-loop protection
}
```

### Defining Cycles via `OutputCondition`
You bind a function to a Task that analyzes its output. Depending on the return string, the task routes to a new path.

```go
codeTask := &tasks.Task{
    Description: "Write a sorting algorithm in Go.",
    Agent:       coder,
}

testTask := &tasks.Task{
    Description: "Test the provided code. If it fails, output 'FAIL'. If it passes, output 'PASS'.",
    Agent:       qaAgent,
    Dependencies: []*tasks.Task{codeTask}, // testTask waits for codeTask
}

// Map the outcomes
testTask.NextPaths = map[string]*tasks.Task{
    "retry":   codeTask, // Cycle backwards!
    "success": deployTask, // Move forwards
}

// Evaluate the specific outcome
testTask.OutputCondition = func(result interface{}) string {
    resStr := fmt.Sprintf("%v", result)
    if strings.Contains(resStr, "FAIL") {
        return "retry"
    }
    return "success"
}
```
If `testTask` fails, the Engine natively rewinds the state, marks `codeTask` as incomplete, and kicks off the execution loop again.

---

## 3. Structural JSON Output Extraction

CrewAI in Python utilizes `Pydantic` heavily to coerce LLMs into outputting JSON schemas. Crew-GO achieves this entirely natively using Go compilation structs.

### How it works
You pass a pointer of a Go Struct to a `Task`. Crew-GO dynamically reads the `json` tags via reflection, builds a JSON-Schema definition, injects it into the LLM system prompt, forces `JSON Mode` on the OpenAI protocol, and unmarshals the response back into your pointer.

```go
// 1. Define your Strict Schema
type FinancialReport struct {
    CompanyTicker string   `json:"company_ticker"`
    BullPoints    []string `json:"bull_points"`
    BearPoints    []string `json:"bear_points"`
    RiskScore     int      `json:"risk_score_1_to_10"`
}

var finalReport FinancialReport

// 2. Bind the pointer to the Task
analystTask := &tasks.Task{
    Description:  "Analyze AAPL's latest Q3 earnings.",
    Agent:        analyst,
    OutputSchema: &finalReport, // MAGIC HAPPENS HERE
}

myCrew.Kickoff(ctx)

// 3. Immediately use Native Go Structs!
fmt.Printf("Risk Score: %d\n", finalReport.RiskScore)
for _, point := range finalReport.BullPoints {
    fmt.Println("+", point)
}
```

---

## 4. Multi-Crew Orchestration (`pkg/flow`)

If `crew.Crew` is a microservice, `flow.Flow` is the Kubernetes that connects them all together. Flows are reactive State Machines that pipe shared state dictionaries `map[string]interface{}` between distinct functional Nodes.

### Creating a Flow
```go
import "github.com/Ecook14/crewai-go/pkg/flow"

// Initial State Dictionary
initialState := flow.State{"company": "OpenAI"}
f := flow.NewFlow(initialState)
```

### Adding Parallel Nodes
Run multiple independent AI crews or Go functions concurrently, modifying the shared state.
```go
f.AddParallelNodes([]flow.Node{
    fetchFinancialsNode,
    fetchNewsNode,
    fetchSocialSentimentNode,
})
```

### Adding Conditional Branhes (Routers)
Evaluate the state and completely branch your application execution.
```go
f.AddRouter(&flow.RouterNode{
    Routes: []flow.Route{
        {
            Name: "buy-stock",
            Pred: func(s flow.State) bool { return s["sentiment"] == "very_bullish" },
            Node: executeTradeCrew, // Kickoff a trading Crew
        },
        {
            Name: "hold-position",
            Pred: func(s flow.State) bool { return s["sentiment"] == "neutral" },
            Node: logAlertCrew, // Kickoff an alerting Crew
        },
    },
    DefaultNode: fallbackNode,
})
```

```go
// Kickoff the entire Flow infrastructure
finalMasterState, err := f.Kickoff(context.Background())
```
