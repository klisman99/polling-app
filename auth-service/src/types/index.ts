export interface User {
	id: number;
	email: string;
	password: string;
	role: 'user' | 'admin';
}

export interface UserInput {
	email: string;
	password: string;
	role?: 'user' | 'admin';
}

export interface JwtPayload {
	id: number;
	role: 'user' | 'admin';
}