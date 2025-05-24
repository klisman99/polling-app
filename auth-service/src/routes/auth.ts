import { NextFunction, Request, Response, Router } from 'express';
import { AuthService } from '../services/auth';
import { validateSignup, validateLogin } from '../middleware/validate';

const router = Router();

router.post('/signup', validateSignup, async (req: Request, res: Response, next: NextFunction) => {
	try {
		const token = await AuthService.signup(req.body);
		res.status(201).json({ token });
	} catch (err) {
		next(err);
	}
});

router.post('/login', validateLogin, async (req: Request, res: Response, next: NextFunction) => {
	try {
		const token = await AuthService.login(req.body.email, req.body.password);
		res.json({ token });
	} catch (err) {
		next(err);
	}
});

export default router;