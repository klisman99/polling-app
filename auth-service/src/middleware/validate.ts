import { Request, Response, NextFunction } from 'express';
import { body, validationResult } from 'express-validator';

export const validateSignup = [
    body('email').isEmail().withMessage('Invalid email format'),
    body('password')
        .isLength({ min: 8 })
        .withMessage('Password must be at least 8 characters'),
    body('role')
        .optional()
        .isIn(['user', 'admin'])
        .withMessage('Role must be user or admin'),
    (req: Request, res: Response, next: NextFunction) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            res.status(400).json({ errors: errors.array() });
        }
        next();
    },
];

export const validateLogin = [
    body('email').isEmail().withMessage('Invalid email format'),
    body('password').notEmpty().withMessage('Password is required'),
    (req: Request, res: Response, next: NextFunction) => {
        const errors = validationResult(req);
        if (!errors.isEmpty()) {
            res.status(400).json({ errors: errors.array() });
        }
        next();
    },
];