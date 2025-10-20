import { Request, Response, NextFunction } from "express";

export const validateCreateUser = (
	req: Request,
	res: Response,
	next: NextFunction,
) => {
	const { name, email } = req.body;

	if (!name || typeof name !== "string" || name.trim() === "") {
		return res.status(400).json({
			success: false,
			message: "Name is required and must be a non-empty string",
		});
	}

	if (!email || typeof email !== "string") {
		return res.status(400).json({
			success: false,
			message: "Email is required and must be a string",
		});
	}

	// Basic email validation
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	if (!emailRegex.test(email)) {
		return res.status(400).json({
			success: false,
			message: "Invalid email format",
		});
	}

	next();
};

export const validateUpdateUser = (
	req: Request,
	res: Response,
	next: NextFunction,
) => {
	const { name, email } = req.body;

	// At least one field should be provided
	if (!name && !email) {
		return res.status(400).json({
			success: false,
			message: "At least one field (name or email) must be provided",
		});
	}

	if (name !== undefined && (typeof name !== "string" || name.trim() === "")) {
		return res.status(400).json({
			success: false,
			message: "Name must be a non-empty string",
		});
	}

	if (email !== undefined) {
		if (typeof email !== "string") {
			return res.status(400).json({
				success: false,
				message: "Email must be a string",
			});
		}

		const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		if (!emailRegex.test(email)) {
			return res.status(400).json({
				success: false,
				message: "Invalid email format",
			});
		}
	}

	next();
};
