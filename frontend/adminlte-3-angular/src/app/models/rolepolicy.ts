import { Policy } from "./policy";
import { Role } from "./role";

export class RolePolicy {
    roleId: number;
    role: Role;
    policy: Policy;
    policyId: number;
    read: boolean;
    write: boolean;

    constructor() {
        this.roleId = 0;
        this.role = new Role();
        this.policy = new Policy();
        this.policyId = 0;
        this.read = false;
        this.write = false;
    }
}