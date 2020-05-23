import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
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
import { DashboardComponent } from './dashboard/dashboard.component';
import { JWTInterceptor } from './services/jwt-interceptor';

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
    AppRoutingModule,
    BrowserModule,
    FormsModule,
    HttpClientModule,
    SharedModule,
  ],
  providers: [
    BasePathService,
    { provide: HTTP_INTERCEPTORS, useClass: JWTInterceptor, multi: true },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
