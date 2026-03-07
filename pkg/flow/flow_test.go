package flow

import (
	"context"
	"testing"
)

func TestFlowKickoff(t *testing.T) {
	f := NewFlow(State{"counter": 0})

	f.AddNode(func(ctx context.Context, state State) (State, error) {
		state["counter"] = 1
		return state, nil
	})

	f.AddNode(func(ctx context.Context, state State) (State, error) {
		count := state["counter"].(int)
		state["counter"] = count + 1
		return state, nil
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if result["counter"] != 2 {
		t.Errorf("expected counter=2, got %v", result["counter"])
	}
}

func TestFlowStateMerge(t *testing.T) {
	f := NewFlow(State{"a": "original"})

	f.AddNode(func(ctx context.Context, state State) (State, error) {
		state["b"] = "added"
		return state, nil
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if result["a"] != "original" {
		t.Errorf("expected a=original, got %v", result["a"])
	}
	if result["b"] != "added" {
		t.Errorf("expected b=added, got %v", result["b"])
	}
}

func TestFlowContextCancellation(t *testing.T) {
	f := NewFlow(nil)
	f.AddNode(func(ctx context.Context, state State) (State, error) {
		state["reached"] = true
		return state, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := f.Kickoff(ctx)
	if err == nil {
		t.Error("expected context cancellation error, got nil")
	}
}

func TestConditionalNode(t *testing.T) {
	f := NewFlow(State{"should_run": true})

	// This node should execute because should_run is true
	f.AddConditionalNode("test-cond", func(state State) bool {
		return state["should_run"] == true
	}, func(ctx context.Context, state State) (State, error) {
		state["executed"] = true
		return state, nil
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if result["executed"] != true {
		t.Error("conditional node should have executed")
	}
}

func TestConditionalNodeSkipped(t *testing.T) {
	f := NewFlow(State{"should_run": false})

	f.AddConditionalNode("test-skip", func(state State) bool {
		return state["should_run"] == true
	}, func(ctx context.Context, state State) (State, error) {
		state["executed"] = true
		return state, nil
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if _, exists := result["executed"]; exists {
		t.Error("conditional node should NOT have executed")
	}
}

func TestRouterNode(t *testing.T) {
	f := NewFlow(State{"mode": "fast"})

	f.AddRouter(&RouterNode{
		Routes: []Route{
			{
				Name: "fast-route",
				Pred: func(state State) bool { return state["mode"] == "fast" },
				Node: func(ctx context.Context, state State) (State, error) {
					state["route_taken"] = "fast"
					return state, nil
				},
			},
			{
				Name: "slow-route",
				Pred: func(state State) bool { return state["mode"] == "slow" },
				Node: func(ctx context.Context, state State) (State, error) {
					state["route_taken"] = "slow"
					return state, nil
				},
			},
		},
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if result["route_taken"] != "fast" {
		t.Errorf("expected fast route, got %v", result["route_taken"])
	}
}

func TestParallelNodes(t *testing.T) {
	f := NewFlow(nil)

	f.AddParallelNodes([]Node{
		func(ctx context.Context, state State) (State, error) {
			state["worker_a"] = "done"
			return state, nil
		},
		func(ctx context.Context, state State) (State, error) {
			state["worker_b"] = "done"
			return state, nil
		},
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if result["worker_a"] != "done" || result["worker_b"] != "done" {
		t.Errorf("expected both workers done, got a=%v, b=%v", result["worker_a"], result["worker_b"])
	}
}

func TestListenerNode(t *testing.T) {
	f := NewFlow(State{"ready": true})

	f.AddListener(&ListenerConfig{
		Key:      "ready",
		Expected: true,
		Node: func(ctx context.Context, state State) (State, error) {
			state["listener_fired"] = true
			return state, nil
		},
	})

	result, err := f.Kickoff(context.Background())
	if err != nil {
		t.Fatalf("Flow.Kickoff failed: %v", err)
	}

	if result["listener_fired"] != true {
		t.Error("listener should have fired")
	}
}
