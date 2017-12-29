import { AccountModule } from './account/account.module';
import { AuthGuard } from './auth.guard';
import { AuthService } from './auth.service';
import { LoginComponent } from './login/login.component';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AdminRoutingModule } from './admin-routing.module';
import { DashboardComponent } from './dashboard/dashboard.component';
import { AppComponent } from './app/app.component';
import { ShareSiteModule } from './share-site/share-site.module';
import { ShareSiteService, JWTInterceptor } from './services/share-site.service';
import { PathService } from './services/path.service';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { SessionService } from './services/session.service';
import { NotFoundComponent } from './not-found/not-found.component';

@NgModule({
  imports: [
    CommonModule,
    AdminRoutingModule,
    FormsModule,
    HttpClientModule
  ],
  providers: [
    AuthService,
    AuthGuard,
    ShareSiteService,
    PathService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: JWTInterceptor,
      multi: true
    },
    SessionService
  ],
  declarations: [
    DashboardComponent,
    AppComponent,
    LoginComponent,
    NotFoundComponent
  ]
})
export class AdminModule { }

