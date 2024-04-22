import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormArray } from '@angular/forms';
import { RoleService } from '@services/roles/role.service';

@Component({
  selector: 'app-role-create',
  templateUrl: './role-create.component.html',
  styleUrl: './role-create.component.scss'
})
export class RoleCreateComponent implements OnInit {
  roleForm: FormGroup;
  @Output() onClose = new EventEmitter<void>();
  policies = [
    { id: 1, name: 'UserManagement', code: 'UserManagement', description: '' },
    { id: 2, name: 'SubjectManagement', code: 'SubjectManagement', description: '' },
    { id: 3, name: 'SomeOtherManagement', code: 'SomeOtherManagement', description: '' }
  ];

  constructor(private fb: FormBuilder, private roleService: RoleService) {
    this.roleForm = this.fb.group({
      Name: ['', [Validators.required, Validators.minLength(4)]],
      Code: ['', [Validators.required, Validators.minLength(4)]],
      rolePolicies: this.fb.array([])
    });

    this.addRolePolicies();
  }

  ngOnInit(): void {}

  get rolePolicies(): FormArray {
    return this.roleForm.get('rolePolicies') as FormArray;
  }

  addRolePolicies() {
    this.policies.forEach(policy => {
      this.rolePolicies.push(this.fb.group({
        policyId: [policy.id],
        policyName: [policy.name],
        read: [false],
        write: [false]
      }));
    });
  }

  openModal() {
    $('#createModal').modal('show');
  }

  closeModal() {
    $('#createModal').modal('hide');
    this.onClose.emit();
  }

  saveChanges() {
    if (this.roleForm.valid) {
      this.roleService.insertRole(this.roleForm.value).subscribe({
        next: response => {
          console.log('Role inserted successfully', response);
          this.closeModal();
        },
        error: error => {
          console.error('Error inserting role: ', error);
        }
      });
    } else {
      console.log('Role form is not valid');
    }
  }
}