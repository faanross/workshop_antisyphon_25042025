// WebSocket connection management service
import { ref } from 'vue';

// Reactive state
const isConnected = ref(false);
const messages = ref([]);

// WebSocket instance
let socket = null;

// WebSocket server URL - adjust if needed
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

        // Attempt to reconnect after a delay
        setTimeout(() => {
            console.log('Attempting to reconnect...');
            connect();
        }, 3000);
    };

    // Connection error event
    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        isConnected.value = false;
    };

    // Incoming message event
    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log('Received message:', data);

            // Add message to history if it's a command response
            if (data.type === 'response') {
                messages.value.push(data);
            }
        } catch (error) {
            console.error('Error parsing message:', error);
        }
    };
}

// Send a command to the server
function sendCommand(command) {
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        console.error('Cannot send message, WebSocket not connected');
        return false;
    }

    const message = {
        type: 'command',
        command
    };

    socket.send(JSON.stringify(message));
    return true;
}

// Start connection when service is imported
connect();

// Export the service
export default {
    isConnected,
    messages,
    connect,
    sendCommand
};