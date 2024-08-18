import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {provideHttpClient, withInterceptorsFromDi} from "@angular/common/http";
import {RouterLink, RouterLinkActive, RouterOutlet} from "@angular/router";
import {MatToolbarModule} from "@angular/material/toolbar";
import {MatButtonModule} from "@angular/material/button";
import {MatIconModule} from "@angular/material/icon";
import {MatSidenavModule} from "@angular/material/sidenav";
import {MatListModule} from "@angular/material/list";

@NgModule({
  declarations: [],
  providers: [provideHttpClient(withInterceptorsFromDi())],
  imports: [
    CommonModule,

  ]
})
export class AppModule {
}
