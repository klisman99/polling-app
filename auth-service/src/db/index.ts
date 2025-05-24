import { Pool } from 'pg';
import dotenv from 'dotenv';

dotenv.config();

const pool = new Pool({
	connectionString: process.env.DATABASE_URL,
});

export async function query(text: string, params: any[] = []) {
	const client = await pool.connect();
	try {
		const res = await client.query(text, params);
		return res;
	} finally {
		client.release();
	}
}

export default pool;