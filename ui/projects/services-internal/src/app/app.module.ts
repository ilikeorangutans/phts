import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { CookieService } from 'ngx-cookie-service';
import { StoreModule } from '@ngrx/store';
import { sessionReducer } from "./reducers/session.reducer";
import { EffectsModule } from "@ngrx/effects";

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BasePathService } from 'projects/shared/src/app/services/base-path.service';
import { AuthInterceptor } from './services/auth.interceptor';
import { SessionEffects } from './effects/session.effects';
import { StoreDevtoolsModule } from '@ngrx/store-devtools';
import { environment } from '../environments/environment';

@NgModule({
  declarations: [AppComponent],
  imports: [AppRoutingModule, BrowserModule, HttpClientModule,
    StoreModule.forRoot({ session: sessionReducer }, {}),
    EffectsModule.forRoot([SessionEffects]),
    StoreDevtoolsModule.instrument({ maxAge: 25, logOnly: environment.production })],    
  providers: [
    BasePathService,
    CookieService,
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true },
  ],
  bootstrap: [AppComponent],
})
export class AppModule { }
