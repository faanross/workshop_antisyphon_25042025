<template>
  <div class="connection-status">
    <div class="indicator" :class="{ connected: isConnected }"></div>
    <span>{{ isConnected ? 'Connected' : 'Disconnected' }}</span>
  </div>
</template>

<script>
import websocketService from '../services/websocket';

export default {
  name: 'ConnectionStatus',
  setup() {
    return {
      isConnected: websocketService.isConnected
    };
  }
};
</script>

<style scoped>
.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
}

.indicator {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: #ff3b30;
  transition: background-color 0.3s ease;
}

.indicator.connected {
  background-color: #34c759;
}
</style>