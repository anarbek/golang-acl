import {AppState} from '@/store/state';
import {UiState} from '@/store/ui/state';
import {Component, HostBinding, OnInit} from '@angular/core';
import {Store} from '@ngrx/store';
import {AppService} from '@services/app.service';
import { AuthService } from '@services/auth.service';
import {Observable} from 'rxjs';

const BASE_CLASSES = 'main-sidebar elevation-4';
@Component({
    selector: 'app-menu-sidebar',
    templateUrl: './menu-sidebar.component.html',
    styleUrls: ['./menu-sidebar.component.scss']
})
export class MenuSidebarComponent implements OnInit {
    @HostBinding('class') classes: string = BASE_CLASSES;
    public ui: Observable<UiState>;
    public user;
    public menu: MenuItem[] = [];

    constructor(
        public appService: AppService,
        private store: Store<AppState>,
        private auth: AuthService
    ) {}

    ngOnInit() {
        this.ui = this.store.select('ui');
        this.ui.subscribe((state: UiState) => {
            this.classes = `${BASE_CLASSES} ${state.sidebarSkin}`;
        });
        this.user = this.appService.user;
        console.log(this.user);
        this.menu = this.filterMenuByPermissions(MENU);
    }

    filterMenuByPermissions(menu: any[]): any[] {
        const permissions = this.auth.getPermissions();
        return menu.filter(item => {
            // If the item has a permission property, check if it's in the permissions array
            if (item.permission && !permissions.includes(item.permission)) {
                return false;
            }
            // If the item has children, filter them recursively
            if (item.children) {
                item.children = this.filterMenuByPermissions(item.children);
            }
            return true;
        });
    }
}

interface MenuItem {
    name: string;
    permission: string;
    iconClasses: string;
    path?: string[];
    children?: MenuItem[];
  }

export const MENU = [
    {
        name: 'Dashboard',
        permission: 'UserManagement.r',
        iconClasses: 'fas fa-tachometer-alt',
        path: ['/']
    },
    {
        name: 'Dashboard2',
        permission: 'test.r',
        iconClasses: 'fas fa-tachometer-alt',
        path: ['/']
    },
    {
        name: 'Blank',
        permission: 'UserManagement.r',
        iconClasses: 'fas fa-file',
        path: ['/blank']
    },
    {
        name: 'Main Menu',
        permission: 'UserManagement.r',
        iconClasses: 'fas fa-folder',
        children: [
            {
                name: 'Roles',
                permission: 'UserManagement.w',
                iconClasses: 'far fa-address-book',
                path: ['/roles']
            },
            {
                name: 'Sub Menu',
                permission: 'UserManagement.r',
                iconClasses: 'far fa-address-book',
                path: ['/sub-menu-1']
            },
            {
                name: 'Blank',
                permission: 'UserManagement.r',
                iconClasses: 'fas fa-file',
                path: ['/sub-menu-2']
            }
        ]
    }
];
