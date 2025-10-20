export interface User {
	id: string;
	name: string;
	email: string;
	createdAt: Date;
	updatedAt: Date;
}

export interface CreateUserDto {
	name: string;
	email: string;
}

export interface UpdateUserDto {
	name?: string;
	email?: string;
}

// In-memory database (replace with real database later)
export const users: User[] = [];
