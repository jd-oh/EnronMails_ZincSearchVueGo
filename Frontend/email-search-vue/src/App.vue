<template>
  <div id="app">
    <header>
      <h1>Enron Mail Explorer</h1>
    </header>
    <main>
      <router-view /> 
    </main>

  </div>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style>

#app {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
}
header {
  background: #4CAF50;
  color: white;
  padding: 1rem;
  text-align: center;
}
footer {
  background: #f1f1f1;
  color: #333;
  text-align: center;
  padding: 1rem;
  position: fixed;
  bottom: 0;
  width: 100%;
}
</style>
