import { AccountRoutingModule } from './account-routing.module';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AccountDashboardComponent } from './account-dashboard/account-dashboard.component';

@NgModule({
  imports: [
    CommonModule,
    AccountRoutingModule
  ],
  declarations: [AccountDashboardComponent]
})
export class AccountModule { }
