import express from 'express';
import helmet from 'helmet';
import { query } from './db';
import authRouter from './routes/auth';
import { errorHandler } from './middleware/error';

const app = express();

// Security middleware
app.use(helmet());
app.use(express.json());

// Routes
app.use('/auth', authRouter);

// Health check
app.get('/health', async (req, res) => {
    try {
        await query('SELECT 1');
        res.json({ status: 'Auth service running, DB connected' });
    } catch (err) {
        res.status(500).json({ status: 'DB connection failed', error: (err as Error).message });
    }
});

// Error handling
app.use(errorHandler);

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
    console.log(`Auth service running on port ${PORT}`);
});