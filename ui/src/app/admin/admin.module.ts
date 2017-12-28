import { AccountModule } from './account/account.module';
import { AuthGuard } from './auth.guard';
import { AuthService } from './auth.service';
import { LoginComponent } from './login/login.component';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AdminRoutingModule } from './admin-routing.module';
import { DashboardComponent } from './dashboard/dashboard.component';
import { AppComponent } from './app/app.component';

@NgModule({
  imports: [
    CommonModule,
    AdminRoutingModule,
    FormsModule,
    HttpModule
  ],
  providers: [AuthService, AuthGuard],
  declarations: [
    DashboardComponent,
    AppComponent,
    LoginComponent
  ]
})
export class AdminModule { }
