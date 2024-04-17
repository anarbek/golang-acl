import { Policy } from '@/models/policy';
import { Component, Input } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { RoleService } from '@services/roles/role.service';

@Component({
  selector: 'app-role-edit',
  templateUrl: './role-edit.component.html',
  styleUrl: './role-edit.component.scss'
})
export class RoleEditComponent {
  @Input() roleId: string;
  roleForm: FormGroup;
  policies: Policy[] = [];
  constructor(private roleService: RoleService, private fb: FormBuilder) {
    this.roleForm = this.fb.group({
      Id: [''],  // Role ID
      Name: ['', [Validators.required, Validators.minLength(4)]],
      Code: ['', [Validators.required, Validators.minLength(4)]],
      rolePolicies: this.fb.array([])
    });
    this.policies = [
      { id: 1, name: 'UserManagement', code: 'UserManagement', description: '' },
      { id: 2, name: 'RoleManagement', code: 'RoleManagement', description: '' },
      { id: 3, name: 'SubjectManagement', code: 'SubjectManagement', description: '' }
    ];
  }

  createRolePolicy(policy: any): FormGroup {
    return this.fb.group({      
      policyId: [policy.policy.id],
      roleId: [policy.roleId],
      policyName: [policy.policy.name],
      read: [policy.read],
      write: [policy.write]
    });
  }

  get rolePolicies(): FormArray {
    return this.roleForm.get('rolePolicies') as FormArray;
  }

  openEditModal(data: any) {
    this.roleId = data.id;
    console.log('id: ', data.id);

    // Convert roleId to number
    const roleIdNumber = Number(this.roleId);

    // Get role detail
    this.roleService.getRole(roleIdNumber).subscribe(role => {
      // TODO: Do something with the role data
      console.log(role);
      this.roleForm.patchValue({
        Id: role.id,
        Name: role.name,
        Code: role.code
      });

      // Create form groups for the role policies
      const rolePoliciesFGs = this.policies.map(policy => {
        // Find if the current policy is in the role's policies
        const rolePolicy = role.rolePolicies ? role.rolePolicies.find(rp => rp.policy.id === policy.id) : null;

        // If the policy is found, use its read and write values, otherwise set them to false
        return this.createRolePolicy({
          policy: policy,
          roleId: rolePolicy ? rolePolicy.roleId : -1,
          read: rolePolicy ? rolePolicy.read : false,
          write: rolePolicy ? rolePolicy.write : false
        });
      });

      const rolePoliciesFormArray = this.fb.array(rolePoliciesFGs);
      this.roleForm.setControl('rolePolicies', rolePoliciesFormArray);
      // Open the modal
      $('#editModal').modal('show');
    }, error => {
      console.error('Error getting role: ', error);
    });
  }

  closeModal() {
    $('#editModal').modal('hide');
  }

  saveChanges() {
    if (this.roleForm.valid) {
      const roleIdNumber = Number(this.roleId);
      this.roleService.updateRole(roleIdNumber, this.roleForm.value).subscribe({
        next: response => {
          console.log('Role updated successfully', response);
          this.closeModal();
        },
        error: error => {
          console.error('Error updating role: ', error);
        }
      });
    } else {
      console.log('Role form is not valid');
    }
  }
}
