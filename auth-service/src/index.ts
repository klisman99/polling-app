import express from 'express';
import { query } from './db';

const app = express();
app.use(express.json());

// Health check
app.get('/health', async (req, res) => {
  try {
    await query('SELECT 1'); // Test DB connection
    res.json({ status: 'Auth service running, DB connected' });
  } catch (err) {
    res.status(500).json({ status: 'DB connection failed', error: (err as Error).message });
  }
});

// Test users table
app.get('/users', async (req, res) => {
  try {
    const result = await query('SELECT * FROM users');
    res.json(result.rows);
  } catch (err) {
    res.status(500).json({ error: (err as Error).message });
  }
});

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  console.log(`Auth service running on port ${PORT}`);
});