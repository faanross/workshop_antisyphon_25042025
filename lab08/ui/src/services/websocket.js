// WebSocket connection management service
import { ref } from 'vue';

// Reactive state to track connection status
const isConnected = ref(false);

// WebSocket instance
let socket = null;

// WebSocket server URL
const wsUrl = 'ws://localhost:8080/ws';

// Connect to WebSocket server
function connect() {
    // Don't connect if already connected
    if (socket && (socket.readyState === WebSocket.CONNECTING ||
        socket.readyState === WebSocket.OPEN)) {
        return;
    }

    // Create new WebSocket connection
    socket = new WebSocket(wsUrl);

    // Connection opened event
    socket.onopen = () => {
        console.log('WebSocket connected');
        isConnected.value = true;
    };

    // Connection closed event
    socket.onclose = () => {
        console.log('WebSocket disconnected');
        isConnected.value = false;
    };

    // Connection error event
    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

// Start connection when service is imported
connect();

// Export the service
export default {
    isConnected,
    connect
};