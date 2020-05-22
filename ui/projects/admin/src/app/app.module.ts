import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { BasePathService } from 'projects/shared/src/app/services/base-path.service';
import { JoinComponent } from './join/join.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { FormComponent } from './join/form/form.component';
import { AdminShellComponent } from './admin-shell/admin-shell.component';
import { LoginComponent } from './login/login.component';
import { SharedModule } from './shared/shared.module';
import { AccountModule } from './account/account.module';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ShareSiteModule } from './share-site/share-site.module';
import { CollectionModule } from './collection/collection.module';

@NgModule({
  declarations: [
    AdminShellComponent,
    AppComponent,
    DashboardComponent,
    FormComponent,
    JoinComponent,
    LoginComponent,
    NotFoundComponent,
  ],
  imports: [
    AccountModule,
    AppRoutingModule,
    BrowserModule,
    CollectionModule,
    FormsModule,
    HttpClientModule,
    SharedModule,
    ShareSiteModule,
  ],
  providers: [BasePathService],
  bootstrap: [AppComponent],
})
export class AppModule {}
