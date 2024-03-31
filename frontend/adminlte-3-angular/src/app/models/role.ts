import { RolePolicy } from "./rolepolicy";

export class Role {
    id: number;
    code: string;
    name: string;
    description: string;
    rolePolicies: RolePolicy[];

    constructor() {
        this.id = 0;
        this.code = '';
        this.name = '';
        this.description = '';
        this.rolePolicies = [];
    }
}
