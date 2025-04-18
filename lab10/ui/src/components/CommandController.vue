<template>
  <div class="command-controller">
    <div class="command-buttons">
      <button
          v-for="cmd in commands"
          :key="cmd"
          @click="executeCommand(cmd)"
          :disabled="!isConnected"
      >
        {{ cmd }}
      </button>
    </div>

    <div class="results-panel">
      <div class="results-header">
        <h3>Command Results</h3>
        <button @click="clearResults" class="clear-button">Clear</button>
      </div>

      <div class="results-container">
        <div v-if="messages.length === 0" class="no-results">
          No commands executed yet
        </div>

        <div v-for="(result, index) in messages" :key="index" class="result-item">
          <div class="result-header">
            <div class="command-name">$ {{ result.command }}</div>
            <div class="status-badge" :class="result.status">{{ result.status }}</div>
          </div>
          <pre class="output">{{ result.output }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import websocketService from '../services/websocket';

export default {
  name: 'CommandController',
  setup(props) {
    // Available commands
    const commands = ['pwd', 'whoami', 'hostname'];

    // Execute a command
    function executeCommand(command) {
      websocketService.sendCommand(command);
    }

    // Clear results
    function clearResults() {
      websocketService.messages.value = [];
    }

    return {
      commands,
      executeCommand,
      clearResults,
      isConnected: websocketService.isConnected,
      messages: websocketService.messages
    };
  }
};
</script>

<style scoped>
.command-controller {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.command-buttons {
  display: flex;
  gap: 10px;
}

.command-buttons button {
  padding: 8px 16px;
  border: none;
  background-color: #007aff;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}

.command-buttons button:disabled {
  background-color: #97a5b5;
  cursor: not-allowed;
}

.results-panel {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background-color: #3b3b3b;
  border-bottom: 1px solid #ddd;
}

.results-header h3 {
  margin: 0;
  font-size: 16px;
}

.clear-button {
  padding: 4px 8px;
  background: none;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.results-container {
  max-height: 400px;
  overflow-y: auto;
  padding: 12px;
}

.no-results {
  color: #8e8e93;
  font-style: italic;
  text-align: center;
  padding: 20px 0;
}

.result-item {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.result-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.command-name {
  font-family: monospace;
  font-weight: bold;
}

.status-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 10px;
}



.output {
  background-color: #1d1d1d;
  color: #f5f5f5;
  padding: 12px;
  border-radius: 4px;
  margin: 0;
  overflow-x: auto;
  font-size: 14px;
  line-height: 1.5;
}
</style>