import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';

import { query } from '../db';
import { User, UserInput, JwtPayload } from '../types';

export class AuthService {
    static async signup(input: UserInput): Promise<string> {
        const { email, password, role = 'user' } = input;

        // Check if user exists
        const existing = await query('SELECT id FROM users WHERE email = $1', [email]);
        if ((existing.rowCount ?? 0) > 0) {
            throw new Error('User already exists');
        }

        // Hash password
        const SALT_ROUNDS = 10;
        const hashedPassword = await bcrypt.hash(password, SALT_ROUNDS);

        // Insert user
        const result = await query(
            'INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id, role',
            [email, hashedPassword, role]
        );
        const user = result.rows[0];

        // Generate JWT
        const payload: JwtPayload = { id: user.id, role: user.role };
        const token = jwt.sign(payload, process.env.JWT_SECRET!, { expiresIn: '1h' });

        return token;
    }

    static async login(email: string, password: string): Promise<string> {
        // Find user
        const result = await query('SELECT * FROM users WHERE email = $1', [email]);
        if (result.rowCount === 0) {
            throw new Error('Invalid credentials');
        }
        const user: User = result.rows[0];

        // Verify password
        const isValid = await bcrypt.compare(password, user.password);
        if (!isValid) {
            throw new Error('Invalid credentials');
        }

        // Generate JWT
        const payload: JwtPayload = { id: user.id, role: user.role };
        const token = jwt.sign(payload, process.env.JWT_SECRET!, { expiresIn: '1h' });

        return token;
    }
}