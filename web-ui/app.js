const terminal = document.getElementById('terminal');
const agentList = document.getElementById('agent-list');
const statusDot = document.getElementById('status-dot');
const statusText = document.getElementById('status-text');

let socket;
const agents = new Map();

function connect() {
    socket = new WebSocket('ws://' + window.location.host + '/ws');

    socket.onopen = () => {
        statusDot.classList.add('active');
        statusText.innerText = 'CONNECTED';
        addLog('system', 'WebSocket connection established. Streaming live telemetry.');
    };

    socket.onclose = () => {
        statusDot.classList.remove('active');
        statusText.innerText = 'DISCONNECTED';
        addLog('system', 'Connection lost. Retrying in 5 seconds...');
        setTimeout(connect, 5000);
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        handleEvent(data);
    };
}

function handleEvent(event) {
    const { type, agent_role, payload, timestamp } = event;
    const time = new Date(timestamp).toLocaleTimeString();

    switch (type) {
        case 'agent_started':
            updateAgent(agent_role, 'Working', payload.input);
            addLog('system', `Agent [${agent_role}] initialized task execution.`, 'success');
            break;
        case 'agent_thinking':
            updateAgent(agent_role, 'Thinking');
            addLog('thinking', `Agent [${agent_role}] is processing iteration #${payload.iteration}...`);
            break;
        case 'tool_started':
            updateAgent(agent_role, 'Using Tool');
            addLog('tool', `Agent [${agent_role}] invoking tool: ${payload.tool}`);
            break;
        case 'tool_finished':
            addLog('tool', `Tool [${payload.tool}] returned result: ${payload.result ? payload.result.substring(0, 100) + '...' : 'OK'}`);
            break;
        case 'agent_finished':
            updateAgent(agent_role, 'Idle');
            addLog('system', `Agent [${agent_role}] successfully completed task.`, 'success');
            break;
        case 'task_started':
            addLog('system', `New Task Started: ${payload.description.substring(0, 100)}...`);
            break;
        case 'task_finished':
            addLog('system', `Task Completed successfully.`, 'success');
            break;
    }
}

function updateAgent(role, status, goal = '') {
    if (!agents.has(role)) {
        const agentEl = document.createElement('div');
        agentEl.className = 'agent-item';
        agentEl.id = `agent-${role.replace(/\s+/g, '-')}`;
        agentEl.innerHTML = `
            <div class="agent-info">
                <h4>${role}</h4>
                <p class="agent-status">${status}</p>
            </div>
            <div class="agent-dot active"></div>
        `;
        agentList.appendChild(agentEl);
        agents.set(role, { el: agentEl });
    } else {
        const agentData = agents.get(role);
        agentData.el.querySelector('.agent-status').innerText = status;
    }
}

function addLog(type, message, category = '') {
    const entry = document.createElement('div');
    entry.className = `log-entry type-${category || type}`;
    
    const timeSpan = document.createElement('span');
    timeSpan.className = 'log-time';
    timeSpan.innerText = new Date().toLocaleTimeString();

    const typeSpan = document.createElement('span');
    typeSpan.className = 'log-type';
    typeSpan.innerText = `[${type.toUpperCase()}]`;

    entry.appendChild(timeSpan);
    entry.appendChild(typeSpan);
    entry.appendChild(document.createTextNode(message));

    terminal.appendChild(entry);
    terminal.scrollTop = terminal.scrollHeight;

    // Limit log entries
    if (terminal.childElementCount > 100) {
        terminal.removeChild(terminal.firstChild);
    }
}

connect();
