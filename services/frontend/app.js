const express = require('express');
const axios = require('axios');
const app = express();
const PORT = 3000;

app.get('/', async (req, res) => {
  try {
    // Calling the Go Backend
    const response = await axios.get('http://localhost:8080/data');
    res.json({ 
      message: 'Hello from Frontend!', 
      backend_says: response.data 
    });
  } catch (error) {
    res.status(500).json({ error: 'Backend is unreachable' });
  }
});

app.listen(PORT, () => {
  console.log(`Frontend running on http://localhost:${PORT}`);
});