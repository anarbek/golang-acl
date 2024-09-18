import { Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { RolesComponent } from './components/roles/roles.component';

export const routes: Routes = [
    { 
        path: '', 
        component: AppComponent 
      },
      { 
        path: 'roles', 
        component: RolesComponent,
        data: { roles: ['ADMIN'] }
      },
    //   { 
    //     path: 'login', 
    //     component: LoginComponent
    //   },
      { 
        path: '**', 
        redirectTo: '' 
      }
];
