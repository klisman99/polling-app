import { mockDeep } from 'jest-mock-extended';
import { Pool } from 'pg';

process.env.JWT_SECRET = 'test-secret';

jest.mock('../src/db', () => ({
	query: jest.fn(),
	default: mockDeep<Pool>(),
}));

jest.spyOn(console, 'error').mockImplementation(() => {});