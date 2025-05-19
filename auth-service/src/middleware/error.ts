import { Request, Response, NextFunction } from 'express';

export function errorHandler(
    err: Error,
    req: Request,
    res: Response,
    next: NextFunction
) {
    console.error(err.stack);

    if (err.message === 'User already exists') {
        res.status(409).json({ error: err.message });
    }

    if (err.message === 'Invalid credentials') {
        res.status(401).json({ error: err.message });
    }

    res.status(500).json({ error: 'Internal server error' });
}