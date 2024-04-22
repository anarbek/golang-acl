import { AfterViewInit, Component, ComponentFactoryResolver, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { RoleService } from '@services/roles/role.service';
import { RoleEditComponent } from './role-edit/role-edit.component';
import 'bootstrap';
import { RoleCreateComponent } from './role-create/role-create.component';
import { BehaviorSubject } from 'rxjs';

@Component({
  selector: 'app-roles',
  templateUrl: './roles.component.html',
  styleUrl: './roles.component.scss'
})

export class RolesComponent implements OnInit , AfterViewInit {
  @ViewChild(RoleEditComponent) roleEditComponent: RoleEditComponent;
  @ViewChild(RoleCreateComponent) roleCreateComponent: RoleCreateComponent;
  dtOptions: DataTables.Settings = {};
  displayTable: boolean = false;  
  private rolesSubject = new BehaviorSubject<any[]>([]);
  public roles$ = this.rolesSubject.asObservable();
  
  constructor(private roleService: RoleService){

  }
  ngOnInit(): void {
    this.roles$.subscribe(data => {
      this.dtOptions = {
        data: data,
        columns: [{
          title: 'ID',
          data: 'id'
        }, {
          title: 'Name',
          data: 'name'
        }, {
          title: 'Code',
          data: 'code'
        },
        {
          title: 'Edit',
          data: null,
          defaultContent: '<button class="btn btn-primary editBtn">Edit</button>',
          orderable: false
        }],
        rowCallback: (row: Node, data: any[] | Object, index: number) => {          
          $('.editBtn', row).unbind('click');
          $('.editBtn', row).bind('click', () => {
             this.roleEditComponent.openEditModal(data);
          });          
          return row;
        }
      };
    });
    this.getRoles();
  }

  getRoles() {
    this.displayTable = false;
    this.roleService.getRoles().subscribe(res => {
      this.rolesSubject.next(res);
      this.displayTable = true;
    });
  }

  openCreateModal() {
    this.roleCreateComponent.openModal();
  }

  ngAfterViewInit(): void {
    /*const table = $('#yourTableId').DataTable(); // replace 'yourTableId' with the id of your table

    $('.editBtn').on('click', function() {
      const data = table.row($(this).parents('tr')).data();
      this.roleEditComponent.openEditModal(data);
    }.bind(this));*/
  }
}
