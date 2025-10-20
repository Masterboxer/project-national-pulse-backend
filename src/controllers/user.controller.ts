import { Request, Response } from "express";
import {
	User,
	CreateUserDto,
	UpdateUserDto,
	users,
} from "../models/user.model";
import { v4 as uuidv4 } from "uuid";

export const getAllUsers = (req: Request, res: Response) => {
	res.status(200).json({
		success: true,
		data: users,
	});
};

export const getUserById = (req: Request, res: Response) => {
	const { id } = req.params;
	const user = users.find((u) => u.id === id);

	if (!user) {
		return res.status(404).json({
			success: false,
			message: "User not found",
		});
	}

	res.status(200).json({
		success: true,
		data: user,
	});
};

export const createUser = (req: Request, res: Response) => {
	const { name, email }: CreateUserDto = req.body;

	// Check if email already exists
	const existingUser = users.find((u) => u.email === email);
	if (existingUser) {
		return res.status(400).json({
			success: false,
			message: "Email already exists",
		});
	}

	const newUser: User = {
		id: uuidv4(),
		name,
		email,
		createdAt: new Date(),
		updatedAt: new Date(),
	};

	users.push(newUser);

	res.status(201).json({
		success: true,
		data: newUser,
	});
};

export const updateUser = (req: Request, res: Response) => {
	const { id } = req.params;
	const { name, email }: UpdateUserDto = req.body;

	const userIndex = users.findIndex((u) => u.id === id);

	if (userIndex === -1) {
		return res.status(404).json({
			success: false,
			message: "User not found",
		});
	}

	// Check if email already exists (for other users)
	if (email) {
		const existingUser = users.find((u) => u.email === email && u.id !== id);
		if (existingUser) {
			return res.status(400).json({
				success: false,
				message: "Email already exists",
			});
		}
	}

	const updatedUser: User = {
		...users[userIndex],
		...(name && { name }),
		...(email && { email }),
		updatedAt: new Date(),
	};

	users[userIndex] = updatedUser;

	res.status(200).json({
		success: true,
		data: updatedUser,
	});
};

export const deleteUser = (req: Request, res: Response) => {
	const { id } = req.params;
	const userIndex = users.findIndex((u) => u.id === id);

	if (userIndex === -1) {
		return res.status(404).json({
			success: false,
			message: "User not found",
		});
	}

	users.splice(userIndex, 1);

	res.status(200).json({
		success: true,
		message: "User deleted successfully",
	});
};
