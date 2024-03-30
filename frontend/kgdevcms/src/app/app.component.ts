import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { MainFooterComponent } from './shared/components/main-footer/main-footer.component';
import { ControlSidebarComponent } from './shared/components/control-sidebar/control-sidebar.component';
import { ContentWrapperComponent } from './shared/components/content-wrapper/content-wrapper.component';
import { MainHeaderComponent } from './shared/components/main-header/main-header.component';
import { MainSidebarComponent } from './shared/components/main-sidebar/main-sidebar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, 
    MainFooterComponent, 
    ControlSidebarComponent,
     ContentWrapperComponent,
     MainHeaderComponent,
     MainSidebarComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'kgdevcms';
}
