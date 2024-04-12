import { Component, OnInit } from '@angular/core';
import { RoleService } from '@services/roles/role.service';

@Component({
  selector: 'app-roles',
  templateUrl: './roles.component.html',
  styleUrl: './roles.component.scss'
})
export class RolesComponent implements OnInit {
  dtOptions: DataTables.Settings = {};
  displayTable: boolean = false;
  constructor(private roleService: RoleService){

  }
  ngOnInit(): void {
    this.roleService.getRoles().subscribe(data => {
      this.displayTable = true;
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
        }]
      };
    });
  }
}
