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
  constructor(private roleService: RoleService, private fb: FormBuilder){
    this.roleForm = this.fb.group({
        roleName: ['', [Validators.required, Validators.minLength(5)]],
        roleCode: ['', [Validators.required, Validators.minLength(5)]],
        rolePolicies: this.fb.array([])
    });
    this.policies = [
      { id: 1, name: 'UserManagement', code: 'UserManagement', description:'' },
      { id: 2, name: 'RoleManagement', code: 'RoleManagement', description:''  },
      { id: 3, name: 'SubjectManagement', code: 'SubjectManagement', description:''  }
  ];
  }

  createRolePolicy(policy: any): FormGroup {
      return this.fb.group({
          policyId: [policy.policy.id],
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
        roleName: role.name,
        roleCode: role.code
      });
      
      // Create form groups for the role policies
      const rolePoliciesFGs = this.policies.map(policy => {
        // Find if the current policy is in the role's policies
        const rolePolicy = role.rolePolicies.find(rp => rp.policy.id === policy.id);

        // If the policy is found, use its read and write values, otherwise set them to false
        return this.createRolePolicy({
          policy: policy,
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
    // Add your code to save changes here.
    console.log('Saving changes...');
  }
}
