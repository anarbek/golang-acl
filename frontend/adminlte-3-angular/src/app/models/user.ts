import { Role } from "./role";

export class User {
    tenantId: number;
    id: number;
    username: string;
    email: string;
    role: Role; // Assuming Role is another class you've defined
    roleId: number;

    constructor() {
        this.tenantId = 0;
        this.id = 0;
        this.username = '';
        this.email = '';
        this.role = new Role(); // Assuming Role is another class you've defined
        this.roleId = 0;
    }
}