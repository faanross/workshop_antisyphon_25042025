<template>
  <div class="app">
    <header>
      <h1>Orlok C2 Command Center</h1>
      <ConnectionStatus />
    </header>

    <main>
      <CommandController />
    </main>
  </div>
</template>

<script>
// Script stays the same
import ConnectionStatus from './components/ConnectionStatus.vue';
import CommandController from './components/CommandController.vue';

export default {
  components: {
    ConnectionStatus,
    CommandController
  }
}
</script>

<style>
body {
  background-color: #1a1a1a;
  color: #f0f0f0;
  font-family: sans-serif;
  margin: 0;
  padding: 20px;
}

.app {
  max-width: 800px;
  margin: 0 auto;
}

header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  margin-bottom: 30px;
}

h1 {
  font-size: 28px;
  margin: 0 0 16px 0;
  color: #e6e978;
}

main {
  margin-top: 20px;
}
</style>