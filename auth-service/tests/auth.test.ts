import supertest from 'supertest';
import express from 'express';
import helmet from 'helmet';

import { query } from '../src/db';
import authRouter from '../src/routes/auth';
import { errorHandler } from '../src/middleware/error';

const mockQuery = query as jest.Mock;

const app = express();
app.use(helmet());
app.use(express.json());
app.use('/auth', authRouter);
app.use(errorHandler);

describe('Auth Endpoints', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    describe('POST /auth/signup', () => {
        it('should create a new user and return a JWT', async () => {
            mockQuery
                .mockResolvedValueOnce({ rowCount: 0 })
                .mockResolvedValueOnce({
                    rows: [{ id: 1, role: 'user' }],
                    rowCount: 1,
                });

            const response = await supertest(app)
                .post('/auth/signup')
                .send({ email: 'test@example.com', password: 'securepassword123', role: 'user' });
            console.log(response.body);
            expect(response.status).toBe(201);
            expect(response.body).toHaveProperty('token');
            expect(mockQuery).toHaveBeenCalledTimes(2);
        });

        it('should return 409 if user already exists', async () => {
            mockQuery.mockResolvedValueOnce({ rowCount: 1 });

            const response = await supertest(app)
                .post('/auth/signup')
                .send({ email: 'test@example.com', password: 'securepassword123' });

            expect(response.status).toBe(409);
            expect(response.body).toEqual({ error: 'User already exists' });
            expect(mockQuery).toHaveBeenCalledTimes(1);
        });

        it('should return 400 for invalid input', async () => {
            const response = await supertest(app)
                .post('/auth/signup')
                .send({ email: 'invalid', password: 'short' });

            expect(response.status).toBe(400);
            expect(response.body.errors).toBeInstanceOf(Array);
            expect(response.body.errors).toContainEqual(
                expect.objectContaining({ msg: 'Invalid email format' })
            );
            expect(response.body.errors).toContainEqual(
                expect.objectContaining({ msg: 'Password must be at least 8 characters' })
            );
            expect(mockQuery).not.toHaveBeenCalled();
        });
    });

    describe('POST /auth/login', () => {
        it('should log in a user and return a JWT', async () => {
            mockQuery.mockResolvedValueOnce({
                rows: [
                    {
                        id: 1,
                        email: 'test@example.com',
                        password: await import('bcrypt').then(bcrypt => bcrypt.hash('securepassword123', 10)),
                        role: 'user',
                    },
                ],
                rowCount: 1,
            });

            const response = await supertest(app)
                .post('/auth/login')
                .send({ email: 'test@example.com', password: 'securepassword123' });

            expect(response.status).toBe(200);
            expect(response.body).toHaveProperty('token');
            expect(mockQuery).toHaveBeenCalledTimes(1);
        });

        it('should return 401 for invalid credentials', async () => {
            mockQuery.mockResolvedValueOnce({ rowCount: 0 });

            const response = await supertest(app)
                .post('/auth/login')
                .send({ email: 'test@example.com', password: 'wrongpassword' });

            expect(response.status).toBe(401);
            expect(response.body).toEqual({ error: 'Invalid credentials' });
            expect(mockQuery).toHaveBeenCalledTimes(1);
        });

        it('should return 400 for invalid input', async () => {
            const response = await supertest(app)
                .post('/auth/login')
                .send({ email: 'invalid', password: '' });

            console.log('Signup invalid input response:', response.body); // Debug
            console.log('mockQuery calls:', mockQuery.mock.calls); // Debug

            expect(response.status).toBe(400);
            expect(response.body.errors).toBeInstanceOf(Array);
            expect(response.body.errors).toContainEqual(
                expect.objectContaining({ msg: 'Invalid email format' })
            );
            expect(response.body.errors).toContainEqual(
                expect.objectContaining({ msg: 'Password is required' })
            );
            expect(mockQuery).not.toHaveBeenCalled();
        });
    });
});