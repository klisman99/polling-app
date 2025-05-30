import express from 'express';
import helmet from 'helmet';
import { query } from './db';
import authRouter from './routes/auth';
import { errorHandler } from './middleware/error';

const app = express();

app.use(helmet());
app.use(express.json());

app.use('/auth', authRouter);

app.get('/health', async (req, res) => {
	try {
		await query('SELECT 1');
		res.json({ status: 'Auth service running, DB connected' });
	} catch (err) {
		res.status(500).json({ status: 'DB connection failed', error: (err as Error).message });
	}
});

app.use(errorHandler);

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
	console.log(`Auth service running on port ${PORT}`);
});