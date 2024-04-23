import { Policy } from '@/models/policy';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { RoleService } from '@services/roles/role.service';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'app-role-edit',
  templateUrl: './role-edit.component.html',
  styleUrl: './role-edit.component.scss'
})
export class RoleEditComponent {
  @Input() roleId: string;
  @Output() onClose = new EventEmitter<void>();
  roleForm: FormGroup;
  policies: Policy[] = [];
  constructor(private roleService: RoleService, private fb: FormBuilder,
    private toastr: ToastrService,
  ) {
    this.roleForm = this.fb.group({
      Id: [''],  // Role ID
      Name: ['', [Validators.required, Validators.minLength(4)]],
      Code: ['', [Validators.required, Validators.minLength(4)]],
      rolePolicies: this.fb.array([])
    });
    this.policies = [
      { id: 1, name: 'UserManagement', code: 'UserManagement', description: '' },
      { id: 2, name: 'SubjectManagement', code: 'SubjectManagement', description: '' },
      { id: 3, name: 'SomeOtherManagement', code: 'SomeOtherManagement', description: '' }
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
          roleId: role.id,
          policyId: policy.id,
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
    this.onClose.emit();
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
          let msg = `Error updating role: ${error.error.error}`
          console.error(msg);
          this.toastr.error(msg);          
        }
      });
    } else {
      console.log('Role form is not valid');
    }
  }
}
