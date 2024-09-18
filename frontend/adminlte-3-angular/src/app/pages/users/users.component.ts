import { AfterViewInit, Component, OnInit, ViewChild } from '@angular/core';
import { UsersService } from '@services/users/users.service';
import 'bootstrap';
import { BehaviorSubject } from 'rxjs';
import Settings, { Config } from 'datatables.net';

@Component({
  selector: 'app-users',  
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'] 
})


export class UsersComponent implements OnInit , AfterViewInit {
  dtOptions: Config = {};
  displayTable: boolean = false;  
  private itemsSubject = new BehaviorSubject<any[]>([]);
  public items$ = this.itemsSubject.asObservable();

  constructor(private itemsService: UsersService){
    
  }

  ngOnInit(): void {    
    this.items$.subscribe(data => {
      this.dtOptions = {
        data: data,
        columns: [{
          title: 'ID',
          data: 'id'
        }, {
          title: 'Username',
          data: 'username'
        }, {
          title: 'Email',
          data: 'email'
        },
        {
          title: 'Edit',
          data: null,
          defaultContent: '<button class="btn btn-primary editBtn">Edit</button>',
          orderable: false
        }],
        rowCallback: (row: Node, data: any[] | Object, index: number) => {          
          //$('.editBtn', row).unbind('click');
          // $('.editBtn', row).bind('click', () => {
          //    this.roleEditComponent.openEditModal(data);
          // });          
          return row;
        }
      };
    });
    this.getItems();
  }

  getItems() {
    this.displayTable = false;
    this.itemsService.getItems().subscribe(res => {
      this.itemsSubject.next(res);
      this.displayTable = true;
    });
  }
  
  ngAfterViewInit(): void {
    
  }

}
